package sid

import (
	"crypto/rand"
	"unsafe"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

const mask = 248 // floor(256/62)*62，保证均匀分布

// New 生成 n 位随机字符串
func New(n int) string {
	buf := make([]byte, n)
	tmp := make([]byte, 32)

	for i := 0; i < n; {
		_, _ = rand.Read(tmp)
		for _, b := range tmp {
			if i >= n {
				break
			}
			if b < mask {
				buf[i] = alphabet[b%62]
				i++
			}
		}
	}
	return unsafe.String(unsafe.SliceData(buf), n)
}
