package api

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/EducationEKT/EKT/crypto"
	"github.com/EducationEKT/EKT/db"
	"github.com/EducationEKT/EKT/log"

	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_http/x_router"
)

func init() {
	x_router.Post("/db/api/get", GetValue)
	x_router.Get("/db/api/getByHex", GetValueByHexHash)
}

var (
	InvalidKey = errors.New("Invalid Key")
	NotFound   = errors.New("Not found")
)

func GetValueByHexHash(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	hash := req.MustGetString("hash")
	key, err := hex.DecodeString(hash)
	if err != nil {
		return x_resp.Return(key, err)
	}
	v, err := GetValueByHash(key)
	return validate(key, v, err)
}

func GetValue(req *x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr) {
	v, err := GetValueByHash(req.Body)
	return validate(req.Body, v, err)
}

func GetValueByHash(key []byte) ([]byte, error) {
	if len(key) != 32 {
		log.Info("Remote peer want a db value that len(key) is not 32 byte, return fail.")
		return nil, InvalidKey
	}
	return db.GetDBInst().Get(key)
}

func validate(k, v []byte, err error) (*x_resp.XRespContainer, *x_err.XErr) {
	if err != nil {
		return x_resp.Return(nil, err)
	}
	if !bytes.Equal(crypto.Sha3_256(v), k) {
		log.Info("This key is not the hash of the db value, return fail.")
		return x_resp.Fail(-403, "Invalid Key", string(k)), nil
	}
	return &x_resp.XRespContainer{
		HttpCode: 200,
		Body:     v,
	}, nil
}
