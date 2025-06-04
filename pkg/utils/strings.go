package utils

import "strings"

func extractString(input string, prefix string) string {
	index := strings.Index(input, prefix)

	if index >= 0 {
		// 提取前缀之后的所有内容
		result := input[index+len(prefix):]
		// 去除开头可能的多余空格
		return strings.TrimSpace(result)
	}

	return "" // 未找到匹配
}
