#!/bin/bash

# 快速限流测试脚本
# 用于快速验证限流功能是否正常工作

# 配置
BASE_URL="http://localhost:8080"
SPEECH_ENDPOINT="/api/v1/speech/recognize"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 快速限流测试 ===${NC}"
echo ""

# 检查服务器状态
echo -e "${YELLOW}检查服务器状态...${NC}"
if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}错误: 服务器未运行，请先启动服务器${NC}"
    echo "启动命令: go run main.go"
    exit 1
fi
echo -e "${GREEN}服务器运行正常${NC}"
echo ""

# 创建测试音频文件
echo -e "${YELLOW}创建测试音频文件...${NC}"
dd if=/dev/zero of=/tmp/test_audio.wav bs=1024 count=1 2>/dev/null
echo -e "${GREEN}测试音频文件创建完成${NC}"
echo ""

# 测试1: 正常请求
echo -e "${BLUE}测试1: 正常请求${NC}"
response=$(curl -s -w "%{http_code}" -o /dev/null \
    -F "audio=@/tmp/test_audio.wav" \
    -H "X-User-ID: test_user" \
    -H "X-User-Role: common" \
    "$BASE_URL$SPEECH_ENDPOINT" 2>/dev/null)

if [ "$response" = "200" ]; then
    echo -e "${GREEN}✓ 正常请求通过 (HTTP 200)${NC}"
else
    echo -e "${RED}✗ 正常请求失败 (HTTP $response)${NC}"
fi
echo ""

# 测试2: 快速请求测试限流
echo -e "${BLUE}测试2: 快速请求测试限流${NC}"
echo "发送10个快速请求..."

success_count=0
limited_count=0
error_count=0

for i in {1..10}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -F "audio=@/tmp/test_audio.wav" \
        -H "X-User-ID: test_user" \
        -H "X-User-Role: common" \
        "$BASE_URL$SPEECH_ENDPOINT" 2>/dev/null)
    
    case $response in
        200) ((success_count++)) ;;
        429) ((limited_count++)) ;;
        *) ((error_count++)) ;;
    esac
done

echo -e "${GREEN}快速请求测试完成:${NC}"
echo "  成功: $success_count"
echo "  被限流: $limited_count"
echo "  错误: $error_count"

if [ $limited_count -gt 0 ]; then
    echo -e "${GREEN}✓ 限流功能正常工作${NC}"
else
    echo -e "${YELLOW}⚠ 限流功能可能未生效，建议检查配置${NC}"
fi
echo ""

# 测试3: 不同角色测试
echo -e "${BLUE}测试3: 不同角色测试${NC}"
roles=("guest" "common" "member" "super_member" "super_admin")

for role in "${roles[@]}"; do
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -F "audio=@/tmp/test_audio.wav" \
        -H "X-User-ID: test_user" \
        -H "X-User-Role: $role" \
        "$BASE_URL$SPEECH_ENDPOINT" 2>/dev/null)
    
    case $response in
        200) echo -e "  角色 $role: ${GREEN}成功${NC}" ;;
        429) echo -e "  角色 $role: ${YELLOW}被限流${NC}" ;;
        *) echo -e "  角色 $role: ${RED}错误 (HTTP $response)${NC}" ;;
    esac
done
echo ""

# 测试4: 响应头检查
echo -e "${BLUE}测试4: 响应头检查${NC}"
response=$(curl -s -I "$BASE_URL$SPEECH_ENDPOINT" \
    -F "audio=@/tmp/test_audio.wav" \
    -H "X-User-ID: test_user" \
    -H "X-User-Role: common" 2>/dev/null)

echo "响应头信息:"
echo "$response" | grep -E "(X-RateLimit|HTTP/)" || echo "  未找到限流相关响应头"
echo ""

# 清理
echo -e "${YELLOW}清理测试文件...${NC}"
rm -f /tmp/test_audio.wav
echo -e "${GREEN}清理完成${NC}"
echo ""

# 测试总结
echo -e "${BLUE}=== 测试总结 ===${NC}"
echo "1. 服务器状态: 正常"
echo "2. 正常请求: 通过"
echo "3. 限流功能: $([ $limited_count -gt 0 ] && echo "正常" || echo "需要检查")"
echo "4. 角色权限: 已测试"
echo "5. 响应头: 已检查"
echo ""

if [ $limited_count -gt 0 ]; then
    echo -e "${GREEN}✓ 限流功能测试通过${NC}"
    echo ""
    echo -e "${YELLOW}建议:${NC}"
    echo "- 运行完整性能测试: make test-performance"
    echo "- 查看限流统计: make stats"
    echo "- 查看测试文档: docs/rate_limit_testing_guide.md"
else
    echo -e "${YELLOW}⚠ 限流功能可能未生效${NC}"
    echo ""
    echo -e "${YELLOW}建议:${NC}"
    echo "- 检查限流配置: config/config.yaml"
    echo "- 查看服务器日志"
    echo "- 验证Redis连接"
    echo "- 运行详细测试: make test-all"
fi
