package handler

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func writeJSON(rw http.ResponseWriter, status int, data any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(data) //nolint:errcheck
}

func writeError(rw http.ResponseWriter, status int, msg string) {
	writeJSON(rw, status, errorResponse{Message: msg})
}
