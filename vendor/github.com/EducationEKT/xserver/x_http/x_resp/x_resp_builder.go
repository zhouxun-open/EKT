package x_resp

import (
	"encoding/json"
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
