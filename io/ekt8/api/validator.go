package api

import (
	"encoding/hex"
	"encoding/json"
	"github.com/EducationEKT/EKT/io/ekt8/crypto"
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
)

func ValidateSign(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	sign := req.Param["sign"]
	msg := req.Param["msg"]
	data, _ := json.Marshal(msg)
	data = crypto.Sha3_256(data)
	signByte, _ := hex.DecodeString(sign.(string))
	if crypto.Verify(signByte, pubKey, data) {
		return nil, nil
	}
	return nil, x_err.New(-1, "Invalid Signature")
}
