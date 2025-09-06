package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试配置
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 10 // 10 QPS
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 20 // 20个瞬时并发
		},
		KeyFunc: func(c *gin.Context) string {
			return "test_key"
		},
	}

	// 创建限流器
	limiter := NewRateLimiter(config)

	// 创建测试路由
	r := gin.New()
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 测试正常请求
	t.Run("正常请求应该通过", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("X-RateLimit-Limit"), "10")
	})

	// 测试限流
	t.Run("超过限制的请求应该被限流", func(t *testing.T) {
		// 快速发送大量请求
		limitedCount := 0
		for i := 0; i < 25; i++ { // 超过burst限制
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			if w.Code == http.StatusTooManyRequests {
				limitedCount++
			}
		}

		// 应该有一些请求被限流
		assert.Greater(t, limitedCount, 0, "应该有请求被限流")
	})

	// 测试跳过功能
	t.Run("跳过限流的请求应该通过", func(t *testing.T) {
		skipConfig := RateLimitConfig{
			GetRateFunc: func(c *gin.Context) int {
				return 1 // 很低的限制
			},
			GetBurstFunc: func(c *gin.Context) int {
				return 1
			},
			KeyFunc: func(c *gin.Context) string {
				return "test_key"
			},
			SkipFunc: func(c *gin.Context) bool {
				return c.GetHeader("X-Skip-RateLimit") == "true"
			},
		}

		skipLimiter := NewRateLimiter(skipConfig)
		skipRouter := gin.New()
		skipRouter.Use(skipLimiter.Middleware())
		skipRouter.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// 发送跳过限流的请求
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Skip-RateLimit", "true")
		skipRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRateLimitMonitor(t *testing.T) {
	// 创建监控器
	monitor := NewRateLimitMonitor(1 * time.Minute)

	// 记录一些请求
	monitor.RecordRequest("key1", false, 100*time.Millisecond)
	monitor.RecordRequest("key1", true, 50*time.Millisecond)
	monitor.RecordRequest("key2", false, 200*time.Millisecond)

	// 获取统计信息
	stats := monitor.GetStats("key1")
	assert.NotNil(t, stats)
	assert.Equal(t, int64(2), stats.TotalRequests)
	assert.Equal(t, int64(1), stats.LimitedRequests)
	assert.Equal(t, 0.5, stats.LimitRate)

	// 获取所有统计信息
	allStats := monitor.GetAllStats()
	assert.Len(t, allStats, 2)
	assert.Contains(t, allStats, "key1")
	assert.Contains(t, allStats, "key2")

	// 测试获取被限流最多的键
	topKeys := monitor.GetTopLimitedKeys(1)
	assert.Len(t, topKeys, 1)
	assert.Equal(t, "key1", topKeys[0])
}

func TestSpeechRateLimit(t *testing.T) {
	// 测试语音识别限流配置
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			switch role {
			case "member":
				return 200
			case "super_member":
				return 500
			case "super_admin":
				return 1000
			default:
				return 50
			}
		},
		GetBurstFunc: func(c *gin.Context) int {
			role := getUserRole(c)
			switch role {
			case "member":
				return 400
			case "super_member":
				return 1000
			case "super_admin":
				return 2000
			default:
				return 100
			}
		},
		KeyFunc: func(c *gin.Context) string {
			userID := getUserID(c)
			if userID > 0 {
				return fmt.Sprintf("user:%d", userID)
			}
			return fmt.Sprintf("ip:%s", c.ClientIP())
		},
		SkipFunc: func(c *gin.Context) bool {
			role := getUserRole(c)
			skipRoles := []string{"super_admin", "member"}
			for _, skipRole := range skipRoles {
				if role == skipRole {
					return true
				}
			}
			return false
		},
	}

	assert.NotNil(t, config.GetRateFunc)
	assert.NotNil(t, config.GetBurstFunc)
	assert.NotNil(t, config.KeyFunc)
	assert.NotNil(t, config.SkipFunc)

	// 测试键生成函数
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "192.168.1.1:8080"

	// 测试IP限流（没有用户ID时）
	key := config.KeyFunc(c)
	assert.Contains(t, key, "192.168.1.1")

	// 测试用户限流
	c.Set("id", int64(123))
	key = config.KeyFunc(c)
	assert.Equal(t, "user:123", key)

	// 测试动态限流参数
	// 普通用户
	rate := config.GetRateFunc(c)
	burst := config.GetBurstFunc(c)
	assert.Equal(t, 50, rate)
	assert.Equal(t, 100, burst)

	// 会员用户
	c.Set("role", "member")
	rate = config.GetRateFunc(c)
	burst = config.GetBurstFunc(c)
	assert.Equal(t, 200, rate)
	assert.Equal(t, 400, burst)

	// 超级管理员
	c.Set("role", "super_admin")
	rate = config.GetRateFunc(c)
	burst = config.GetBurstFunc(c)
	assert.Equal(t, 1000, rate)
	assert.Equal(t, 2000, burst)

	// 测试跳过功能
	shouldSkip := config.SkipFunc(c)
	assert.True(t, shouldSkip)
}

func TestConcurrentAccess(t *testing.T) {
	// 测试并发访问
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 100
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 200
		},
		KeyFunc: func(c *gin.Context) string {
			return "concurrent_test"
		},
	}

	limiter := NewRateLimiter(config)

	// 并发测试
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			allowed := limiter.Allow(c, "concurrent_test")
			// 大部分请求应该被允许
			assert.True(t, allowed)
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 100; i++ {
		<-done
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 1000
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 2000
		},
		KeyFunc: func(c *gin.Context) string {
			return "benchmark_key"
		},
	}

	limiter := NewRateLimiter(config)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			limiter.Allow(c, "benchmark_key")
		}
	})
}

func BenchmarkRateLimitMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)

	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 1000
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 2000
		},
		KeyFunc: func(c *gin.Context) string {
			return "benchmark_key"
		},
	}

	limiter := NewRateLimiter(config)

	r := gin.New()
	r.Use(limiter.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)
		}
	})
}

func TestGetUserID(t *testing.T) {
	// 创建测试用的Gin上下文
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// 测试没有设置用户ID的情况
	userID := getUserID(c)
	assert.Equal(t, int64(0), userID, "没有设置用户ID时应该返回0")

	// 测试设置用户ID的情况
	c.Set("id", int64(123))
	userID = getUserID(c)
	assert.Equal(t, int64(123), userID, "设置用户ID后应该返回正确的值")

	// 测试设置错误类型的情况
	c.Set("id", "invalid")
	userID = getUserID(c)
	assert.Equal(t, int64(0), userID, "设置错误类型的用户ID时应该返回0")
}

func TestGetUserRole(t *testing.T) {
	// 创建测试用的Gin上下文
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// 测试没有设置角色的情况
	role := getUserRole(c)
	assert.Equal(t, "", role, "没有设置角色时应该返回空字符串")

	// 测试设置角色的情况
	c.Set("role", "admin")
	role = getUserRole(c)
	assert.Equal(t, "admin", role, "设置角色后应该返回正确的值")

	// 测试设置错误类型的情况
	c.Set("role", 123)
	role = getUserRole(c)
	assert.Equal(t, "", role, "设置错误类型的角色时应该返回空字符串")
}
