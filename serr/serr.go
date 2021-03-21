package serr

import (
	"github.com/pkg/errors"
)

func New(msg string) error {
	return errors.New(msg)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

// func WithStack(err error) error {
// 	return errors.WithStack(err)
// }
