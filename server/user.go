package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
)

func (s *Server) handleUserValidate(w http.ResponseWriter, r *http.Request) {

	type user struct {
		Password   string `valid:"required" json:"password"`
		Phone      string `valid:"required" json:"phone"`
		IsProvider bool   `json:"provider"`
	}
	var usr user

	jsonStr, err := json.Marshal(allMpFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &usr); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	usr.IsProvider, _ = strconv.ParseBool(r.URL.Query().Get("provider"))

	_, err = govalidator.ValidateStruct(usr)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//usr.Phone = govalidator.Trim(usr.Phone, "")
	//usr.Password = govalidator.Trim(usr.Password, "")

	log.Println("[LOG] phone: [" + usr.Phone + "] password: [" + usr.Password + "]")

	err = s.UsrSvc.ValidateUser(r.Context(), usr.Phone, usr.Password, usr.IsProvider)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "Incorrect phone number or password", http.StatusUnauthorized)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsg(w, "User is valid")

}
