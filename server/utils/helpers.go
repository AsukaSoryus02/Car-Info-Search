package utils

import (
	"strconv"
	"strings"
)

// FormatInt36 将整数格式化为base36字符串
// 用于生成更短的ID字符串
func FormatInt36(n int64) string {
	return strings.ToLower(strconv.FormatInt(n, 36))
}

// ParseInt36 将base36字符串解析为整数
func ParseInt36(s string) (int64, error) {
	return strconv.ParseInt(strings.ToLower(s), 36, 64)
}

// TruncateString 截断字符串到指定长度，并添加省略号
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ContainsString 检查字符串切片是否包含指定字符串
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// RemoveString 从字符串切片中移除指定字符串
func RemoveString(slice []string, s string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}
