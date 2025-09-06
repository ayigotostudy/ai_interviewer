// Package middleware 提供限流监控功能
// 用于收集、统计和分析限流相关的性能指标
// 支持实时监控、历史统计和健康状态评估
package middleware

import (
	"sync"
	"time"
)

// RateLimitStats 限流统计信息结构体
// 记录特定限流键的详细统计信息
type RateLimitStats struct {
	// TotalRequests 总请求数
	// 记录该限流键接收到的所有请求总数
	TotalRequests int64 `json:"total_requests"`

	// LimitedRequests 被限流的请求数
	// 记录因超过限流阈值而被拒绝的请求数量
	LimitedRequests int64 `json:"limited_requests"`

	// LimitRate 限流率
	// 计算公式：LimitedRequests / TotalRequests
	// 用于评估限流策略的有效性和系统负载情况
	LimitRate float64 `json:"limit_rate"`

	// ActiveConnections 当前活跃连接数
	// 记录当前正在处理的连接数量
	// 用于监控系统并发负载
	ActiveConnections int64 `json:"active_connections"`

	// AvgResponseTime 平均响应时间，单位：毫秒
	// 记录该限流键下请求的平均处理时间
	// 用于性能分析和优化
	AvgResponseTime float64 `json:"avg_response_time"`

	// TimeWindow 统计时间窗口
	// 定义统计数据的时间范围
	// 超过此时间窗口的统计数据会被清理
	TimeWindow time.Duration `json:"time_window"`

	// LastUpdate 最后更新时间
	// 记录统计数据最后一次更新的时间戳
	// 用于判断数据的新鲜度和清理过期数据
	LastUpdate time.Time `json:"last_update"`
}

// RateLimitMonitor 限流监控器结构体
// 负责收集、存储和分析限流相关的统计数据
type RateLimitMonitor struct {
	// stats 统计数据映射表
	// key: 限流键（如"ip:192.168.1.1"或"user:12345"）
	// value: 对应的统计信息
	stats map[string]*RateLimitStats

	// mu 读写锁，保护stats映射表的并发访问
	// 使用读写锁是因为读操作（查询统计）比写操作（更新统计）更频繁
	mu sync.RWMutex

	// timeWindow 统计时间窗口
	// 定义统计数据的时间范围，超过此时间的统计数据会被清理
	// 例如：5分钟表示只保留最近5分钟的统计数据
	timeWindow time.Duration

	// responseTimes 响应时间历史记录
	// key: 限流键
	// value: 该键下所有请求的响应时间历史记录
	// 用于计算平均响应时间和性能分析
	responseTimes map[string][]time.Duration
}

// NewRateLimitMonitor 创建新的限流监控器实例
// 参数：
//   - timeWindow: 统计时间窗口，超过此时间的统计数据会被清理
//
// 返回值：
//   - *RateLimitMonitor: 初始化好的监控器实例
//
// 功能说明：
//  1. 初始化统计数据映射表和响应时间记录
//  2. 设置统计时间窗口
//  3. 启动后台清理协程，定期清理过期数据
//  4. 返回可用的监控器实例
func NewRateLimitMonitor(timeWindow time.Duration) *RateLimitMonitor {
	monitor := &RateLimitMonitor{
		stats:         make(map[string]*RateLimitStats),
		timeWindow:    timeWindow,
		responseTimes: make(map[string][]time.Duration),
	}

	// 启动后台清理协程
	// 定期清理过期的统计数据，防止内存泄漏
	go monitor.startCleanup()

	return monitor
}

// RecordRequest 记录请求统计信息
// 参数：
//   - key: 限流键，用于标识不同的限流对象
//   - limited: 是否被限流，true表示请求被限流，false表示请求通过
//   - responseTime: 请求响应时间
//
// 功能说明：
//  1. 更新总请求数和限流请求数
//  2. 计算限流率
//  3. 记录响应时间历史
//  4. 计算平均响应时间
//  5. 更新最后修改时间
func (m *RateLimitMonitor) RecordRequest(key string, limited bool, responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取或创建该键的统计信息
	stats, exists := m.stats[key]
	if !exists {
		// 创建新的统计信息记录
		stats = &RateLimitStats{
			TimeWindow: m.timeWindow,
			LastUpdate: time.Now(),
		}
		m.stats[key] = stats
		// 初始化响应时间记录，预分配100个元素的容量
		m.responseTimes[key] = make([]time.Duration, 0, 100)
	}

	// 更新请求计数
	stats.TotalRequests++
	if limited {
		stats.LimitedRequests++
	}

	// 计算限流率
	// 限流率 = 被限流请求数 / 总请求数
	if stats.TotalRequests > 0 {
		stats.LimitRate = float64(stats.LimitedRequests) / float64(stats.TotalRequests)
	}

	// 记录响应时间历史
	responseTimes := m.responseTimes[key]
	responseTimes = append(responseTimes, responseTime)

	// 保持最近100个响应时间记录，避免内存无限增长
	if len(responseTimes) > 100 {
		responseTimes = responseTimes[len(responseTimes)-100:]
	}
	m.responseTimes[key] = responseTimes

	// 计算平均响应时间
	if len(responseTimes) > 0 {
		var total time.Duration
		for _, rt := range responseTimes {
			total += rt
		}
		// 转换为毫秒并计算平均值
		stats.AvgResponseTime = float64(total.Milliseconds()) / float64(len(responseTimes))
	}

	// 更新最后修改时间
	stats.LastUpdate = time.Now()
}

