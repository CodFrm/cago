package utils

import "testing"

func TestRandString(t *testing.T) {
	// 随机10次数字
	for i := 0; i < 10; i++ {
		result := RandString(6, Number)
		// 检查是否只包含数字
		for _, v := range result {
			if v < '0' || v > '9' {
				t.Error("RandString(6, Number) failed")
			}
		}
	}
	// 随机10次数字+小写字母
	for i := 0; i < 10; i++ {
		result := RandString(6, Letter)
		// 检查是否只包含数字和小写字母
		for _, v := range result {
			if v < '0' || (v > '9' && v < 'a') || v > 'z' {
				t.Error("RandString(6, Letter) failed")
			}
		}
	}
}
