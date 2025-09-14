// Package middleware 提供限流中间件功能
// 基于令牌桶算法实现高并发场景下的请求限流
// 支持多维度限流策略：按IP、用户ID、IP+用户ID等
// 特别针对语音交互模块优化，保证低延迟体验
package middleware

import (
	"ai_jianli_go/config"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitConfig 限流配置结构体
// 定义了限流器的核心参数和行为，支持基于用户角色的动态限流
type RateLimitConfig struct {
	// GetRateFunc 动态获取令牌生成速率的函数
	// 根据用户角色返回不同的QPS限制
	// 例如：普通用户返回50，会员用户返回200，管理员返回1000
	GetRateFunc func(*gin.Context) int

	// GetBurstFunc 动态获取桶容量的函数
	// 根据用户角色返回不同的burst值
	// 例如：普通用户返回100，会员用户返回400，管理员返回2000
	GetBurstFunc func(*gin.Context) int

	// KeyFunc 用于生成限流键的函数
	// 支持按不同维度限流：IP、用户ID、IP+用户ID等
	// 返回的键用于区分不同的限流对象
	// 例如：按IP限流返回 "ip:192.168.1.1"
	//      按用户限流返回 "user:12345"
	KeyFunc func(*gin.Context) string

	// SkipFunc 跳过限流的条件函数
	// 当返回true时，该请求将跳过限流检查
	// 通常用于管理员用户或特殊权限用户
	// 例如：检查用户角色是否为super_admin
	SkipFunc func(*gin.Context) bool

	// OnLimitReached 达到限流时的回调函数
	// 当请求被限流时调用，用于自定义限流响应
	// 参数：gin.Context和限流键
	// 通常用于设置HTTP 429状态码和Retry-After头
	OnLimitReached func(*gin.Context, string)
}

// RateLimiter 限流器结构体
// 管理多个限流器实例，每个限流键对应一个独立的令牌桶
// 增强版RateLimiter结构体，添加最后访问时间记录
type RateLimiter struct {
    // limiters 限流器映射表
    limiters map[string]*rate.Limiter
    
    // lastAccess 记录每个限流器最后访问时间
    lastAccess map[string]time.Time
    
    // mu 读写锁，保护并发访问
    mu sync.RWMutex
    
    // config 限流配置
    config RateLimitConfig
    
    // cleanupThreshold 清理阈值，超过这个时间未访问的限流器将被清理
    cleanupThreshold time.Duration
}

// NewRateLimiter 创建新的限流器实例
// 参数：
//   - config: 限流配置，包含速率、桶容量、键生成函数等
//
// 返回值：
//   - *RateLimiter: 初始化好的限流器实例
//
// 功能说明：
//  1. 设置默认的键生成函数（如果未提供）
//  2. 设置默认的限流回调函数（如果未提供）
//  3. 初始化限流器映射表和配置
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
    // 保留原有的默认配置设置
    if config.KeyFunc == nil {
        config.KeyFunc = func(c *gin.Context) string {
            return c.ClientIP()
        }
    }
    if config.OnLimitReached == nil {
        config.OnLimitReached = func(c *gin.Context, key string) {
            rate := config.GetRateFunc(c)
            c.Header("X-RateLimit-Limit", strconv.Itoa(rate))
            c.Header("X-RateLimit-Remaining", "0")
            c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
            c.Header("Retry-After", "1")
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error":       "请求过于频繁，请稍后重试",
                "code":        "RATE_LIMIT_EXCEEDED",
                "retry_after": 1,
            })
        }
    }

    return &RateLimiter{
        limiters:        make(map[string]*rate.Limiter),
        lastAccess:      make(map[string]time.Time), // 初始化最后访问时间映射表
        config:          config,
        cleanupThreshold: 30 * time.Minute, // 默认30分钟未访问则清理
    }
}

// getLimiter 获取或创建限流器实例（修改版）
func (rl *RateLimiter) getLimiter(c *gin.Context, key string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    rateValue := rl.config.GetRateFunc(c)
    burstValue := rl.config.GetBurstFunc(c)

    limiter, exists := rl.limiters[key]
    if !exists {
        limiter = rate.NewLimiter(rate.Limit(rateValue), burstValue)
        rl.limiters[key] = limiter
    }
    
    // 更新最后访问时间
    rl.lastAccess[key] = time.Now()
    
    return limiter
}

