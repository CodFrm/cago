package errs

import "errors"

// WarnError 包装一个错误，表示这个错误是一个警告
type WarnError struct {
	error
}

func (s *WarnError) Error() string {
	return s.error.Error()
}

func (s *WarnError) Unwrap() error {
	return s.error
}

// Warn 包装一个错误，表示这个错误是一个警告
func Warn(err error) error {
	return &WarnError{
		error: err,
	}
}

// IsWarn 判断一个错误是否是一个警告
func IsWarn(err error) *WarnError {
	var e *WarnError
	if errors.As(err, &e) {
		return e
	}
	return nil
}
