package randutil

import (
	"math/rand"
	"time"
)

var (
	characters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString 生成随机字符串
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}
