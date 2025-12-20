package sid

import (
	"math/rand"
	"strings"
)

// 自定义字符集
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func New(n int) string {
	var sb strings.Builder
	length := len(alphabet)

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(length)])
	}

	return sb.String()
}
