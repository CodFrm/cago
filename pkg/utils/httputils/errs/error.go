package errs

import (
	"go.uber.org/zap"
)

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

func WarpStack(err error) error {
	return &Error{
		err:   err,
		field: []zap.Field{zap.StackSkip("stack", 2)},
	}
}

func Warp(err error, msg string, field ...zap.Field) error {
	return &Error{
		err:   err,
		msg:   msg,
		field: field,
	}
}
