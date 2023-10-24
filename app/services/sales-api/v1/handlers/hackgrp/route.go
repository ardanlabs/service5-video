package hackgrp

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// Routes adds specific routes for this group.
func Routes(mux *httptreemux.ContextMux) {
	mux.Handle(http.MethodGet, "/hack", Hack)
}
