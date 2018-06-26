package xlog

import (
	"testing"
	"time"
)

func TestDailyLog(t *testing.T) {
	log := NewDailyLog("/tmp/daily.log")
	log.Info("hello world")
	log.Debug("I'm here")
	log.Error("file not found")
	log.Warn("Listen to me")
	log.Crit("I'm crash")
	log.Debug("hello %s", "golang")
	time.Sleep(3 * time.Second)
}
