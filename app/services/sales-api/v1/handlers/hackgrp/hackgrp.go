package hackgrp

import (
	"encoding/json"
	"net/http"
)

func Hack(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(status)
}
