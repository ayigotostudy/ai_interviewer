#!/bin/bash

# 限流性能测试脚本
# 用于全面测试限流器的性能表现

# 配置
BASE_URL="http://localhost:8080"
SPEECH_ENDPOINT="/api/v1/speech/recognize"
USER_ENDPOINT="/api/v1/user/login"
RATE_LIMIT_ENDPOINT="/api/v1/ratelimit"
RESUME_ENDPOINT="/api/v1/resume"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_REQUESTS=0
SUCCESS_REQUESTS=0
RATE_LIMITED_REQUESTS=0
ERROR_REQUESTS=0
TOTAL_TIME=0

# 创建测试音频文件
create_test_audio() {
    echo -e "${YELLOW}创建测试音频文件...${NC}"
    dd if=/dev/zero of=/tmp/test_audio.wav bs=1024 count=10 2>/dev/null
    echo -e "${GREEN}测试音频文件创建完成${NC}"
}

# 清理测试文件
cleanup() {
    echo -e "${YELLOW}清理测试文件...${NC}"
    rm -f /tmp/test_audio.wav
    rm -f /tmp/rate_limit_results.json
    echo -e "${GREEN}清理完成${NC}"
}

# 检查服务器状态
check_server() {
    echo -e "${YELLOW}检查服务器状态...${NC}"
    if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
        echo -e "${RED}错误: 服务器未运行，请先启动服务器${NC}"
        exit 1
    fi
    echo -e "${GREEN}服务器运行正常${NC}"
}

# 测试1: 基础限流功能测试
test_basic_rate_limit() {
    echo -e "${BLUE}=== 测试1: 基础限流功能测试 ===${NC}"
    
    local test_name="基础限流"
    local requests=50
    local concurrent=10
    
    echo "发送 $requests 个请求，并发数: $concurrent"
    
    local start_time=$(date +%s.%N)
    local success=0
    local limited=0
    local error=0
    
    for i in $(seq 1 $requests); do
        (
            response=$(curl -s -w "%{http_code}" -o /dev/null \
                -F "audio=@/tmp/test_audio.wav" \
                -H "X-User-ID: test_user_$i" \
                -H "X-User-Role: common" \
                "$BASE_URL$SPEECH_ENDPOINT" 2>/dev/null)
            
            case $response in
                200) ((success++)) ;;
                429) ((limited++)) ;;
                *) ((error++)) ;;
            esac
        ) &
        
        # 控制并发数
        if [ $((i % concurrent)) -eq 0 ]; then
            wait
        fi
    done
    
    wait
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    echo -e "${GREEN}测试完成:${NC}"
    echo "  成功请求: $success"
    echo "  被限流请求: $limited"
    echo "  错误请求: $error"
    echo "  总耗时: ${duration}秒"
    echo "  QPS: $(echo "scale=2; $requests / $duration" | bc)"
    echo "  限流率: $(echo "scale=2; $limited * 100 / $requests" | bc)%"
    echo ""
}

# 测试2: 高并发压力测试
test_high_concurrency() {
    echo -e "${BLUE}=== 测试2: 高并发压力测试 ===${NC}"
    
    local test_name="高并发压力"
    local requests=1000
    local concurrent=50
    
    echo "发送 $requests 个请求，并发数: $concurrent"
    
    local start_time=$(date +%s.%N)
    local success=0
    local limited=0
    local error=0
    
    # 使用xargs控制并发
    seq 1 $requests | xargs -n 1 -P $concurrent -I {} bash -c '
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -F "audio=@/tmp/test_audio.wav" \
            -H "X-User-ID: test_user_{}" \
            -H "X-User-Role: common" \
            "'$BASE_URL$SPEECH_ENDPOINT'" 2>/dev/null)
        
        case $response in
            200) echo "SUCCESS" ;;
            429) echo "LIMITED" ;;
            *) echo "ERROR" ;;
        esac
    ' | while read result; do
        case $result in
            SUCCESS) ((success++)) ;;
            LIMITED) ((limited++)) ;;
            ERROR) ((error++)) ;;
        esac
    done
    
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    echo -e "${GREEN}压力测试完成:${NC}"
    echo "  成功请求: $success"
    echo "  被限流请求: $limited"
    echo "  错误请求: $error"
    echo "  总耗时: ${duration}秒"
    echo "  QPS: $(echo "scale=2; $requests / $duration" | bc)"
    echo "  限流率: $(echo "scale=2; $limited * 100 / $requests" | bc)%"
    echo ""
}

