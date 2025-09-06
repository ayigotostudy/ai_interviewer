package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// BenchmarkRateLimiterPerformance 基准测试限流器性能
func BenchmarkRateLimiterPerformance(b *testing.B) {
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
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow(c, "benchmark_key")
		}
	})
}

// BenchmarkRateLimitMiddlewarePerformance 基准测试限流中间件性能
func BenchmarkRateLimitMiddlewarePerformance(b *testing.B) {
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

// BenchmarkConcurrentAccess 并发访问基准测试
func BenchmarkConcurrentAccess(b *testing.B) {
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 1000
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 2000
		},
		KeyFunc: func(c *gin.Context) string {
			return "concurrent_test"
		},
	}

	limiter := NewRateLimiter(config)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			limiter.Allow(c, "concurrent_test")
		}
	})
}

// BenchmarkRoleBasedLimiting 基于角色的限流基准测试
func BenchmarkRoleBasedLimiting(b *testing.B) {
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
	}

	limiter := NewRateLimiter(config)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Request.RemoteAddr = "192.168.1.1:8080"
			c.Set("id", int64(123))
			c.Set("role", "member")

			limiter.Allow(c, "user:123")
		}
	})
}

// BenchmarkRateLimitMonitor 限流监控基准测试
func BenchmarkRateLimitMonitor(b *testing.B) {
	monitor := NewRateLimitMonitor(1 * time.Minute)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.RecordRequest("test_key", false, 100*time.Millisecond)
		}
	})
}

// TestRateLimiterPerformance 性能测试
func TestRateLimiterPerformance(t *testing.T) {
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 1000
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 2000
		},
		KeyFunc: func(c *gin.Context) string {
			return "performance_test"
		},
	}

	limiter := NewRateLimiter(config)

	// 测试参数
	concurrentUsers := 100
	requestsPerUser := 1000
	totalRequests := concurrentUsers * requestsPerUser

	// 创建测试上下文
	createTestContext := func(userID int64) *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.RemoteAddr = fmt.Sprintf("192.168.1.%d:8080", userID%255)
		c.Set("id", userID)
		c.Set("role", "common")
		return c
	}

	// 性能测试
	start := time.Now()
	var wg sync.WaitGroup
	var successCount, limitedCount, errorCount int64
	var mu sync.Mutex

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int64) {
			defer wg.Done()

			c := createTestContext(userID)
			key := fmt.Sprintf("user:%d", userID)

			for j := 0; j < requestsPerUser; j++ {
				allowed := limiter.Allow(c, key)

				mu.Lock()
				if allowed {
					successCount++
				} else {
					limitedCount++
				}
				mu.Unlock()
			}
		}(int64(i + 1))
	}

	wg.Wait()
	duration := time.Since(start)

	// 计算性能指标
	qps := float64(totalRequests) / duration.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100
	limitRate := float64(limitedCount) / float64(totalRequests) * 100

	t.Logf("性能测试结果:")
	t.Logf("  总请求数: %d", totalRequests)
	t.Logf("  并发用户数: %d", concurrentUsers)
	t.Logf("  每用户请求数: %d", requestsPerUser)
	t.Logf("  总耗时: %v", duration)
	t.Logf("  QPS: %.2f", qps)
	t.Logf("  成功请求: %d (%.2f%%)", successCount, successRate)
	t.Logf("  被限流请求: %d (%.2f%%)", limitedCount, limitRate)
	t.Logf("  错误请求: %d", errorCount)

	// 性能断言
	if qps < 1000 {
		t.Errorf("QPS过低: %.2f, 期望 >= 1000", qps)
	}

	if successRate < 50 {
		t.Errorf("成功率过低: %.2f%%, 期望 >= 50%%", successRate)
	}
}

// TestRateLimiterLatency 延迟测试
func TestRateLimiterLatency(t *testing.T) {
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 100
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 200
		},
		KeyFunc: func(c *gin.Context) string {
			return "latency_test"
		},
	}

	limiter := NewRateLimiter(config)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 测试延迟
	iterations := 1000
	var latencies []time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		limiter.Allow(c, "latency_test")
		latency := time.Since(start)
		latencies = append(latencies, latency)
	}

	// 计算延迟统计
	var total time.Duration
	for _, latency := range latencies {
		total += latency
	}
	avgLatency := total / time.Duration(len(latencies))

	// 排序计算百分位数
	sortLatencies(latencies)
	p50 := latencies[len(latencies)*50/100]
	p90 := latencies[len(latencies)*90/100]
	p95 := latencies[len(latencies)*95/100]
	p99 := latencies[len(latencies)*99/100]

	t.Logf("延迟测试结果:")
	t.Logf("  平均延迟: %v", avgLatency)
	t.Logf("  P50延迟: %v", p50)
	t.Logf("  P90延迟: %v", p90)
	t.Logf("  P95延迟: %v", p95)
	t.Logf("  P99延迟: %v", p99)
	t.Logf("  最大延迟: %v", latencies[len(latencies)-1])
	t.Logf("  最小延迟: %v", latencies[0])

	// 延迟断言
	if avgLatency > 1*time.Millisecond {
		t.Errorf("平均延迟过高: %v, 期望 < 1ms", avgLatency)
	}

	if p99 > 10*time.Millisecond {
		t.Errorf("P99延迟过高: %v, 期望 < 10ms", p99)
	}
}

