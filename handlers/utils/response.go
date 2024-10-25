package utils

import (
	"encoding/json"
	"net/http"
)

func Send(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func SendError(w http.ResponseWriter, err string) {
	type response struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	Send(w, response{
		Status: http.StatusText(http.StatusBadRequest),
		Error:  err,
	})
}