# 测试3: 延迟测试
test_latency() {
    echo -e "${BLUE}=== 测试3: 延迟测试 ===${NC}"
    
    local requests=100
    local latencies=()
    
    echo "测试 $requests 个请求的延迟..."
    
    for i in $(seq 1 $requests); do
        start_time=$(date +%s.%N)
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -F "audio=@/tmp/test_audio.wav" \
            -H "X-User-ID: test_user_$i" \
            -H "X-User-Role: common" \
            "$BASE_URL$SPEECH_ENDPOINT" 2>/dev/null)
        end_time=$(date +%s.%N)
        
        if [ "$response" = "200" ]; then
            latency=$(echo "($end_time - $start_time) * 1000" | bc)
            latencies+=($latency)
        fi
        
        # 避免过于频繁的请求
        sleep 0.1
    done
    
    # 计算延迟统计
    if [ ${#latencies[@]} -gt 0 ]; then
        # 排序
        IFS=$'\n' sorted=($(sort -n <<<"${latencies[*]}"))
        unset IFS
        
        local count=${#sorted[@]}
        local p50_idx=$((count / 2))
        local p90_idx=$((count * 9 / 10))
        local p95_idx=$((count * 95 / 100))
        local p99_idx=$((count * 99 / 100))
        
        local sum=0
        for latency in "${sorted[@]}"; do
            sum=$(echo "$sum + $latency" | bc)
        done
        local avg=$(echo "scale=2; $sum / $count" | bc)
        
        echo -e "${GREEN}延迟统计:${NC}"
        echo "  平均延迟: ${avg}ms"
        echo "  P50延迟: ${sorted[$p50_idx]}ms"
        echo "  P90延迟: ${sorted[$p90_idx]}ms"
        echo "  P95延迟: ${sorted[$p95_idx]}ms"
        echo "  P99延迟: ${sorted[$p99_idx]}ms"
        echo "  最大延迟: ${sorted[-1]}ms"
        echo "  最小延迟: ${sorted[0]}ms"
    else
        echo -e "${RED}没有成功的请求用于延迟统计${NC}"
    fi
    echo ""
}

# 测试4: 角色差异化限流测试
test_role_based_limiting() {
    echo -e "${BLUE}=== 测试4: 角色差异化限流测试 ===${NC}"
    
    local roles=("guest" "common" "member" "super_member" "super_admin")
    local requests_per_role=20
    
    for role in "${roles[@]}"; do
        echo "测试角色: $role"
        
        local success=0
        local limited=0
        local error=0
        
        for i in $(seq 1 $requests_per_role); do
            response=$(curl -s -w "%{http_code}" -o /dev/null \
                -F "audio=@/tmp/test_audio.wav" \
                -H "X-User-ID: test_user_$i" \
                -H "X-User-Role: $role" \
                "$BASE_URL$SPEECH_ENDPOINT" 2>/dev/null)
            
            case $response in
                200) ((success++)) ;;
                429) ((limited++)) ;;
                *) ((error++)) ;;
            esac
        done
        
        local limit_rate=$(echo "scale=2; $limited * 100 / $requests_per_role" | bc)
        echo "  成功: $success, 被限流: $limited, 错误: $error, 限流率: ${limit_rate}%"
    done
    echo ""
}

# 测试5: 内存和CPU使用率监控
test_resource_usage() {
    echo -e "${BLUE}=== 测试5: 资源使用率监控 ===${NC}"
    
    echo "开始监控资源使用率..."
    
    # 启动资源监控
    (
        while true; do
            echo "$(date '+%H:%M:%S'),$(ps -o pid,pcpu,pmem,comm -p $(pgrep -f 'ai_jianli_go') | tail -1 | awk '{print $2","$3","$4}')"
            sleep 1
        done
    ) > /tmp/resource_monitor.log &
    
    local monitor_pid=$!
    
    # 执行压力测试
    echo "执行压力测试..."
    for i in $(seq 1 200); do
        curl -s -o /dev/null \
            -F "audio=@/tmp/test_audio.wav" \
            -H "X-User-ID: test_user_$i" \
            -H "X-User-Role: common" \
            "$BASE_URL$SPEECH_ENDPOINT" &
        
        if [ $((i % 20)) -eq 0 ]; then
            wait
        fi
    done
    wait
    
    # 停止监控
    kill $monitor_pid 2>/dev/null
    wait $monitor_pid 2>/dev/null
    
    # 分析资源使用情况
    if [ -f /tmp/resource_monitor.log ]; then
        echo -e "${GREEN}资源使用统计:${NC}"
        awk -F',' 'NR>1 {cpu+=$2; mem+=$3; count++} END {
            if(count>0) {
                printf "  平均CPU使用率: %.2f%%\n", cpu/count
                printf "  平均内存使用率: %.2f%%\n", mem/count
                printf "  监控时长: %d秒\n", count
            }
        }' /tmp/resource_monitor.log
        rm -f /tmp/resource_monitor.log
    fi
    echo ""
}

