package xlog

import (
	"fmt"
	"os"
	"time"
)

// Writer object
type Writer struct {
	path string
	fp   *os.File
}

// NewDailyWriter create a daily writer
func NewDailyWriter(path string) *Writer {
	return &Writer{path: path}
}

func (w *Writer) getLogPath() string {
	return fmt.Sprintf("%s.%s", w.path, time.Now().Format("2006-01-02"))
}

func (w *Writer) checkLogFile() error {
	path := w.getLogPath()
	if w.fp == nil {
		fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		w.fp = fp
		os.Remove(w.path)
		os.Symlink(path, w.path)
	}

	if w.fp.Name() != path {
		w.fp.Close()
		fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		w.fp = fp
		os.Remove(w.path)
		os.Symlink(path, w.path)
	}

	return nil
}

func (w *Writer) Write(b []byte) (n int, err error) {
	err = w.checkLogFile()
	if err != nil {
		return 0, err
	}
	return w.fp.Write(b)
}