// Cleanup 清理过期的限流器（增强版）
func (rl *RateLimiter) Cleanup() {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    count := 0
    
    // 遍历所有限流器，清理超过阈值未访问的实例
    for key, lastAccessTime := range rl.lastAccess {
        if now.Sub(lastAccessTime) > rl.cleanupThreshold {
            delete(rl.limiters, key)
            delete(rl.lastAccess, key)
            count++
        }
    }
    
    // 可选：记录清理日志
    if count > 0 {
        // 注意：这里应使用实际的日志记录函数
        fmt.Printf("RateLimiter: 清理了%d个过期的限流器实例\n", count)
    }
}

// Allow 检查是否允许请求通过
// 参数：
//   - c: gin.Context，用于获取用户角色和动态限流参数
//   - key: 限流键，用于标识不同的限流对象
//
// 返回值：
//   - bool: true表示允许请求，false表示被限流
//
// 功能说明：
//  1. 获取对应键的限流器实例
//  2. 调用令牌桶的Allow方法检查是否有可用令牌
//  3. 如果有令牌则消耗一个并返回true，否则返回false
//  4. 这是一个非阻塞操作，立即返回结果
func (rl *RateLimiter) Allow(c *gin.Context, key string) bool {
	limiter := rl.getLimiter(c, key)
	return limiter.Allow()
}

// Wait 等待直到允许请求通过
// 参数：
//   - c: gin.Context，用于获取用户角色和动态限流参数
//   - key: 限流键，用于标识不同的限流对象
//
// 返回值：
//   - error: 如果等待过程中出现错误（如上下文取消）则返回错误
//
// 功能说明：
//  1. 获取对应键的限流器实例
//  2. 调用令牌桶的Wait方法等待可用令牌
//  3. 这是一个阻塞操作，会等待直到有令牌可用或上下文取消
//  4. 适用于需要保证请求最终被处理的场景
func (rl *RateLimiter) Wait(c *gin.Context, key string) error {
	limiter := rl.getLimiter(c, key)
	return limiter.Wait(c.Request.Context())
}

// Middleware 返回Gin中间件函数
// 返回值：
//   - gin.HandlerFunc: 符合Gin中间件接口的处理函数
//
// 功能说明：
//  1. 检查是否需要跳过限流（如管理员用户）
//  2. 生成限流键（按IP、用户ID等维度）
//  3. 检查是否允许请求通过
//  4. 记录请求统计信息到监控系统
//  5. 设置限流相关的响应头
//  6. 继续执行后续中间件或业务逻辑
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间，用于计算响应时间
		startTime := time.Now()

		// 检查是否需要跳过限流检查
		// 通常用于管理员用户或特殊权限用户
		if rl.config.SkipFunc != nil && rl.config.SkipFunc(c) {
			c.Next()
			return
		}

		// 使用配置的键生成函数生成限流键
		// 支持按IP、用户ID、IP+用户ID等不同维度限流
		key := rl.config.KeyFunc(c)

		// 检查是否允许请求通过
		// 如果返回false，说明请求被限流
		allowed := rl.Allow(c, key)

		// 使用defer确保请求统计信息被记录
		// 无论请求是否被限流，都需要记录统计信息
		defer func() {
			responseTime := time.Since(startTime)
			RecordRequest(key, !allowed, responseTime)
		}()

		// 如果请求被限流，执行限流回调并终止请求处理
		if !allowed {
			rl.config.OnLimitReached(c, key)
			c.Abort()
			return
		}

		// 请求通过限流检查，设置限流相关的响应头
		// 这些头信息帮助客户端了解当前的限流状态
		limiter := rl.getLimiter(c, key)
		rate := rl.config.GetRateFunc(c)
		c.Header("X-RateLimit-Limit", strconv.Itoa(rate))                                        // 限流速率
		c.Header("X-RateLimit-Remaining", strconv.Itoa(limiter.Burst()))                         // 剩余令牌数
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10)) // 重置时间

		// 继续执行后续的中间件或业务逻辑
		c.Next()
	}
}


// 辅助函数：从Gin上下文获取用户角色
func getUserRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	roleStr, ok := role.(string)
	if !ok {
		return ""
	}
	return roleStr
}

// 辅助函数：从Gin上下文获取用户ID
func getUserID(c *gin.Context) int64 {
	userID, exists := c.Get("id")
	if !exists {
		return 0
	}
	userIDInt64, ok := userID.(int64)
	if !ok {
		return 0
	}
	return userIDInt64
}

// 辅助函数：根据用户角色获取限流参数
func getRateLimitForRole(role string, defaultRate, defaultBurst int, roleLimits map[string]config.RoleRateLimit) (int, int) {
	// 如果角色在配置中存在，使用角色特定的限流参数
	if roleLimit, exists := roleLimits[role]; exists {
		return roleLimit.Rate, roleLimit.Burst
	}
	// 否则使用默认限流参数
	return defaultRate, defaultBurst
}

