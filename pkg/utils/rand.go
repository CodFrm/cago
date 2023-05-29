package utils

import (
	"math/rand"
	"time"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6 // 6 bits to represent a letter index
)

var src = rand.NewSource(time.Now().UnixNano())

type RandStringType int

const (
	Number RandStringType = iota // 数字
	Letter                       // 小写数字
	Mix                          // 大小写数字混合
)

func RandString(n int, randType RandStringType) string {
	b := make([]byte, n)
	l := int(10 + (randType * 26))
	for i := 0; i < n; i++ {
		b[i] = letterBytes[rand.Intn(l)]
	}
	return string(b)
}
