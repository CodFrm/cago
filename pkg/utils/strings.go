package utils

import "fmt"

func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func ToNumber[T int8 | int16 | int | int32 | int64](str string) T {
	var val T
	_, err := fmt.Sscanf(str, "%d", &val)
	if err != nil {
		return 0
	}
	return val
}