// 从配置文件获取限流配置的函数
// 这些函数从config包读取配置，避免硬编码

// GetSpeechRateLimitConfig 获取语音识别接口限流配置
func GetSpeechRateLimitConfig() RateLimitConfig {
	cfg := config.GetRateLimitConfig()
	speechCfg := cfg.Speech

	return RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			rate, _ := getRateLimitForRole(role, speechCfg.DefaultRate, speechCfg.DefaultBurst, speechCfg.RoleLimits)
			return rate
		},
		GetBurstFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			_, burst := getRateLimitForRole(role, speechCfg.DefaultRate, speechCfg.DefaultBurst, speechCfg.RoleLimits)
			return burst
		},
		KeyFunc: func(c *gin.Context) string {
			// 优先按用户ID限流，提供个性化服务
			// 如果没有用户ID，则回退到IP限流
			userID := getUserID(c)
			if userID > 0 {
				return fmt.Sprintf("user:%d", userID)
			}
			return fmt.Sprintf("ip:%s", c.ClientIP())
		},
		SkipFunc: func(c *gin.Context) bool {
			// 检查用户角色是否在跳过列表中
			role := getUserRole(c)
			for _, skipRole := range speechCfg.SkipRoles {
				if role == skipRole {
					return true
				}
			}
			return false
		},
	}
}

// GetGeneralRateLimitConfig 获取通用API限流配置
func GetGeneralRateLimitConfig() RateLimitConfig {
	cfg := config.GetRateLimitConfig()
	generalCfg := cfg.General

	return RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			rate, _ := getRateLimitForRole(role, generalCfg.DefaultRate, generalCfg.DefaultBurst, generalCfg.RoleLimits)
			return rate
		},
		GetBurstFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			_, burst := getRateLimitForRole(role, generalCfg.DefaultRate, generalCfg.DefaultBurst, generalCfg.RoleLimits)
			return burst
		},
		KeyFunc: func(c *gin.Context) string {
			// 按IP限流，简单有效
			return fmt.Sprintf("ip:%s", c.ClientIP())
		},
		SkipFunc: func(c *gin.Context) bool {
			// 检查用户角色是否在跳过列表中
			role := getUserRole(c)
			for _, skipRole := range generalCfg.SkipRoles {
				if role == skipRole {
					return true
				}
			}
			return false
		},
	}
}

// GetUploadRateLimitConfig 获取文件上传接口限流配置
func GetUploadRateLimitConfig() RateLimitConfig {
	cfg := config.GetRateLimitConfig()
	uploadCfg := cfg.Upload

	return RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			rate, _ := getRateLimitForRole(role, uploadCfg.DefaultRate, uploadCfg.DefaultBurst, uploadCfg.RoleLimits)
			return rate
		},
		GetBurstFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			_, burst := getRateLimitForRole(role, uploadCfg.DefaultRate, uploadCfg.DefaultBurst, uploadCfg.RoleLimits)
			return burst
		},
		KeyFunc: func(c *gin.Context) string {
			// 优先按用户限流，如果没有用户ID则按IP限流
			userID := getUserID(c)
			if userID > 0 {
				return fmt.Sprintf("user:%d", userID)
			}
			return fmt.Sprintf("ip:%s", c.ClientIP())
		},
		SkipFunc: func(c *gin.Context) bool {
			// 检查用户角色是否在跳过列表中
			role := getUserRole(c)
			for _, skipRole := range uploadCfg.SkipRoles {
				if role == skipRole {
					return true
				}
			}
			return false
		},
	}
}

// GetAuthRateLimitConfig 获取认证接口限流配置
func GetAuthRateLimitConfig() RateLimitConfig {
	cfg := config.GetRateLimitConfig()
	authCfg := cfg.Auth

	return RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			rate, _ := getRateLimitForRole(role, authCfg.DefaultRate, authCfg.DefaultBurst, authCfg.RoleLimits)
			return rate
		},
		GetBurstFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			_, burst := getRateLimitForRole(role, authCfg.DefaultRate, authCfg.DefaultBurst, authCfg.RoleLimits)
			return burst
		},
		KeyFunc: func(c *gin.Context) string {
			// 按IP限流，防止单个IP进行大量认证尝试
			return fmt.Sprintf("ip:%s", c.ClientIP())
		},
		SkipFunc: func(c *gin.Context) bool {
			// 检查用户角色是否在跳过列表中
			role := getUserRole(c)
			for _, skipRole := range authCfg.SkipRoles {
				if role == skipRole {
					return true
				}
			}
			return false
		},
	}
}

