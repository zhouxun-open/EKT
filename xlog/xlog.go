package xlog

import (
	"fmt"
	"log"
	"os"
	"path"
)

const (
	logChanSize = 1 << 16
)

type xlog struct {
	c chan string
	w *log.Logger
}

// XLog is the log interface
type XLog interface {
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Crit(msg string, args ...interface{})
}

func checklogPath(filename string) {
	fileinfo, err := os.Stat(filename)
	if err == nil && fileinfo.IsDir() {
		panic("log path must be is file")
	}
	dir := path.Dir(filename)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		panic(fmt.Sprintf("create log dir failed, %s", err.Error()))
	}
}

// NewDailyLog create a daily log wrapper
func NewDailyLog(path string) XLog {
	checklogPath(path)
	_log := xlog{w: log.New(NewDailyWriter(path), "", log.LstdFlags), c: make(chan string, logChanSize)}
	go _log.writer()
	return &_log
}

func (log *xlog) Debug(msg string, args ...interface{}) {
	log.c <- fmt.Sprintf("[debug] %s", fmt.Sprintf(msg, args...))
}

func (log *xlog) Warn(msg string, args ...interface{}) {
	log.c <- fmt.Sprintf("[warn] %s", fmt.Sprintf(msg, args...))
}

func (log *xlog) Info(msg string, args ...interface{}) {
	log.c <- fmt.Sprintf("[info] %s", fmt.Sprintf(msg, args...))
}

func (log *xlog) Error(msg string, args ...interface{}) {
	log.c <- fmt.Sprintf("[error] %s", fmt.Sprintf(msg, args...))
}

func (log *xlog) Crit(msg string, args ...interface{}) {
	log.c <- fmt.Sprintf("[crit] %s", fmt.Sprintf(msg, args...))
}

func (log *xlog) writer() {
	for {
		log.w.Println(<-log.c)
	}
}
