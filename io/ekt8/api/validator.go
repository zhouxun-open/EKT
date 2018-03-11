package api

import (
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	// 客户端发送一条消息和一个签名到节点上，节点对签名进行校验，校验成功返回world代表用户的签名正确，否则标识签名错误
	// 此接口一般用于用户刚更新完公私钥对，在转账前对新私钥的校验
	x_router.All("/hello", ValidateSign, world)
}

func ValidateSign(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	//sign := req.Param["sign"]
	//msg := req.Param["msg"]
	//data, _ := json.Marshal(msg)
	//data = crypto.Sha3_256(data)
	//signByte, _ := hex.DecodeString(sign.(string))
	//if crypto.Verify(signByte, pubKey, data) {
	//	return nil, nil
	//}
	//return nil, x_err.New(-1, "Invalid Signature")
	return nil, nil
}

func world(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	return &x_resp.XRespContainer{
		Body: []byte("world"),
	}, nil
}
