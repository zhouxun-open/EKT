package log

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/EducationEKT/EKT/io/ekt8/conf"
)

func TestXLoger(t *testing.T) {
	os.Chdir("/opt/gopath/src/github.com/EducationEKT/EKT")
	conf.InitConfig("genesis.conf")
	err := InitLog()
	if err != nil {
		fmt.Printf("%v\n", err)
		t.FailNow()
	}
	xLogger := GetLogInst()
	xLogger.LogCrit("this is a crit log 2")
	xLogger.LogCrit("this is a crit log 3")
	xLogger.LogCrit("this is a crit log 1")
	xLogger.LogCrit("this is a crit log 5")
	xLogger.LogCrit("this is a crit log 4")
	xLogger.LogInfo("this is a Info log 2")
	xLogger.LogInfo("this is a Info log 3")
	xLogger.LogInfo("this is a Info log 1")
	xLogger.LogInfo("this is a Info log 5")
	xLogger.LogInfo("this is a Info log 4")
	xLogger.LogDebug("this is a Debug log 2")
	xLogger.LogDebug("this is a Debug log 3")
	xLogger.LogDebug("this is a Debug log 1")
	xLogger.LogDebug("this is a Debug log 5")
	xLogger.LogDebug("this is a Debug log 4")
	xLogger.LogErr("this is a Err log 2")
	xLogger.LogErr("this is a Err log 3")
	xLogger.LogErr("this is a Err log 1")
	xLogger.LogErr("this is a Err log 5")
	xLogger.LogErr("this is a Err log 4")
	time.Sleep(1 * time.Second)
}
