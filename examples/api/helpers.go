package main

import (
	"encoding/json"
	"net/http"
)

// fail writes a json response with error msg and status header
func fail(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)

	data := struct {
		Error string `json:"error"`
	}{Error: msg}
	resp, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// ok writes data to response with 200 status
func ok(w http.ResponseWriter, data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		fail(w, "Oops something evil has happened", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
