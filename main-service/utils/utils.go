package utils

import (
	"encoding/json"
	"net/http"
)

// ResponseJSON standard json response helper
func ResponseJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// ResponseError standard json response helper
func ResponseError(w http.ResponseWriter, code int, message string) {
	ResponseJSON(w, code, map[string]string{"error": message})
}
