package sid

import (
	"math/rand"
	"strings"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Generate 生成长度为 n 的可视化随机字符串（A-Z, a-z, 0-9）
func Generate() string {
	var sb strings.Builder
	sb.Grow(8)
	for i := 0; i < 8; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
