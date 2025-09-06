# AI简历助手限流性能测试 Makefile

.PHONY: help test-rate-limit test-benchmark test-performance test-stress test-all clean

# 默认目标
help:
	@echo "AI简历助手限流性能测试"
	@echo ""
	@echo "可用命令:"
	@echo "  quick-test          - 快速限流测试"
	@echo "  test-rate-limit     - 运行基础限流测试"
	@echo "  test-benchmark      - 运行基准测试"
	@echo "  test-performance    - 运行性能测试"
	@echo "  test-stress         - 运行压力测试"
	@echo "  test-custom         - 运行自定义测试工具"
	@echo "  analyze-results     - 分析测试结果"
	@echo "  test-all            - 运行所有测试"
	@echo "  clean               - 清理测试文件"
	@echo ""
	@echo "示例:"
	@echo "  make test-rate-limit"
	@echo "  make test-benchmark"
	@echo "  make test-performance"

# 快速测试
quick-test:
	@echo "运行快速限流测试..."
	@chmod +x scripts/quick_test.sh
	@./scripts/quick_test.sh

# 基础限流测试
test-rate-limit:
	@echo "运行基础限流测试..."
	@chmod +x scripts/rate_limit_test.sh
	@./scripts/rate_limit_test.sh

# 限流性能测试
test-performance:
	@echo "运行限流性能测试..."
	@chmod +x scripts/rate_limit_performance_test.sh
	@./scripts/rate_limit_performance_test.sh

# 基准测试
test-benchmark:
	@echo "运行基准测试..."
	@go test -bench=. -benchmem -run=^$$ ./internal/middleware/

# 性能测试
test-performance-go:
	@echo "运行Go性能测试..."
	@go test -run=TestRateLimiterPerformance -v ./internal/middleware/

# 延迟测试
test-latency:
	@echo "运行延迟测试..."
	@go test -run=TestRateLimiterLatency -v ./internal/middleware/

# 内存测试
test-memory:
	@echo "运行内存使用测试..."
	@go test -run=TestRateLimiterMemoryUsage -v ./internal/middleware/

# 压力测试
test-stress:
	@echo "运行压力测试..."
	@go test -run=TestRateLimiterStress -v -timeout=5m ./internal/middleware/

# 并发测试
test-concurrent:
	@echo "运行并发测试..."
	@go test -run=TestConcurrentAccess -v ./internal/middleware/

# 角色测试
test-role:
	@echo "运行角色差异化限流测试..."
	@go test -run=TestSpeechRateLimit -v ./internal/middleware/

# 监控测试
test-monitor:
	@echo "运行限流监控测试..."
	@go test -run=TestRateLimitMonitor -v ./internal/middleware/

# 运行所有测试
test-all: test-rate-limit test-benchmark test-performance-go test-latency test-memory test-concurrent test-role test-monitor
	@echo "所有测试完成"

# 运行自定义测试工具
test-custom:
	@echo "运行自定义测试工具..."
	@go run tools/tester/rate_limit_tester.go tools/rate_limit_test_config.json

# 运行自定义测试工具（使用默认配置）
test-custom-default:
	@echo "运行自定义测试工具（默认配置）..."
	@go run tools/tester/rate_limit_tester.go

# 生成测试报告
test-report:
	@echo "生成测试报告..."
	@go test -v ./internal/middleware/ > test_report.txt 2>&1
	@echo "测试报告已生成: test_report.txt"

# 分析测试结果
analyze-results:
	@echo "分析测试结果..."
	@if [ -f "rate_limit_test_result.json" ]; then \
		go run tools/analyzer/analyze_test_results.go rate_limit_test_result.json; \
	else \
		echo "未找到测试结果文件: rate_limit_test_result.json"; \
		echo "请先运行测试: make test-custom"; \
	fi

# 清理测试文件
clean:
	@echo "清理测试文件..."
	@rm -f /tmp/test_audio.wav
	@rm -f /tmp/rate_limit_results.json
	@rm -f /tmp/resource_monitor.log
	@rm -f /tmp/rate_limit_performance_report_*.txt
	@rm -f test_report.txt
	@rm -f rate_limit_test_result.json
	@echo "清理完成"

# 安装依赖
install-deps:
	@echo "安装测试依赖..."
	@go mod tidy
	@go mod download

# 检查服务器状态
check-server:
	@echo "检查服务器状态..."
	@curl -s http://localhost:8080/health > /dev/null && echo "服务器运行正常" || echo "服务器未运行"

# 启动服务器（开发模式）
start-server:
	@echo "启动服务器..."
	@go run main.go

# 停止服务器
stop-server:
	@echo "停止服务器..."
	@pkill -f "go run main.go" || echo "服务器未运行"

# 重启服务器
restart-server: stop-server start-server

# 查看限流统计
stats:
	@echo "查看限流统计..."
	@curl -s -H "X-User-Role: super_admin" http://localhost:8080/api/v1/ratelimit/stats | jq '.' || echo "需要管理员权限或jq工具"

# 查看被限流最多的键
top-limited:
	@echo "查看被限流最多的键..."
	@curl -s -H "X-User-Role: super_admin" http://localhost:8080/api/v1/ratelimit/top-limited?limit=10 | jq '.' || echo "需要管理员权限或jq工具"

# 查看限流健康状态
health:
	@echo "查看限流健康状态..."
	@curl -s http://localhost:8080/api/v1/ratelimit/health | jq '.' || echo "需要jq工具"

# 运行完整测试套件
full-test: check-server test-all test-report
	@echo "完整测试套件执行完成"

# 快速测试（仅基础功能）
quick-test-full: check-server test-rate-limit test-benchmark
	@echo "快速测试完成"
