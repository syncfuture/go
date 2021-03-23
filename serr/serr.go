package serr

import (
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func New(msg string) error {
	return errors.New(msg)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// func Wrap(err error, msg string) error {
// 	return errors.Wrap(err, msg)
// }

// func Wrapf(err error, format string, args ...interface{}) error {
// 	return errors.Wrapf(err, format, args...)
// }

// func Cause(err error) error {
// 	return errors.Cause(err)
// }

// func Is(err, target error) bool {
// 	return errors.Is(err, target)
// }

func Wrap(err error) error {
	_, ok := err.(stackTracer)
	if ok {
		return err
	}

	return errors.WithStack(err)
}
