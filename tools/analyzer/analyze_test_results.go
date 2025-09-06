package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// 测试结果结构
type AnalysisTestResult struct {
	TotalRequests    int           `json:"total_requests"`
	SuccessRequests  int           `json:"success_requests"`
	LimitedRequests  int           `json:"limited_requests"`
	ErrorRequests    int           `json:"error_requests"`
	TotalDuration    time.Duration `json:"total_duration"`
	QPS              float64       `json:"qps"`
	SuccessRate      float64       `json:"success_rate"`
	LimitRate        float64       `json:"limit_rate"`
	AverageLatency   time.Duration `json:"average_latency"`
	P50Latency       time.Duration `json:"p50_latency"`
	P90Latency       time.Duration `json:"p90_latency"`
	P95Latency       time.Duration `json:"p95_latency"`
	P99Latency       time.Duration `json:"p99_latency"`
	MaxLatency       time.Duration `json:"max_latency"`
	MinLatency       time.Duration `json:"min_latency"`
}

// 性能等级
type PerformanceLevel struct {
	Level       string
	QPS         float64
	SuccessRate float64
	AvgLatency  time.Duration
	Description string
}

// 性能等级定义
var performanceLevels = []PerformanceLevel{
	{"优秀", 1000, 90, 1 * time.Millisecond, "系统性能优秀，可以处理高并发请求"},
	{"良好", 500, 80, 5 * time.Millisecond, "系统性能良好，可以处理中等并发请求"},
	{"一般", 200, 70, 10 * time.Millisecond, "系统性能一般，需要优化"},
	{"较差", 100, 60, 20 * time.Millisecond, "系统性能较差，需要重点优化"},
	{"很差", 50, 50, 50 * time.Millisecond, "系统性能很差，需要全面优化"},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用方法: go run analyze_test_results.go <测试结果文件>")
		fmt.Println("示例: go run analyze_test_results.go rate_limit_test_result.json")
		os.Exit(1)
	}

	filename := os.Args[1]
	result, err := loadTestResult(filename)
	if err != nil {
		fmt.Printf("加载测试结果失败: %v\n", err)
		os.Exit(1)
	}

	analyzeResult(result)
}

