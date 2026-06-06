package sid

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// 自定义字符集
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func New(n int) string {
	var sb strings.Builder
	length := big.NewInt(int64(len(alphabet)))

	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, length)
		if err != nil {
			panic(err)
		}
		sb.WriteByte(alphabet[idx.Int64()])
	}

	return sb.String()
}
