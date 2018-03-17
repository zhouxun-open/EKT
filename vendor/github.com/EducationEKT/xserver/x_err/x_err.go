package x_err

import (
	"encoding/json"
)

type XErr struct {
	Level           int    `json:"level"`
	Msg             string `json:"msg"`
	SuggestContinue bool   `json:"Continue"`
}

type LogicalErr struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func (xerr *XErr) Error() string {
	bytes, _ := json.Marshal(xerr)
	return string(bytes)
}

func New(level int, msg string) *XErr {
	return &XErr{level, msg, false}
}

func NewParamErr() *LogicalErr {
	return &LogicalErr{-7, "param is missing"}
}

func NewXErr(err error) *XErr {
	if err == nil {
		return nil
	}
	return &XErr{-1, err.Error(), false}
}
