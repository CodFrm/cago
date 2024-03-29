package errs

import (
	"go.uber.org/zap"
)

// Unwrap 用于获取错误的原始错误
type Unwrap interface {
	Unwrap() error
}

type Error struct {
	err   error
	msg   string
	field []zap.Field
}

func (s *Error) Error() string {
	return s.err.Error()
}

func (s *Error) Unwrap() error {
	return s.err
}

func (s *Error) Field() []zap.Field {
	return s.field
}

// WrapStack 包装一个错误，同时记录堆栈信息
func WrapStack(err error) error {
	return Wrap(err, "", zap.StackSkip("stack", 2))
}

// Wrap 包装一个错误，同时记录错误信息
func Wrap(err error, msg string, field ...zap.Field) error {
	return &Error{
		err:   err,
		msg:   msg,
		field: field,
	}
}
