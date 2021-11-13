package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	app "github.com/andrwkng/hudumaapp"
)

// allFormValues returns a map that contains all the form values.
// Instead of multiple values for the same key, only the first value is returned as string.
func allFormValues(r *http.Request) map[string]string {
	r.ParseForm()
	m := make(map[string]string)
	for k, v := range r.PostForm {
		m[k] = v[0]
	}
	return m
}

func handleUnathorised(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := make(map[string]string)
	resp["message"] = "Unauthorized"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func handleSuccess(w http.ResponseWriter, resource interface{}) {
	jsonResp, err := json.Marshal(resource)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func handleSuccessText(w http.ResponseWriter, resource interface{}) {
	jsonResp, err := json.Marshal(resource)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func handleError(w http.ResponseWriter, err error, code int) {
	resp := make(map[string]app.Error)
	resp["error"] = app.Error{
		Code:    strconv.Itoa(code),
		Message: err.Error(),
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(jsonResp)
}
