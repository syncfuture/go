package slog

import (
	"fmt"
	"os"
	"runtime"

	"github.com/syncfuture/go/config"

	"github.com/kataras/golog"
)

var (
	Level      = "warn"
	_detailMap = map[string]int{
		"all":   0,
		"debug": 1,
		"info":  2,
		"warn":  3,
		"error": 4,
		"fatal": 5,
	}
	_detailLevel = _detailMap["warn"]
)

func Init(args ...string) {
	var configFile string
	if len(args) > 0 {
		configFile = args[0]
	} else {
		configFile = "configs.json"
	}

	_, err := os.Stat(configFile)
	if err == nil {
		cp := config.NewJsonConfigProvider(configFile)
		Level = cp.GetStringDefault("Log.Level", "warn")
		detailLevel := cp.GetString("Log.DetailLevel") // 显示文件行数与否的级别
		if detailLevel == "" {
			detailLevel = Level
		}
		_detailLevel = _detailMap[detailLevel]
	}

	golog.SetLevel(Level)
}

func Debug(v ...interface{}) {
	if _detailLevel <= 1 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Debugf("<%s:%d> %v", file, line, v)
			return
		}
	}

	golog.Debug(v...)
}
func Debugf(format string, args ...interface{}) {
	if _detailLevel <= 1 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Debugf("<%s:%d> %v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Debugf(format, args...)
}

func Info(v ...interface{}) {
	if _detailLevel <= 2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Infof("<%s:%d> %v", file, line, v)
			return
		}
	}

	golog.Info(v...)
}
func Infof(format string, args ...interface{}) {
	if _detailLevel <= 2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Infof("<%s:%d> %v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Infof(format, args...)
}

func Warn(v ...interface{}) {
	if _detailLevel <= 3 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Warnf("<%s:%d> %v", file, line, v)
			return
		}
	}

	golog.Warn(v...)
}
func Warnf(format string, args ...interface{}) {
	if _detailLevel <= 3 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Warnf("<%s:%d> %v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Warnf(format, args...)
}

func Error(v ...interface{}) {
	if _detailLevel <= 4 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Errorf("<%s:%d> %v", file, line, v)
			return
		}
	}

	golog.Error(v...)
}
func Errorf(format string, args ...interface{}) {
	if _detailLevel <= 4 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Errorf("<%s:%d> %v", file, line, fmt.Sprintf(format, args...))
			return
		}
	}

	golog.Errorf(format, args...)
}

func Fatal(v ...interface{}) {
	if _detailLevel <= 5 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Fatalf("<%s:%d> %v", file, line, v)
			return
		}
	}

	golog.Fatal(v...)
}
func Fatalf(format string, args ...interface{}) {
	if _detailLevel <= 5 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			golog.Fatalf("<%s:%d> %v", file, line, fmt.Sprintf(format, args...))
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
