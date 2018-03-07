package log

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/conf"
)

const MAX_LOG_LENGTH = 1 << 16

var xLoger *XLoger

type XLoger struct {
	dictPath string
	Info     chan string
	Err      chan string
	Debug    chan string
	Crit     chan string
}

func GetLogInst() *XLoger {
	return xLoger
}

func InitLog() error {
	fileInfo, err := os.Stat(conf.EKTConfig.LogPath)
	if err != nil {
		err = os.MkdirAll(conf.EKTConfig.LogPath, 0777)
		if err != nil {
			return err
		}
		return InitLog()
	}
	if !fileInfo.IsDir() {
		return errors.New("dir needed, given file")
	}
	xLoger = &XLoger{dictPath: conf.EKTConfig.LogPath, Info: make(chan string, MAX_LOG_LENGTH), Err: make(chan string, MAX_LOG_LENGTH), Crit: make(chan string, MAX_LOG_LENGTH), Debug: make(chan string, MAX_LOG_LENGTH)}
	go xLoger.run()
	return nil
}

func (xLoger *XLoger) LogInfo(formatter string, args ...interface{}) {
	info := fmt.Sprintf(formatter, args...)
	xLoger.Info <- info
}

func (xLogger *XLoger) LogErr(formatter string, args ...interface{}) {
	err := fmt.Sprintf(formatter, args...)
	xLoger.Err <- err
}

func (xLoger *XLoger) LogCrit(formatter string, args ...interface{}) {
	crit := fmt.Sprintf(formatter, args...)
	xLoger.Crit <- crit
}

func (xLogger *XLoger) LogDebug(formatter string, args ...interface{}) {
	debug := fmt.Sprintf(formatter, args...)
	xLoger.Debug <- debug
}

func (xLog *XLoger) run() {
	for {
		select {
		case acc, ok := <-xLog.Info:
			if ok {
				logInfo(xLog, acc)
			}
		case err, ok := <-xLog.Err:
			if ok {
				logErr(xLog, err)
			}
		case crit, ok := <-xLog.Crit:
			if ok {
				logCrit(xLog, crit)
			}
		case debug, ok := <-xLog.Debug:
			if ok {
				logDebug(xLog, debug)
			}
		}
	}
}

func logInfo(xLog *XLoger, acc string) {
	path := xLog.jointLogPath("info")
	AppendFile(path, acc)
}

func logErr(xLog *XLoger, err string) {
	path := xLog.jointLogPath("err")
	AppendFile(path, err)
}

func logCrit(xLog *XLoger, crit string) {
	path := xLog.jointLogPath("crit")
	AppendFile(path, crit)
}

func logDebug(xLog *XLoger, debug string) {
	path := xLog.jointLogPath("debug")
	AppendFile(path, debug)
}

func (xLog *XLoger) jointLogPath(logLevel string) string {
	b := bytes.Buffer{}
	b.WriteString(xLog.dictPath)
	dict := b.String()
	if !strings.HasSuffix(dict, "/") {
		b.WriteString("/")
	}
	os.MkdirAll(dict, os.ModePerm)
	b.WriteString(logLevel)
	b.WriteString(time.Now().Format(".2006-01-02.log"))
	return b.String()
}

func AppendFile(filename string, data string) {
	logFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0666)
	defer logFile.Close()
	logFile.WriteString(time.Now().String() + "   " + data + "\n")
}
