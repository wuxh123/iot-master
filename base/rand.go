package base

import (
	"bytes"
	"math/rand"
	"time"
)

const _chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(size int) string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i ++ {
		s.WriteByte(_chars[rand.Int63() % int64(len(_chars))])
	}
	return s.String()
}
