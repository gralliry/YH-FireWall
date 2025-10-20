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

func NewWithNoRepeat(n int) string {
	runes := []rune(alphabet) // 支持 Unicode
	length := len(runes)

	if n > length {
		n = length // 防止 n 太大
	}

	// Fisher–Yates 洗牌
	for i := length - 1; i > length-1-n; i-- {
		j := rand.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes[length-n:])
}