// 加载测试结果
func loadTestResult(filename string) (*AnalysisTestResult, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result AnalysisTestResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// 分析测试结果
func analyzeResult(result *AnalysisTestResult) {
	fmt.Println("=== 限流性能测试结果分析 ===")
	fmt.Println()

	// 基本信息
	fmt.Println("📊 基本信息:")
	fmt.Printf("  总请求数: %d\n", result.TotalRequests)
	fmt.Printf("  成功请求: %d (%.2f%%)\n", result.SuccessRequests, result.SuccessRate)
	fmt.Printf("  被限流请求: %d (%.2f%%)\n", result.LimitedRequests, result.LimitRate)
	fmt.Printf("  错误请求: %d\n", result.ErrorRequests)
	fmt.Printf("  总耗时: %v\n", result.TotalDuration)
	fmt.Println()

	// 性能指标
	fmt.Println("⚡ 性能指标:")
	fmt.Printf("  QPS: %.2f\n", result.QPS)
	fmt.Printf("  平均延迟: %v\n", result.AverageLatency)
	fmt.Printf("  P50延迟: %v\n", result.P50Latency)
	fmt.Printf("  P90延迟: %v\n", result.P90Latency)
	fmt.Printf("  P95延迟: %v\n", result.P95Latency)
	fmt.Printf("  P99延迟: %v\n", result.P99Latency)
	fmt.Printf("  最大延迟: %v\n", result.MaxLatency)
	fmt.Printf("  最小延迟: %v\n", result.MinLatency)
	fmt.Println()

	// 性能等级评估
	level := evaluatePerformance(result)
	fmt.Printf("🎯 性能等级: %s\n", level.Level)
	fmt.Printf("   描述: %s\n", level.Description)
	fmt.Println()

	// 详细分析
	analyzeDetails(result)
	fmt.Println()

	// 优化建议
	provideRecommendations(result)
	fmt.Println()

	// 生成报告
	generateReport(result, level)
}

// 评估性能等级
func evaluatePerformance(result *AnalysisTestResult) PerformanceLevel {
	for _, level := range performanceLevels {
		if result.QPS >= level.QPS && result.SuccessRate >= level.SuccessRate && result.AverageLatency <= level.AvgLatency {
			return level
		}
	}
	return performanceLevels[len(performanceLevels)-1] // 返回最差的等级
}

// 详细分析
func analyzeDetails(result *AnalysisTestResult) {
	fmt.Println("🔍 详细分析:")

	// QPS分析
	if result.QPS >= 1000 {
		fmt.Println("  ✅ QPS表现优秀，系统可以处理高并发请求")
	} else if result.QPS >= 500 {
		fmt.Println("  ✅ QPS表现良好，系统可以处理中等并发请求")
	} else if result.QPS >= 200 {
		fmt.Println("  ⚠️  QPS表现一般，建议优化系统性能")
	} else {
		fmt.Println("  ❌ QPS表现较差，需要重点优化")
	}

	// 成功率分析
	if result.SuccessRate >= 90 {
		fmt.Println("  ✅ 成功率表现优秀，系统稳定性很好")
	} else if result.SuccessRate >= 80 {
		fmt.Println("  ✅ 成功率表现良好，系统稳定性较好")
	} else if result.SuccessRate >= 70 {
		fmt.Println("  ⚠️  成功率表现一般，建议检查系统稳定性")
	} else {
		fmt.Println("  ❌ 成功率表现较差，需要检查系统问题")
	}

	// 延迟分析
	if result.AverageLatency <= 1*time.Millisecond {
		fmt.Println("  ✅ 延迟表现优秀，响应速度很快")
	} else if result.AverageLatency <= 5*time.Millisecond {
		fmt.Println("  ✅ 延迟表现良好，响应速度较快")
	} else if result.AverageLatency <= 10*time.Millisecond {
		fmt.Println("  ⚠️  延迟表现一般，建议优化响应速度")
	} else {
		fmt.Println("  ❌ 延迟表现较差，需要优化响应速度")
	}

	// 限流率分析
	if result.LimitRate <= 10 {
		fmt.Println("  ✅ 限流率适中，系统负载均衡")
	} else if result.LimitRate <= 30 {
		fmt.Println("  ⚠️  限流率较高，建议调整限流参数")
	} else {
		fmt.Println("  ❌ 限流率过高，需要重新配置限流策略")
	}
}

// 提供优化建议
func provideRecommendations(result *AnalysisTestResult) {
	fmt.Println("💡 优化建议:")

	// QPS优化建议
	if result.QPS < 1000 {
		fmt.Println("  🚀 提升QPS:")
		fmt.Println("    - 优化代码实现，减少不必要的计算")
		fmt.Println("    - 使用连接池和缓存")
		fmt.Println("    - 考虑使用Redis集群")
		fmt.Println("    - 调整限流参数，增加rate值")
	}

	// 延迟优化建议
	if result.AverageLatency > 5*time.Millisecond {
		fmt.Println("  ⚡ 降低延迟:")
		fmt.Println("    - 优化数据库查询")
		fmt.Println("    - 使用异步处理")
		fmt.Println("    - 减少网络调用")
		fmt.Println("    - 优化算法复杂度")
	}

	// 稳定性优化建议
	if result.SuccessRate < 80 {
		fmt.Println("  🛡️  提升稳定性:")
		fmt.Println("    - 增加错误处理机制")
		fmt.Println("    - 实现重试机制")
		fmt.Println("    - 监控系统资源使用")
		fmt.Println("    - 优化限流策略")
	}

	// 限流策略优化建议
	if result.LimitRate > 30 {
		fmt.Println("  ⚖️  优化限流策略:")
		fmt.Println("    - 调整rate和burst参数")
		fmt.Println("    - 实现动态限流")
		fmt.Println("    - 基于用户角色的差异化限流")
		fmt.Println("    - 监控限流效果")
	}

	// 通用建议
	fmt.Println("  🔧 通用建议:")
	fmt.Println("    - 定期进行性能测试")
	fmt.Println("    - 监控生产环境指标")
	fmt.Println("    - 建立性能基线")
	fmt.Println("    - 持续优化系统架构")
}

// 生成报告
func generateReport(result *AnalysisTestResult, level PerformanceLevel) {
	report := fmt.Sprintf(`
# 限流性能测试报告

## 测试结果概览
- 测试时间: %s
- 总请求数: %d
- 成功请求: %d (%.2f%%)
- 被限流请求: %d (%.2f%%)
- 错误请求: %d
- 总耗时: %v

## 性能指标
- QPS: %.2f
- 平均延迟: %v
- P50延迟: %v
- P90延迟: %v
- P95延迟: %v
- P99延迟: %v
- 最大延迟: %v
- 最小延迟: %v

## 性能等级
- 等级: %s
- 描述: %s

## 优化建议
1. 根据测试结果调整限流参数
2. 监控生产环境的限流效果
3. 定期进行性能测试
4. 考虑使用Redis集群提高性能

## 下次测试建议
- 测试时间: %s
- 测试场景: 增加更多并发用户
- 监控指标: 重点关注QPS和延迟
`, 
		time.Now().Format("2006-01-02 15:04:05"),
		result.TotalRequests,
		result.SuccessRequests, result.SuccessRate,
		result.LimitedRequests, result.LimitRate,
		result.ErrorRequests,
		result.TotalDuration,
		result.QPS,
		result.AverageLatency,
		result.P50Latency,
		result.P90Latency,
		result.P95Latency,
		result.P99Latency,
		result.MaxLatency,
		result.MinLatency,
		level.Level,
		level.Description,
		time.Now().AddDate(0, 0, 7).Format("2006-01-02 15:04:05"),
	)

	filename := fmt.Sprintf("rate_limit_analysis_report_%s.md", time.Now().Format("20060102_150405"))
	err := ioutil.WriteFile(filename, []byte(report), 0644)
	if err != nil {
		fmt.Printf("生成报告失败: %v\n", err)
		return
	}

	fmt.Printf("📄 详细报告已生成: %s\n", filename)
}