# 测试6: 限流恢复测试
test_rate_limit_recovery() {
    echo -e "${BLUE}=== 测试6: 限流恢复测试 ===${NC}"
    
    echo "测试限流后的恢复能力..."
    
    # 快速发送请求触发限流
    echo "触发限流..."
    for i in $(seq 1 50); do
        curl -s -o /dev/null \
            -F "audio=@/tmp/test_audio.wav" \
            -H "X-User-ID: test_user" \
            -H "X-User-Role: common" \
            "$BASE_URL$SPEECH_ENDPOINT" &
    done
    wait
    
    # 等待限流恢复
    echo "等待限流恢复..."
    sleep 5
    
    # 测试恢复后的请求
    echo "测试恢复后的请求..."
    local success=0
    local limited=0
    
    for i in $(seq 1 10); do
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -F "audio=@/tmp/test_audio.wav" \
            -H "X-User-ID: test_user" \
            -H "X-User-Role: common" \
            "$BASE_URL$SPEECH_ENDPOINT" 2>/dev/null)
        
        case $response in
            200) ((success++)) ;;
            429) ((limited++)) ;;
        esac
        
        sleep 1
    done
    
    echo -e "${GREEN}恢复测试结果:${NC}"
    echo "  成功请求: $success"
    echo "  被限流请求: $limited"
    echo ""
}

# 测试7: 限流统计信息测试
test_rate_limit_stats() {
    echo -e "${BLUE}=== 测试7: 限流统计信息测试 ===${NC}"
    
    echo "查询限流统计信息..."
    
    # 获取全局限流统计
    echo "全局限流统计:"
    curl -s -H "X-User-Role: super_admin" \
        "$BASE_URL$RATE_LIMIT_ENDPOINT/stats" | jq '.' 2>/dev/null || echo "需要管理员权限或jq工具"
    
    echo ""
    echo "被限流最多的键:"
    curl -s -H "X-User-Role: super_admin" \
        "$BASE_URL$RATE_LIMIT_ENDPOINT/top-limited?limit=5" | jq '.' 2>/dev/null || echo "需要管理员权限或jq工具"
    
    echo ""
    echo "限流系统健康状态:"
    curl -s "$BASE_URL$RATE_LIMIT_ENDPOINT/health" | jq '.' 2>/dev/null || echo "需要jq工具"
    
    echo ""
}

# 生成测试报告
generate_report() {
    echo -e "${BLUE}=== 生成测试报告 ===${NC}"
    
    local report_file="/tmp/rate_limit_performance_report_$(date +%Y%m%d_%H%M%S).txt"
    
    {
        echo "限流性能测试报告"
        echo "=================="
        echo "测试时间: $(date)"
        echo "测试环境: $BASE_URL"
        echo ""
        echo "测试项目:"
        echo "1. 基础限流功能测试"
        echo "2. 高并发压力测试"
        echo "3. 延迟测试"
        echo "4. 角色差异化限流测试"
        echo "5. 资源使用率监控"
        echo "6. 限流恢复测试"
        echo "7. 限流统计信息测试"
        echo ""
        echo "建议:"
        echo "- 根据测试结果调整限流参数"
        echo "- 监控生产环境的限流效果"
        echo "- 定期进行性能测试"
    } > "$report_file"
    
    echo -e "${GREEN}测试报告已生成: $report_file${NC}"
}

# 主函数
main() {
    echo -e "${PURPLE}=== AI简历助手限流性能测试 ===${NC}"
    echo ""
    
    # 设置错误处理
    set -e
    trap cleanup EXIT
    
    # 检查依赖
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}错误: 需要安装curl${NC}"
        exit 1
    fi
    
    if ! command -v bc &> /dev/null; then
        echo -e "${RED}错误: 需要安装bc${NC}"
        exit 1
    fi
    
    # 执行测试
    check_server
    create_test_audio
    
    test_basic_rate_limit
    test_high_concurrency
    test_latency
    test_role_based_limiting
    test_resource_usage
    test_rate_limit_recovery
    test_rate_limit_stats
    
    generate_report
    
    echo -e "${GREEN}=== 所有测试完成 ===${NC}"
    echo ""
    echo -e "${YELLOW}测试说明:${NC}"
    echo "1. 基础限流测试: 验证限流器基本功能"
    echo "2. 高并发测试: 测试系统在高并发下的表现"
    echo "3. 延迟测试: 测量请求响应时间"
    echo "4. 角色测试: 验证不同角色的限流策略"
    echo "5. 资源监控: 监控CPU和内存使用情况"
    echo "6. 恢复测试: 测试限流后的恢复能力"
    echo "7. 统计测试: 验证限流统计功能"
    echo ""
    echo -e "${BLUE}性能优化建议:${NC}"
    echo "- 根据测试结果调整限流参数"
    echo "- 监控生产环境的限流效果"
    echo "- 定期进行性能测试"
    echo "- 考虑使用Redis集群提高限流性能"
}

# 运行主函数
main "$@"
