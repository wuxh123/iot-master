package base

import (
	"bytes"
	"math/rand"
	"time"
)

const _chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(size int) string {
	rand.Seed(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i ++ {
		s.WriteByte(_chars[rand.Int63() % int64(len(_chars))])
	}
	return s.String()
}

func RandomNumber(size int) string {
	rand.Seed(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i ++ {
		s.WriteByte('0' + uint8(rand.Intn(10)))
	}
	return s.String()
}
