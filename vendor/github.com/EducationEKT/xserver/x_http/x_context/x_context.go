package x_context

import (
	"net/http"
	"github.com/EducationEKT/xserver/x_utils/x_time"
)

type XContext struct {
	W         http.ResponseWriter    `json:"omitempty"`
	R         *http.Request          `json:"omitempty"`
	StartTime int64                  `json:"startTime"`
	IP        string                 `json:"ip"`
	ReqId     string                 `json:"req_id"`
	LogParam  map[string]interface{} `json:"log_param"`
	LogTime   map[string]int64       `json:"log_time"`
	Sticker   map[string]interface{} `json:"sticker"`
}

func NewXContext(w http.ResponseWriter, r *http.Request) *XContext {
	return &XContext{
		W:         w,
		R:         r,
		IP:        r.RemoteAddr,
		StartTime: x_time.Now(),
		LogTime:   make(map[string]int64),
		LogParam:  make(map[string]interface{}),
		Sticker:   make(map[string]interface{}),
	}
}
