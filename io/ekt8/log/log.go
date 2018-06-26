package log

import (
	"sync"

	"github.com/EducationEKT/EKT/io/ekt8/conf"
	"github.com/EducationEKT/EKT/io/ekt8/xlog"
)

var once sync.Once
var l xlog.XLog

func InitLog() {
	once.Do(func() {
		l = xlog.NewDailyLog(conf.EKTConfig.LogPath)
	})
}

func Debug(msg string, args ...interface{}) {
	l.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	l.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	l.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	l.Warn(msg, args...)
}

func Crit(msg string, args ...interface{}) {
	l.Crit(msg, args...)
}
