package server

import (
	"encoding/json"
	"log"
	"net/http"
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

func handleSuccess(w http.ResponseWriter, resource interface{}) {
	jsonResp, err := json.Marshal(resource)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func handleSuccessMsg(w http.ResponseWriter, msg string) {
	resp := make(map[string]string)
	resp["success"] = msg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func handleSuccessMsgWithRes(w http.ResponseWriter, msg string, res interface{}) {
	resp := make(map[string]interface{})
	resp["success"] = msg
	resp["data"] = res
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func handleError(w http.ResponseWriter, errorMsg string, code int) {
	resp := make(map[string]string)
	resp["error"] = errorMsg
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(jsonResp)
}

func handleUnathorised(w http.ResponseWriter) {
	resp := make(map[string]string)
	resp["error"] = "Unauthorized"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(jsonResp)
}
