package x_time

import (
	"time"
)

const (
	DEFAULT_DATE      = "2006-01-02"
	DEFAULT_TIME      = "15:04:05"
	DEFAULT_DATE_TIME = "2006-01-02 15:04:05"
)

func Now() int64 {
	return time.Now().UnixNano() / 1e6
}

func DateStr() string {
	return time.Now().Format(DEFAULT_DATE)
}

func TimeStr() string {
	return time.Now().Format(DEFAULT_TIME)
}

func DateTimeStr() string {
	return time.Now().Format(DEFAULT_DATE_TIME)
}
