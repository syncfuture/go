package slog

import (
	"runtime"

	log "github.com/kataras/golog"
)

func Debug(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Debugf("<%s:%d> %v", file, line, v)
	} else {
		log.Debug(v...)
	}
}
func Debugf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Debugf("<%s:%d> %v", file, line, args)
	} else {
		log.Debugf(format, args...)
	}
}

func Info(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Infof("<%s:%d> %v", file, line, v)
	} else {
		log.Info(v...)
	}
}
func Infof(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Infof("<%s:%d> %v", file, line, args)
	} else {
		log.Infof(format, args...)
	}
}

func Warn(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Warnf("<%s:%d> %v", file, line, v)
	} else {
		log.Warn(v...)
	}
}
func Warnf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Warnf("<%s:%d> %v", file, line, args)
	} else {
		log.Warnf(format, args...)
	}
}

func Error(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Errorf("<%s:%d> %v", file, line, v)
	} else {
		log.Error(v...)
	}
}
func Errorf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Errorf("<%s:%d> %v", file, line, args)
	} else {
		log.Errorf(format, args...)
	}
}

func Fatal(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Fatalf("<%s:%d> %v", file, line, v)
	} else {
		log.Fatal(v...)
	}
}
func Fatalf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Fatalf("<%s:%d> %v", file, line, args)
	} else {
		log.Fatalf(format, args...)
	}
}
