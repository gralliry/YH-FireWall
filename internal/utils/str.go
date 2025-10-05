package utils

import (
	"math/rand"
	"regexp"
	"strings"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateID 生成长度为 n 的可视化随机字符串（A-Z, a-z, 0-9）
func GenerateID(n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

var re = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// IsValidName 判断字符串是否只由字母、数字、- 和 _ 组成
func IsValidName(s string) bool {
	// ^ 开头，$ 结尾，[] 内为允许字符，+ 表示至少一个
	return re.MatchString(s)
}
