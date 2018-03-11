package api

import (
	"encoding/json"

	"github.com/EducationEKT/EKT/io/ekt8/dispatcher"
	"github.com/EducationEKT/EKT/io/ekt8/event"
	"github.com/EducationEKT/EKT/io/ekt8/p2p"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/user/api/regist")
}

func NewAccount(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	var account event.NewAccountParam
	err := json.Unmarshal(req.Body, &account)
	if err != nil {
		return x_resp.Fail(-1, err.Error(), nil), x_err.New(-1, err.Error())
	}
	evt := &event.Event{EventParam: account, EventType: event.NewAccountEvent}
	dispatcher.GetDisPatcher().NewEvent(evt)
	if !p2p.IsDPosPeer(req.R.RemoteAddr) {
		p2p.BroadcastRequest(req.Path, req.Body)
	}
	return x_resp.Success("success"), nil
}
