package x_router

import (
	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_req"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"sync"
)

const (
	GET  = 1 << iota
	POST = 1 << iota
	ALL  = GET | POST
)

type XHandler func(*x_req.XReq) (*x_resp.XRespContainer, *x_err.XErr)

type XRouterInfo struct {
	HttpMethod int32 //GET POST PUT OPTION DELETE ...
	Handlers   []XHandler
}

type XException struct {
}

type XRouter struct {
	mapping map[string]*XRouterInfo
	Lock    sync.Locker
}

var Router = &XRouter{
	mapping: make(map[string]*XRouterInfo),
	Lock:    new(sync.RWMutex),
}

func (xRouter *XRouter) GetXRouter(url string) *XRouterInfo {
	return xRouter.mapping[url]
}

func (routeInfo *XRouterInfo) Invoke(req *x_req.XReq) *x_resp.XRespContainer {
	var resp *x_resp.XRespContainer
	for _, handler := range routeInfo.Handlers {
		var err *x_err.XErr
		resp, err = handler(req)
		if nil != err {
			if nil == resp {
				body := &x_resp.XRespBody{Status: -1, Msg: err.Msg, Result: make(map[string]interface{})}
				resp = &x_resp.XRespContainer{HttpCode: 200, Body: body.ToBytes()}
			}
			break
		}
	}
	return resp
}

func Get(url string, handlers ...XHandler) {
	Router.Lock.Lock()
	defer Router.Lock.Unlock()
	routerInfo := &XRouterInfo{
		HttpMethod: GET,
		Handlers:   handlers,
	}
	Router.mapping[url] = routerInfo
}

func Post(url string, handlers ...XHandler) {
	Router.Lock.Lock()
	defer Router.Lock.Unlock()
	routeInfo := &XRouterInfo{
		HttpMethod: POST,
		Handlers:   handlers,
	}
	Router.mapping[url] = routeInfo
}

func All(url string, handlers ...XHandler) {
	Router.Lock.Lock()
	defer Router.Lock.Unlock()
	routeInfo := &XRouterInfo{
		HttpMethod: ALL,
		Handlers:   handlers,
	}
	Router.mapping[url] = routeInfo
}

func NotFound(req *x_req.XReq) *x_resp.XRespContainer {
	resp := &x_resp.XRespContainer{
		HttpCode: 404,
		Body:     []byte("{\"status\":-1, \"msg\": \"not found\"}"),
	}
	return resp
}
