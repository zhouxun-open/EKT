package x_resp

import (
	"encoding/json"
	"github.com/EducationEKT/xserver/x_err"
)

type XRespBody struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

func (xRespBody *XRespBody) ToBytes() (bytes []byte) {
	bytes, _ = json.Marshal(xRespBody)
	return
}

func Success(result interface{}) *XRespContainer {
	body := &XRespBody{0, "ok", result}
	return &XRespContainer{
		HttpCode: 200,
		Headers:  make(map[string]string),
		Body:     body.ToBytes(),
	}
}

func Fail(status int, msg string, result interface{}) *XRespContainer {
	body := &XRespBody{status, msg, result}
	return &XRespContainer{
		HttpCode: 200,
		Headers:  make(map[string]string),
		Body:     body.ToBytes(),
	}
}

func Return(result interface{}, err error) (*XRespContainer, *x_err.XErr) {
	return Success(result), x_err.NewXErr(err)
}
