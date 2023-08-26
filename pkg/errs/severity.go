package errs

import "errors"

type WarnError struct {
	error
}

func (s *WarnError) Error() string {
	return s.error.Error()
}

func (s *WarnError) Unwrap() error {
	return s.error
}

func Warn(err error) error {
	return &WarnError{
		error: err,
	}
}

func IsWarn(err error) *WarnError {
	var e *WarnError
	if errors.As(err, &e) {
		return e
	}
	return nil
}
