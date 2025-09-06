// Package ratelimit 提供限流管理相关的API接口
// 用于监控、统计和管理限流系统的运行状态
// 仅限管理员用户访问，提供系统运维和性能分析功能
package ratelimit

import (
	"ai_jianli_go/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RateLimitController 限流管理控制器
// 提供限流系统的监控、统计和管理功能
// 所有接口都需要管理员权限才能访问
type RateLimitController struct{}

// NewRateLimitController 创建新的限流管理控制器实例
// 返回值：
//   - *RateLimitController: 初始化好的控制器实例
func NewRateLimitController() *RateLimitController {
	return &RateLimitController{}
}

// 辅助函数：检查是否为管理员权限
func isAdmin(ctx *gin.Context) bool {
	role, exists := ctx.Get("role")
	if !exists {
		return false
	}
	roleStr, ok := role.(string)
	if !ok {
		return false
	}
	return roleStr == "super_admin"
}

// GetStats 获取全局限流统计信息
// 接口路径：GET /api/v1/ratelimit/stats
// 权限要求：仅限super_admin角色访问
//
// 功能说明：
//  1. 验证用户权限，确保只有管理员可以访问
//  2. 获取所有限流键的统计信息
//  3. 返回详细的统计数据，包括请求数、限流数、响应时间等
//
// 响应格式：
//
//	{
//	  "data": {
//	    "ip:192.168.1.1": {
//	      "total_requests": 1000,
//	      "limited_requests": 50,
//	      "limit_rate": 0.05,
//	      "avg_response_time": 120.5
//	    }
//	  },
//	  "message": "获取限流统计信息成功"
//	}
func (c *RateLimitController) GetStats(ctx *gin.Context) {
	// 检查管理员权限
	if !isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
			"code":  "INSUFFICIENT_PERMISSION",
		})
		return
	}

	// 获取所有限流键的统计信息
	stats := middleware.GetGlobalStats()

	ctx.JSON(http.StatusOK, gin.H{
		"data":    stats,
		"message": "获取限流统计信息成功",
	})
}

// GetTopLimited 获取被限流最多的键
func (c *RateLimitController) GetTopLimited(ctx *gin.Context) {
	// 检查管理员权限
	if !isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
			"code":  "INSUFFICIENT_PERMISSION",
		})
		return
	}

	// 获取limit参数
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 获取被限流最多的键
	topKeys := middleware.GetTopLimitedKeys(limit)

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"top_limited_keys": topKeys,
			"limit":            limit,
		},
		"message": "获取被限流最多的键成功",
	})
}

// GetKeyStats 获取特定键的统计信息
func (c *RateLimitController) GetKeyStats(ctx *gin.Context) {
	// 检查管理员权限
	if !isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
			"code":  "INSUFFICIENT_PERMISSION",
		})
		return
	}

	// 获取键参数
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "键参数不能为空",
			"code":  "INVALID_PARAMETER",
		})
		return
	}

	// 获取特定键的统计信息
	monitor := middleware.GetGlobalMonitor()
	stats := monitor.GetStats(key)

	if stats == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "未找到该键的统计信息",
			"code":  "KEY_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"key":   key,
			"stats": stats,
		},
		"message": "获取键统计信息成功",
	})
}

// GetHealth 获取限流系统健康状态
func (c *RateLimitController) GetHealth(ctx *gin.Context) {
	// 获取所有统计信息
	stats := middleware.GetGlobalStats()

	// 计算总体健康指标
	var totalRequests int64
	var totalLimited int64
	var avgResponseTime float64
	var responseTimeCount int

	for _, stat := range stats {
		totalRequests += stat.TotalRequests
		totalLimited += stat.LimitedRequests
		if stat.AvgResponseTime > 0 {
			avgResponseTime += stat.AvgResponseTime
			responseTimeCount++
		}
	}

	var overallLimitRate float64
	if totalRequests > 0 {
		overallLimitRate = float64(totalLimited) / float64(totalRequests)
	}

	if responseTimeCount > 0 {
		avgResponseTime = avgResponseTime / float64(responseTimeCount)
	}

	// 健康状态评估
	healthStatus := "healthy"
	if overallLimitRate > 0.1 { // 限流率超过10%认为不健康
		healthStatus = "warning"
	}
	if overallLimitRate > 0.3 { // 限流率超过30%认为严重
		healthStatus = "critical"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status": healthStatus,
			"metrics": gin.H{
				"total_requests":       totalRequests,
				"total_limited":        totalLimited,
				"overall_limit_rate":   overallLimitRate,
				"avg_response_time_ms": avgResponseTime,
				"active_keys":          len(stats),
			},
		},
		"message": "获取限流系统健康状态成功",
	})
}
