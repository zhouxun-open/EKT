package ctxlog

import (
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/log"
)

type ContextLog struct {
	ContextInfo string
	Sticker     map[string]interface{}
	Timings     map[string]int64
}

func NewContextLog(name string) *ContextLog {
	return &ContextLog{
		ContextInfo: name,
		Sticker:     make(map[string]interface{}),
		Timings:     make(map[string]int64),
	}
}

func (log *ContextLog) Log(key string, value interface{}) {
	log.Sticker[key] = value
}

func (log *ContextLog) LogTiming(key string, value int64) {
	log.Timings[key] = value
}

func (cLog *ContextLog) Finish() {
	log.Debug(cLog.String())
}

func (cLog *ContextLog) String() string {
	sticker, _ := json.Marshal(cLog.Sticker)
	timings, _ := json.Marshal(cLog.Timings)
	return fmt.Sprintf(`%s : {"sticker": %s, "timings": %s}`, cLog.ContextInfo, string(sticker), string(timings))
}
