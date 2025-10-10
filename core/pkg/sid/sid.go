package sid

import (
	"math/rand"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Generate 生成长度为 n 的可视化随机字符串（A-Z, a-z, 0-9）
func Generate() string {
	l := len(charset)
	sid := make([]byte, 8)
	for i := 0; i < 8; i++ {
		sid[i] = charset[rand.Intn(l)]
	}
	return string(sid)
}
