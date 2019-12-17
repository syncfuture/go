package u

import (
	"errors"
	"reflect"
	"runtime"
	"strings"

	log "github.com/kataras/golog"
)

// JointErrors joint errors to a single error
func JointErrors(errs ...error) error {
	var count = len(errs)
	if count == 0 {
		return nil
	}

	sb := strings.Builder{}

	for _, err := range errs {
		if err != nil && errs != nil {
			sb.WriteString(err.Error() + ";")
		}
	}

	if sb.Len() > 0 {
		err := errors.New(sb.String())
		return err
	}

	return nil
}

func LogFaltal(err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Fatalf("<%s:%d> %v", file, line, err)
		} else {
			log.Fatal(err)
		}
	}
}

// LogError 记录有错误
func LogError(err error) bool {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Errorf("<%s:%d> %v", file, line, err)
		} else {
			log.Error(err)
		}
		return true
	}

	return false
}

// LogErrorMsg 记录错误消息
func LogErrorMsg(err error, mrPtr interface{}) bool {

	if err != nil {
		_, file, line, ok := runtime.Caller(1)

		if ok {
			log.Errorf("<%s:%d> %v", file, line, err)
		} else {
			log.Error(err)
		}

		if reflect.TypeOf(mrPtr).Kind() != reflect.Ptr {
			panic("mrPtr must be a pointer")
		}

		msgCodeField := reflect.ValueOf(mrPtr).Elem().FieldByName("MsgCode")
		if !msgCodeField.CanSet() {
			panic("MsgCode field doesn't exist or cannot be set")
		} else if msgCodeField.Kind() != reflect.String {
			panic("MsgCode must be a string field")
		}

		msgCodeField.SetString(err.Error())

		return true
	}

	return false
}
