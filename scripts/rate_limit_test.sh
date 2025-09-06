#!/bin/bash

# 限流功能测试脚本
# 用于测试语音识别接口的限流效果

# 配置
BASE_URL="http://localhost:8080"
SPEECH_ENDPOINT="/api/v1/speech/recognize"
USER_ENDPOINT="/api/v1/user/login"
RATE_LIMIT_ENDPOINT="/api/v1/ratelimit"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== AI简历助手限流功能测试 ===${NC}"
echo ""

# 检查服务器是否运行
echo -e "${YELLOW}检查服务器状态...${NC}"
if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}错误: 服务器未运行，请先启动服务器${NC}"
    exit 1
fi
echo -e "${GREEN}服务器运行正常${NC}"
echo ""

# 测试1: 语音识别接口限流测试
echo -e "${BLUE}测试1: 语音识别接口限流测试${NC}"
echo "发送100个并发请求到语音识别接口..."

# 创建临时音频文件用于测试
echo "创建测试音频文件..."
dd if=/dev/zero of=/tmp/test_audio.wav bs=1024 count=1 2>/dev/null

# 并发测试
start_time=$(date +%s.%N)
for i in {1..100}; do
    (
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -F "audio=@/tmp/test_audio.wav" \
            -H "X-User-ID: test_user_$i" \
            "$BASE_URL$SPEECH_ENDPOINT")
        
        if [ "$response" = "429" ]; then
            echo "请求 $i: 被限流 (HTTP 429)"
        elif [ "$response" = "200" ]; then
            echo "请求 $i: 成功 (HTTP 200)"
        else
            echo "请求 $i: 其他状态 (HTTP $response)"
        fi
    ) &
done

# 等待所有请求完成
wait
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)

echo -e "${GREEN}并发测试完成，耗时: ${duration}秒${NC}"
echo ""

# 测试2: 用户认证接口限流测试
echo -e "${BLUE}测试2: 用户认证接口限流测试${NC}"
echo "发送30个并发请求到登录接口..."

start_time=$(date +%s.%N)
for i in {1..30}; do
    (
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X POST \
            -H "Content-Type: application/json" \
            -d '{"username":"test","password":"test"}' \
            "$BASE_URL$USER_ENDPOINT")
        
        if [ "$response" = "429" ]; then
            echo "登录请求 $i: 被限流 (HTTP 429)"
        elif [ "$response" = "200" ] || [ "$response" = "401" ]; then
            echo "登录请求 $i: 处理成功 (HTTP $response)"
        else
            echo "登录请求 $i: 其他状态 (HTTP $response)"
        fi
    ) &
done

wait
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)

echo -e "${GREEN}认证接口测试完成，耗时: ${duration}秒${NC}"
echo ""

# 测试3: 限流统计信息查询
echo -e "${BLUE}测试3: 限流统计信息查询${NC}"
echo "查询限流统计信息..."

# 注意：这里需要管理员权限，实际使用时需要提供有效的管理员token
echo "获取全局限流统计信息:"
curl -s -H "X-User-Role: super_admin" \
    "$BASE_URL$RATE_LIMIT_ENDPOINT/stats" | jq '.' 2>/dev/null || echo "需要管理员权限或jq工具"

echo ""
echo "获取被限流最多的键:"
curl -s -H "X-User-Role: super_admin" \
    "$BASE_URL$RATE_LIMIT_ENDPOINT/top-limited?limit=5" | jq '.' 2>/dev/null || echo "需要管理员权限或jq工具"

echo ""
echo "获取限流系统健康状态:"
curl -s "$BASE_URL$RATE_LIMIT_ENDPOINT/health" | jq '.' 2>/dev/null || echo "需要jq工具"

echo ""

# 测试4: 响应头验证
echo -e "${BLUE}测试4: 响应头验证${NC}"
echo "检查限流响应头..."

response=$(curl -s -I "$BASE_URL$SPEECH_ENDPOINT" \
    -F "audio=@/tmp/test_audio.wav" \
    -H "X-User-ID: test_user")

echo "响应头信息:"
echo "$response" | grep -E "(X-RateLimit|HTTP/)"

echo ""

# 测试5: 不同用户角色的限流测试
echo -e "${BLUE}测试5: 不同用户角色限流测试${NC}"

roles=("guest" "common" "member" "super_member" "super_admin")

for role in "${roles[@]}"; do
    echo "测试角色: $role"
    
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -F "audio=@/tmp/test_audio.wav" \
        -H "X-User-ID: test_user" \
        -H "X-User-Role: $role" \
        "$BASE_URL$SPEECH_ENDPOINT")
    
    if [ "$response" = "429" ]; then
        echo "  角色 $role: 被限流 (HTTP 429)"
    elif [ "$response" = "200" ]; then
        echo "  角色 $role: 成功 (HTTP 200)"
    else
        echo "  角色 $role: 其他状态 (HTTP $response)"
    fi
done

echo ""

# 清理临时文件
echo -e "${YELLOW}清理临时文件...${NC}"
rm -f /tmp/test_audio.wav

echo -e "${GREEN}=== 限流功能测试完成 ===${NC}"
echo ""
echo -e "${BLUE}测试总结:${NC}"
echo "1. 语音识别接口限流测试 - 验证高QPS和大burst配置"
echo "2. 用户认证接口限流测试 - 验证严格限流防暴力破解"
echo "3. 限流统计信息查询 - 验证监控功能"
echo "4. 响应头验证 - 验证限流信息传递"
echo "5. 角色权限测试 - 验证不同角色的限流策略"
echo ""
echo -e "${YELLOW}注意: 某些测试需要管理员权限，请确保提供有效的认证信息${NC}"