// GetStats 获取统计信息
func (m *RateLimitMonitor) GetStats(key string) *RateLimitStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if stats, exists := m.stats[key]; exists {
		// 返回副本避免并发修改
		return &RateLimitStats{
			TotalRequests:     stats.TotalRequests,
			LimitedRequests:   stats.LimitedRequests,
			LimitRate:         stats.LimitRate,
			ActiveConnections: stats.ActiveConnections,
			AvgResponseTime:   stats.AvgResponseTime,
			TimeWindow:        stats.TimeWindow,
			LastUpdate:        stats.LastUpdate,
		}
	}
	return nil
}

// GetAllStats 获取所有统计信息
func (m *RateLimitMonitor) GetAllStats() map[string]*RateLimitStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*RateLimitStats)
	for key, stats := range m.stats {
		result[key] = &RateLimitStats{
			TotalRequests:     stats.TotalRequests,
			LimitedRequests:   stats.LimitedRequests,
			LimitRate:         stats.LimitRate,
			ActiveConnections: stats.ActiveConnections,
			AvgResponseTime:   stats.AvgResponseTime,
			TimeWindow:        stats.TimeWindow,
			LastUpdate:        stats.LastUpdate,
		}
	}
	return result
}

// GetTopLimitedKeys 获取被限流最多的键
func (m *RateLimitMonitor) GetTopLimitedKeys(limit int) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	type keyStats struct {
		key   string
		stats *RateLimitStats
	}

	var allStats []keyStats
	for key, stats := range m.stats {
		allStats = append(allStats, keyStats{key, stats})
	}

	// 按限流请求数排序
	for i := 0; i < len(allStats)-1; i++ {
		for j := i + 1; j < len(allStats); j++ {
			if allStats[i].stats.LimitedRequests < allStats[j].stats.LimitedRequests {
				allStats[i], allStats[j] = allStats[j], allStats[i]
			}
		}
	}

	// 返回前limit个
	if limit > len(allStats) {
		limit = len(allStats)
	}

	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = allStats[i].key
	}

	return result
}

// startCleanup 启动清理协程
func (m *RateLimitMonitor) startCleanup() {
	ticker := time.NewTicker(m.timeWindow)
	defer ticker.Stop()

	for range ticker.C {
		m.cleanup()
	}
}

// cleanup 清理过期的统计数据
func (m *RateLimitMonitor) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for key, stats := range m.stats {
		// 如果统计数据超过时间窗口，重置
		if now.Sub(stats.LastUpdate) > m.timeWindow {
			stats.TotalRequests = 0
			stats.LimitedRequests = 0
			stats.LimitRate = 0
			stats.AvgResponseTime = 0
			stats.LastUpdate = now

			// 清空响应时间记录
			m.responseTimes[key] = make([]time.Duration, 0, 100)
		}
	}
}

// 全局限流监控器实例
var globalMonitor = NewRateLimitMonitor(5 * time.Minute)

// GetGlobalMonitor 获取全局限流监控器
func GetGlobalMonitor() *RateLimitMonitor {
	return globalMonitor
}

// RecordRequest 记录请求到全局监控器
func RecordRequest(key string, limited bool, responseTime time.Duration) {
	globalMonitor.RecordRequest(key, limited, responseTime)
}

// GetGlobalStats 获取全局统计信息
func GetGlobalStats() map[string]*RateLimitStats {
	return globalMonitor.GetAllStats()
}

// GetTopLimitedKeys 获取被限流最多的键
func GetTopLimitedKeys(limit int) []string {
	return globalMonitor.GetTopLimitedKeys(limit)
}