// TestRateLimiterMemoryUsage 内存使用测试
func TestRateLimiterMemoryUsage(t *testing.T) {
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 100
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 200
		},
		KeyFunc: func(c *gin.Context) string {
			return "memory_test"
		},
	}

	limiter := NewRateLimiter(config)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 创建大量限流器实例
	keys := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
	}

	// 测试内存使用
	start := time.Now()
	for _, key := range keys {
		limiter.Allow(c, key)
	}
	duration := time.Since(start)

	t.Logf("内存使用测试结果:")
	t.Logf("  键数量: %d", len(keys))
	t.Logf("  处理时间: %v", duration)
	t.Logf("  平均处理时间: %v", duration/time.Duration(len(keys)))

	// 验证限流器数量
	if len(limiter.limiters) != len(keys) {
		t.Errorf("限流器数量不匹配: 期望 %d, 实际 %d", len(keys), len(limiter.limiters))
	}
}

// TestRateLimiterCleanup 清理测试
func TestRateLimiterCleanup(t *testing.T) {
	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 100
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 200
		},
		KeyFunc: func(c *gin.Context) string {
			return "cleanup_test"
		},
	}

	limiter := NewRateLimiter(config)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 创建一些限流器
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		limiter.Allow(c, key)
	}

	// 验证限流器已创建
	if len(limiter.limiters) != len(keys) {
		t.Errorf("限流器创建失败: 期望 %d, 实际 %d", len(keys), len(limiter.limiters))
	}

	// 等待清理
	time.Sleep(2 * time.Second)

	// 验证清理后的状态
	// 注意: 实际清理逻辑需要根据具体实现调整
	t.Logf("清理测试完成，当前限流器数量: %d", len(limiter.limiters))
}

// 辅助函数：排序延迟切片
func sortLatencies(latencies []time.Duration) {
	for i := 0; i < len(latencies)-1; i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}
}

// TestRateLimiterStress 压力测试
func TestRateLimiterStress(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试")
	}

	config := RateLimitConfig{
		GetRateFunc: func(c *gin.Context) int {
			return 1000
		},
		GetBurstFunc: func(c *gin.Context) int {
			return 2000
		},
		KeyFunc: func(c *gin.Context) string {
			return "stress_test"
		},
	}

	limiter := NewRateLimiter(config)

	// 压力测试参数
	duration := 30 * time.Second
	concurrent := 100
	requestsPerSecond := 1000

	// 创建测试上下文
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "192.168.1.1:8080"
	c.Set("id", int64(1))
	c.Set("role", "common")

	// 启动压力测试
	start := time.Now()
	var wg sync.WaitGroup
	var successCount, limitedCount int64
	var mu sync.Mutex

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond/concurrent))
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					allowed := limiter.Allow(c, "stress_test")

					mu.Lock()
					if allowed {
						successCount++
					} else {
						limitedCount++
					}
					mu.Unlock()

					if time.Since(start) >= duration {
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	totalDuration := time.Since(start)

	// 计算性能指标
	totalRequests := successCount + limitedCount
	qps := float64(totalRequests) / totalDuration.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100
	limitRate := float64(limitedCount) / float64(totalRequests) * 100

	t.Logf("压力测试结果:")
	t.Logf("  测试时长: %v", totalDuration)
	t.Logf("  并发数: %d", concurrent)
	t.Logf("  总请求数: %d", totalRequests)
	t.Logf("  QPS: %.2f", qps)
	t.Logf("  成功请求: %d (%.2f%%)", successCount, successRate)
	t.Logf("  被限流请求: %d (%.2f%%)", limitedCount, limitRate)

	// 压力测试断言
	if qps < 500 {
		t.Errorf("QPS过低: %.2f, 期望 >= 500", qps)
	}

	if successRate < 30 {
		t.Errorf("成功率过低: %.2f%%, 期望 >= 30%%", successRate)
	}
}
