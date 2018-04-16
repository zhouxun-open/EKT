package x_resp

import (
	"net/http"
)

const (
	TAB  = "\t"
	LINE = "\n"
)

type XRespContainer struct {
	HttpCode int
	Headers  map[string]string
	Cookies  []http.Cookie
	Body     []byte
}
