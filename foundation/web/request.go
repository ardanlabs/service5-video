package web

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}
