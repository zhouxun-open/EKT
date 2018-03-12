package x_http

import (
	"github.com/EducationEKT/xserver/x_http/x_processer"
	"net/http"
)

func Service(w http.ResponseWriter, r *http.Request) {
	x_processer.Process(w, r)
}
