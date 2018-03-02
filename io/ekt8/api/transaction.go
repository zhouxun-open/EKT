package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

var pubKey, secKey []byte

func init() {
	pubKey, secKey = crypto.GenerateKeyPair()
	fmt.Println(hex.EncodeToString(pubKey))
	msg := make(map[string]interface{})
	msg["Hello"] = "World"
	msg["zhou"] = "xun"
	data, _ := json.Marshal(msg)
	sign, _ := crypto.Crypto(data, secKey)
	fmt.Println(hex.EncodeToString(sign))
	x_router.Post("/transaction/api/newTransaction", ValidateSign, newTransaction)
}

func ValidateSign(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	sign := req.Param["sign"]
	msg := req.Param["msg"]
	data, _ := json.Marshal(msg)
	if crypto.Verify([]byte(sign.(string)), pubKey, data) {
		return nil, nil
	}
	return nil, x_err.New(-1, "Invalid Signature")
}

func newTransaction(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return x_resp.Success(make(map[string]interface{})), nil
}
