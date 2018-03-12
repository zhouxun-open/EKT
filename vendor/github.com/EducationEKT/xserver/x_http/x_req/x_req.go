package x_req

import (
	"encoding/json"

	"github.com/EducationEKT/xserver/x_err"
	"github.com/EducationEKT/xserver/x_http/x_context"
	"github.com/EducationEKT/xserver/x_http/x_resp"
	"github.com/EducationEKT/xserver/x_utils/x_time"
	"github.com/EducationEKT/xserver/x_utils/x_type"

	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type XReq struct {
	R            *http.Request
	W            http.ResponseWriter
	Method       string
	Path         string
	Context      *x_context.XContext
	StartTime    int64
	Query        map[string]interface{}
	Param        map[string]interface{}
	Cookies      map[string]string
	HandlerParam map[string]interface{}
	Body         []byte
}

func (req *XReq) SendResponse(resp *x_resp.XRespContainer) {
	w := req.W
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cost", fmt.Sprintf("%d", x_time.Now()-req.StartTime))
	if resp.Headers != nil && 0 != len(resp.Headers) {
		for key, value := range resp.Headers {
			w.Header().Del(key)
			w.Header()[key] = []string{0: value}
		}
	}
	w.WriteHeader(resp.HttpCode)
	w.Write(resp.Body)
}

func New(r *http.Request, w http.ResponseWriter) *XReq {
	return &XReq{
		R:            r,
		W:            w,
		Path:         r.URL.Path,
		StartTime:    x_time.Now(),
		Context:      x_context.NewXContext(w, r),
		HandlerParam: make(map[string]interface{}),
	}
}

func (req *XReq) ParseRequest() *x_err.XErr {
	req.parseCookie()
	req.Method = req.R.Method
	req.parseQuery()
	if strings.ToUpper(req.Method) == "POST" {
		if err := req.parseBody(); err != nil {
			return err
		}
	}
	return nil
}

func (req *XReq) parseCookie() {
	if cookies := req.R.Cookies(); cookies != nil && len(cookies) != 0 {
		req.Cookies = make(map[string]string)
		for _, cookie := range cookies {
			req.Cookies[cookie.Name] = cookie.Value
		}
	}
}

func (req *XReq) parseQuery() {
	query := req.R.RequestURI[strings.LastIndex(req.R.RequestURI, "?")+1:]
	req.Query = kv2map(query)
}

func (req *XReq) parseBody() *x_err.XErr {
	body, err := ioutil.ReadAll(req.R.Body)
	req.Body = body
	if err != nil {
		return x_err.New(-1, "read body error")
	}
	contentType := req.R.Header.Get("Content-Type")
	if contentType == "application/json" {
		err = json.Unmarshal(body, &req.Param)
		if err != nil {
			return x_err.New(-1, "body parse error")
		}
	} else if contentType == "application/x-www-from-urlencoded" {
		bodyStr := string(body)
		req.Param = kv2map(bodyStr)
	} else {
		return x_err.New(-7, "Content-Type not support")
	}
	return nil
}

func kv2map(kv string) map[string]interface{} {
	pairs := strings.Split(kv, "&")
	result := make(map[string]interface{})
	if len(pairs) > 0 {
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				result[kv[0]] = kv[1]
			}
		}
	}
	return result
}

func (req *XReq) GetParam(name string) (interface{}, bool) {
	v, exist := req.Param[name]
	if !exist {
		v, exist = req.Query[name]
	}
	return v, exist
}

func (req *XReq) MustGetInt64(name string) int64 {
	v, exist := req.GetParam(name)
	if !exist {
		panic(x_err.NewParamErr())
	}
	value, ok := x_type.GetInt64(v)
	if !ok {
		panic(x_err.NewParamErr())
	}
	return value
}

func (req *XReq) MustGetInt32(name string) int32 {
	v, exist := req.GetParam(name)
	if !exist {
		panic(x_err.NewParamErr())
	}
	value, ok := x_type.GetInt32(v)
	if !ok {
		panic(x_err.NewParamErr())
	}
	return value
}

func (req *XReq) MustGetFloat64(name string) float64 {
	v, exist := req.GetParam(name)
	if !exist {
		panic(x_err.NewParamErr())
	}
	value, ok := x_type.GetFloat64(v)
	if !ok {
		panic(x_err.NewParamErr())
	}
	return value
}

func (req *XReq) MustGetFloat32(name string) float64 {
	v, exist := req.GetParam(name)
	if !exist {
		panic(x_err.NewParamErr())
	}
	value, ok := x_type.GetFloat32(v)
	if !ok {
		panic(x_err.NewParamErr())
	}
	return value
}

func (req *XReq) MustGetBool(name string) bool {
	v, exist := req.GetParam(name)
	if !exist {
		panic(x_err.NewParamErr())
	}
	value, ok := x_type.GetBool(v)
	if !ok {
		panic(x_err.NewParamErr())
	}
	return value
}

func (req *XReq) MustGetString(name string) string {
	v, exist := req.GetParam(name)
	if !exist {
		panic(x_err.NewParamErr())
	}
	return x_type.V2String(v)
}
