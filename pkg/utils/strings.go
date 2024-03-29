package utils

import (
	"fmt"
	"strings"
)

// ToString 任意类型转字符串
func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

// ToNumber 字符串转数字
func ToNumber[T int8 | int16 | int | int32 | int64](str string) T {
	var val T
	_, err := fmt.Sscanf(str, "%d", &val)
	if err != nil {
		return 0
	}
	return val
}

// StringReverse 字符串反转
func StringReverse(s string) string {
	a := []rune(s)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return string(a)
}

// Head 返回字符串str中第一个sep之前的字符串和第一个sep之后的字符串
func Head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}
