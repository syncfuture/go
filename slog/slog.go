package slog

import (
	"fmt"
	"runtime"
	"time"

	"github.com/syncfuture/go/sconfig"

	"github.com/kataras/golog"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var (
	_detailMap = map[string]int{
		"all":   0,
		"debug": 1,
		"info":  2,
		"warn":  3,
		"error": 4,
		"fatal": 5,
	}
	_detailLevel = _detailMap["warn"]
	Config       = &LogConfig{
		Level:       "debug",
		DetailLevel: "warn",
	}
	DebugLogger ILogger = new(debugLogger)
)

type LogConfig struct {
	Level       string
	DetailLevel string
	File        string
}

type ILogger interface {
	Printf(format string, args ...interface{})
}

type debugLogger struct{}

func (self *debugLogger) Printf(format string, args ...interface{}) {
	golog.Debugf(format, args...)
}

func Init(configProvider sconfig.IConfigProvider) {
	if configProvider == nil {
		golog.Fatal("configProvider cannot be nil")
	}

	configProvider.GetStruct("Log", &Config)
	if Config == nil {
		golog.Fatal("Cannot find 'Log' section in configuration")
	}

	if Config.Level == "" {
		Config.Level = "debug"
	}
	if Config.DetailLevel == "" {
		Config.DetailLevel = "warn"
	}
	_detailLevel = _detailMap[Config.DetailLevel]

	if Config.Level == "all" { // golog 的debug就会显示所有
		golog.SetLevel("debug")
	} else {
		golog.SetLevel(Config.Level)
	}

	if Config.File != "" {
		rotationSeconds := configProvider.GetIntDefault("Log.RotationSeconds", 86400) // 默认24小时一个新日志文件
		rotationCount := configProvider.GetIntDefault("Log.RotationCount", 7)         // 默认最多保存7个日志文件
		writer, err := rotatelogs.New(
			Config.File+".%Y%m%d%H%M%S",
			rotatelogs.WithRotationTime(time.Duration(rotationSeconds)*time.Second), //
			rotatelogs.WithRotationCount(uint(rotationCount)),
		)
		if err != nil {
			golog.Fatal(err)
		}

		golog.SetOutput(writer)
	}
}

func Debug(v ...interface{}) {
	if _detailLevel <= 1 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Debugf("<%s:%d> %+v", file, line, v)
			return
		}
	}

	golog.Debug(v...)
}
func Debugf(format string, args ...interface{}) {
	if _detailLevel <= 1 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Debugf("<%s:%d> %+v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Debugf(format, args...)
}

func Info(v ...interface{}) {
	if _detailLevel <= 2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Infof("<%s:%d> %+v", file, line, v)
			return
		}
	}

	golog.Info(v...)
}
func Infof(format string, args ...interface{}) {
	if _detailLevel <= 2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Infof("<%s:%d> %+v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Infof(format, args...)
}

func Warn(v ...interface{}) {
	if _detailLevel <= 3 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Warnf("<%s:%d> %+v", file, line, v)
			return
		}
	}

	golog.Warn(v...)
}
func Warnf(format string, args ...interface{}) {
	if _detailLevel <= 3 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Warnf("<%s:%d> %+v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Warnf(format, args...)
}

func Error(v ...interface{}) {
	if _detailLevel <= 4 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Errorf("<%s:%d> %+v", file, line, v)
			return
		}
	}

	golog.Error(v...)
}
func Errorf(format string, args ...interface{}) {
	if _detailLevel <= 4 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Errorf("<%s:%d> %+v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Errorf(format, args...)
}

func Fatal(v ...interface{}) {
	if _detailLevel <= 5 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Fatalf("<%s:%d> %+v", file, line, v)
			return
		}
	}

	golog.Fatal(v...)
}
func Fatalf(format string, args ...interface{}) {
	if _detailLevel <= 5 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Fatalf("<%s:%d> %+v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Fatalf(format, args...)
}

func Println(log string) {
	golog.Println(log)
}

func Print(log string) {
	golog.Print(log)
}
