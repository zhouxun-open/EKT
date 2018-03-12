package x_processer

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_parser"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

var xParser = &x_parser.XParser{}

func Process(w http.ResponseWriter, r *http.Request) {
	req := xParser.Parse(w, r)
	defer func(req *x_req.XReq) {
		if r := recover(); r != nil {
			var resp *x_resp.XRespContainer
			if reflect.TypeOf(r) == reflect.TypeOf(x_err.NewParamErr()) {
				errMsg, _ := json.Marshal(r.(*x_err.LogicalErr))
				resp = &x_resp.XRespContainer{Body: errMsg}
			} else {
				resp = x_resp.Fail(-1, "unknown exception", nil)
			}
			req.SendResponse(resp)
		}
	}(req)
	var resp *x_resp.XRespContainer
	if routerInfo := x_router.Router.GetXRouter(req.Path); routerInfo != nil {
		resp = routerInfo.Invoke(req)
	} else {
		resp = x_router.NotFound(req)
	}
	req.SendResponse(resp)
	accessLogInfo(req)
}

func invoke(req *x_req.XReq) *x_resp.XRespContainer {
	resp := x_resp.Success(req.Context)
	return resp
}

func accessLogInfo(req *x_req.XReq) (logInfo string) {
	//logInfo = x_string.StringJoint(
	//	x_time.DEFAULT_DATE_TIME, TAB,
	//	resp.Context.ReqId, TAB,
	//	resp.Context.IP, TAB,
	//)
	return
}