// 全局限流器实例
// 这些实例在应用启动时创建，供整个应用使用
// 每个实例对应一种特定的限流策略
var (
	// speechLimiter 语音识别接口限流器
	// 从配置文件读取配置，专门为语音交互优化
	speechLimiter *RateLimiter

	// generalLimiter 通用API限流器
	// 从配置文件读取配置，适用于大部分业务接口
	generalLimiter *RateLimiter

	// uploadLimiter 文件上传接口限流器
	// 从配置文件读取配置，专门针对文件上传场景
	uploadLimiter *RateLimiter

	// authLimiter 认证接口限流器
	// 从配置文件读取配置，专门用于防止暴力破解
	authLimiter *RateLimiter
)

// InitRateLimiters 初始化全局限流器实例
// 这个函数应该在应用启动时调用，从配置文件读取配置并创建限流器
func InitRateLimiters() {
	// 检查限流是否启用
	cfg := config.GetRateLimitConfig()
	if !cfg.Enabled {
		// 如果限流被禁用，创建空的限流器
		speechLimiter = nil
		generalLimiter = nil
		uploadLimiter = nil
		authLimiter = nil
		return
	}

	// 创建限流器实例
	speechLimiter = NewRateLimiter(GetSpeechRateLimitConfig())
	generalLimiter = NewRateLimiter(GetGeneralRateLimitConfig())
	uploadLimiter = NewRateLimiter(GetUploadRateLimitConfig())
	authLimiter = NewRateLimiter(GetAuthRateLimitConfig())
}

// 便捷的中间件函数
// 这些函数提供了简单易用的接口，直接返回对应的限流中间件

// SpeechRateLimitMiddleware 返回语音识别接口的限流中间件
// 适用于语音识别、语音转文字等需要低延迟的接口
func SpeechRateLimitMiddleware() gin.HandlerFunc {
	if speechLimiter == nil {
		// 如果限流器未初始化或限流被禁用，返回空中间件
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return speechLimiter.Middleware()
}

// GeneralRateLimitMiddleware 返回通用API的限流中间件
// 适用于简历管理、面试管理、用户管理等常规业务接口
func GeneralRateLimitMiddleware() gin.HandlerFunc {
	if generalLimiter == nil {
		// 如果限流器未初始化或限流被禁用，返回空中间件
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return generalLimiter.Middleware()
}

// UploadRateLimitMiddleware 返回文件上传接口的限流中间件
// 适用于文件上传、图片上传等资源密集型接口
func UploadRateLimitMiddleware() gin.HandlerFunc {
	if uploadLimiter == nil {
		// 如果限流器未初始化或限流被禁用，返回空中间件
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return uploadLimiter.Middleware()
}

// AuthRateLimitMiddleware 返回认证接口的限流中间件
// 适用于登录、注册、密码重置等敏感接口
func AuthRateLimitMiddleware() gin.HandlerFunc {
	if authLimiter == nil {
		// 如果限流器未初始化或限流被禁用，返回空中间件
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return authLimiter.Middleware()
}

// CustomRateLimitMiddleware 创建自定义限流中间件
// 参数：
//   - config: 自定义的限流配置
//
// 返回值：
//   - gin.HandlerFunc: 自定义的限流中间件
//
// 功能说明：
//
//	允许开发者根据特殊需求创建自定义的限流策略
//	适用于有特殊限流要求的业务场景
func CustomRateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config)
	return limiter.Middleware()
}

// StartRateLimitCleanup 启动限流器清理协程
//
// 功能说明：
//  1. 启动后台协程，定期清理过期的限流器
//  2. 清理间隔为5分钟，防止内存泄漏
//  3. 清理所有全局限流器实例的过期数据
//
// 注意事项：
//
//	这个函数应该在应用启动时调用一次
//	多次调用会创建多个清理协程，造成资源浪费
func StartRateLimitCleanup() {
	go func() {
		// 创建5分钟的定时器
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		// 定期执行清理操作
		for range ticker.C {
			// 只清理已初始化的限流器
			if speechLimiter != nil {
				speechLimiter.Cleanup()
			}
			if generalLimiter != nil {
				generalLimiter.Cleanup()
			}
			if uploadLimiter != nil {
				uploadLimiter.Cleanup()
			}
			if authLimiter != nil {
				authLimiter.Cleanup()
			}
		}
	}()
}
