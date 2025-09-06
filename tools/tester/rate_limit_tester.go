package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

// 测试配置
type TestConfig struct {
	BaseURL        string `json:"base_url"`
	Endpoint       string `json:"endpoint"`
	Concurrent     int    `json:"concurrent"`
	TotalRequests  int    `json:"total_requests"`
	TestDuration   int    `json:"test_duration"` // 秒
	UserID         string `json:"user_id"`
	UserRole       string `json:"user_role"`
	AudioFile      string `json:"audio_file"`
}

// 测试结果
type RateLimitTestResult struct {
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

// 请求统计
type RequestStats struct {
	StatusCode int
	Latency    time.Duration
	Success    bool
	Limited    bool
	Error      bool
}

func main() {
	// 默认配置
	config := TestConfig{
		BaseURL:       "http://localhost:8080",
		Endpoint:      "/api/v1/speech/recognize",
		Concurrent:    10,
		TotalRequests: 1000,
		TestDuration:  60,
		UserID:        "test_user",
		UserRole:      "common",
		AudioFile:     "/tmp/test_audio.wav",
	}

	// 解析命令行参数
	if len(os.Args) > 1 {
		configFile := os.Args[1]
		if err := loadConfig(&config, configFile); err != nil {
			fmt.Printf("加载配置文件失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 创建测试音频文件
	if err := createTestAudioFile(config.AudioFile); err != nil {
		fmt.Printf("创建测试音频文件失败: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(config.AudioFile)

	// 执行测试
	result := runTest(config)

	// 输出结果
	printResults(result)

	// 保存结果到文件
	saveResults(result, "rate_limit_test_result.json")
}

// 加载配置文件
func loadConfig(config *TestConfig, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(config)
}

// 创建测试音频文件
func createTestAudioFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入一些测试数据
	_, err = file.Write(make([]byte, 1024))
	return err
}

// 执行测试
func runTest(config TestConfig) RateLimitTestResult {
	fmt.Printf("开始限流性能测试...\n")
	fmt.Printf("配置: %+v\n", config)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var stats []RequestStats

	start := time.Now()
	requestChan := make(chan int, config.TotalRequests)

	// 启动并发请求
	for i := 0; i < config.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestChan {
				stat := makeRequest(config)
				mu.Lock()
				stats = append(stats, stat)
				mu.Unlock()
			}
		}()
	}

	// 发送请求
	go func() {
		defer close(requestChan)
		for i := 0; i < config.TotalRequests; i++ {
			requestChan <- i
		}
	}()

	wg.Wait()
	totalDuration := time.Since(start)

	// 计算统计结果
	return calculateStats(stats, totalDuration)
}

// 发送单个请求
func makeRequest(config TestConfig) RequestStats {
	start := time.Now()

	// 创建multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加音频文件
	file, err := os.Open(config.AudioFile)
	if err != nil {
		return RequestStats{StatusCode: 0, Latency: time.Since(start), Error: true}
	}
	defer file.Close()

	part, err := writer.CreateFormFile("audio", "test_audio.wav")
	if err != nil {
		return RequestStats{StatusCode: 0, Latency: time.Since(start), Error: true}
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return RequestStats{StatusCode: 0, Latency: time.Since(start), Error: true}
	}

	writer.Close()

	// 创建HTTP请求
	req, err := http.NewRequest("POST", config.BaseURL+config.Endpoint, &buf)
	if err != nil {
		return RequestStats{StatusCode: 0, Latency: time.Since(start), Error: true}
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-User-ID", config.UserID)
	req.Header.Set("X-User-Role", config.UserRole)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return RequestStats{StatusCode: 0, Latency: time.Since(start), Error: true}
	}
	defer resp.Body.Close()

	latency := time.Since(start)
	statusCode := resp.StatusCode

	// 判断请求结果
	success := statusCode == 200
	limited := statusCode == 429
	error := statusCode >= 400 && statusCode != 429

	return RequestStats{
		StatusCode: statusCode,
		Latency:    latency,
		Success:    success,
		Limited:    limited,
		Error:      error,
	}
}

// 计算统计结果
func calculateStats(stats []RequestStats, totalDuration time.Duration) RateLimitTestResult {
	if len(stats) == 0 {
		return RateLimitTestResult{}
	}

	var successCount, limitedCount, errorCount int
	var latencies []time.Duration

	for _, stat := range stats {
		if stat.Success {
			successCount++
		} else if stat.Limited {
			limitedCount++
		} else if stat.Error {
			errorCount++
		}

		if stat.Success || stat.Limited {
			latencies = append(latencies, stat.Latency)
		}
	}

	totalRequests := len(stats)
	qps := float64(totalRequests) / totalDuration.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100
	limitRate := float64(limitedCount) / float64(totalRequests) * 100

	// 计算延迟统计
	var avgLatency, p50Latency, p90Latency, p95Latency, p99Latency, maxLatency, minLatency time.Duration
	if len(latencies) > 0 {
		// 排序
		for i := 0; i < len(latencies)-1; i++ {
			for j := i + 1; j < len(latencies); j++ {
				if latencies[i] > latencies[j] {
					latencies[i], latencies[j] = latencies[j], latencies[i]
				}
			}
		}

		// 计算平均值
		var total time.Duration
		for _, latency := range latencies {
			total += latency
		}
		avgLatency = total / time.Duration(len(latencies))

		// 计算百分位数
		p50Latency = latencies[len(latencies)*50/100]
		p90Latency = latencies[len(latencies)*90/100]
		p95Latency = latencies[len(latencies)*95/100]
		p99Latency = latencies[len(latencies)*99/100]
		maxLatency = latencies[len(latencies)-1]
		minLatency = latencies[0]
	}

	return RateLimitTestResult{
		TotalRequests:   totalRequests,
		SuccessRequests: successCount,
		LimitedRequests: limitedCount,
		ErrorRequests:   errorCount,
		TotalDuration:   totalDuration,
		QPS:             qps,
		SuccessRate:     successRate,
		LimitRate:       limitRate,
		AverageLatency:  avgLatency,
		P50Latency:      p50Latency,
		P90Latency:      p90Latency,
		P95Latency:      p95Latency,
		P99Latency:      p99Latency,
		MaxLatency:      maxLatency,
		MinLatency:      minLatency,
	}
}

// 打印结果
func printResults(result RateLimitTestResult) {
	fmt.Printf("\n=== 限流性能测试结果 ===\n")
	fmt.Printf("总请求数: %d\n", result.TotalRequests)
	fmt.Printf("成功请求: %d (%.2f%%)\n", result.SuccessRequests, result.SuccessRate)
	fmt.Printf("被限流请求: %d (%.2f%%)\n", result.LimitedRequests, result.LimitRate)
	fmt.Printf("错误请求: %d\n", result.ErrorRequests)
	fmt.Printf("总耗时: %v\n", result.TotalDuration)
	fmt.Printf("QPS: %.2f\n", result.QPS)
	fmt.Printf("\n延迟统计:\n")
	fmt.Printf("  平均延迟: %v\n", result.AverageLatency)
	fmt.Printf("  P50延迟: %v\n", result.P50Latency)
	fmt.Printf("  P90延迟: %v\n", result.P90Latency)
	fmt.Printf("  P95延迟: %v\n", result.P95Latency)
	fmt.Printf("  P99延迟: %v\n", result.P99Latency)
	fmt.Printf("  最大延迟: %v\n", result.MaxLatency)
	fmt.Printf("  最小延迟: %v\n", result.MinLatency)
}

// 保存结果到文件
func saveResults(result RateLimitTestResult, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("保存结果失败: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Printf("编码结果失败: %v\n", err)
		return
	}

	fmt.Printf("\n结果已保存到: %s\n", filename)
}
