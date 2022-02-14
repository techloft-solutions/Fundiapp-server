package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type user struct {
	Username string `json:"username,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

/*
func (s *Server) handlePortfolioList(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["user_id"]
	// Fetch portfolios from database.
	resp, err := s.PfoSvc.ListPortfoliosByUserId(r.Context(), id)
	if err != nil {
		log.Println(err)
		handleError(w, errors.New("Something went wrong!"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
*/

func (s *Server) handlePortfolioCreate(w http.ResponseWriter, r *http.Request) {
	var portfolio model.Portfolio

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	portfolio.UserID = userID.String()

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	if err := json.Unmarshal(jsonStr, &portfolio); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	photoData := r.PostFormValue("photos")
	portfolio.Photos, err = retrievePhotos(photoData)
	if err != nil {
		handleError(w, "photos: invalid json array value", http.StatusBadRequest)
		return
	}

	err = portfolio.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.PfoSvc.CreatePortfolio(r.Context(), &portfolio)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Portfolio created successfully", portfolio)
}

func (s *Server) handleMyPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	resp, err := s.PfoSvc.ListPortfoliosByUserId(ctx, userID.String())
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}

func (s *Server) handlePortfolio(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, "Id is not a valid UUID", http.StatusBadRequest)
		return
	}

	portfolio, err := s.PfoSvc.FindPortfolioByID(r.Context(), id)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "Portfolio not found", http.StatusNotFound)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, portfolio)
}

func (s *Server) handleMyLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	resp, err := s.LocSvc.ListMyLocations(ctx, userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}

func (s *Server) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}
	location.UserID = userID.String()
	location.ID = uuid.New()

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &location); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	err = location.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.LocSvc.CreateLocation(r.Context(), &location)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Location created successfuly", location)
}

func (s *Server) handleLocationDelete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.LocSvc.RemoveLocation(r.Context(), id)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsg(w, "Location deleted successfuly")
}

func (s *Server) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &user); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	user.IsProvider, _ = strconv.ParseBool(r.URL.Query().Get("provider"))

	err = user.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.UsrSvc.CreateUser(r.Context(), &user)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	//retrieveFirebaseUserData(r.Context(), user.ID.String())

	// obfuscate the password
	user.Password = strings.Repeat("*", len(user.Password))

	handleSuccessMsgWithRes(w, "User created successfuly", user)
}

//func linkFirebaseUserToUser(userID string, firebaseUserID string) error {

func (s *Server) handleUserGet(w http.ResponseWriter, r *http.Request) {
	var res user
	urlquery := r.URL.Query()
	res.Username = urlquery.Get("display_name")
	res.Phone = urlquery.Get("phone")

	switch {
	case res.Username != "":
		_, err := s.UsrSvc.FindUserByUsername(r.Context(), res.Username)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			if err == sql.ErrNoRows {
				handleError(w, "User not found", http.StatusNotFound)
				return
			}
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		handleSuccess(w, res)
	case res.Phone != "":
		_, err := s.UsrSvc.FindUserByPhoneNumber(r.Context(), res.Phone)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			if err == sql.ErrNoRows {
				handleError(w, "User not found", http.StatusNotFound)
				return
			}
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		handleSuccess(w, res)
	default:
		handleError(w, "Invalid request", http.StatusBadRequest)
	}

}

func (s *Server) handleUserValidate(w http.ResponseWriter, r *http.Request) {

	type user struct {
		Password   string `valid:"required" json:"password"`
		Phone      string `valid:"required" json:"phone"`
		IsProvider bool   `json:"provider"`
	}
	var usr user

	jsonStr, err := json.Marshal(allFormValues(r))
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

func (s *Server) handlePasswordNew(w http.ResponseWriter, r *http.Request) {
	newPassWord := r.FormValue("new_password")
	if newPassWord == "" {
		handleError(w, "New password is required", http.StatusBadRequest)
		return
	}

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	log.Println("[http]", r.Method, r.URL.Path, "userID:", userID, "new_password:["+newPassWord+"]")

	err = s.UsrSvc.ResetUserPassword(r.Context(), newPassWord, userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "Failed to reset password", http.StatusNotFound)
			return
		}
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsg(w, "Password reset successfuly")
}

func (s *Server) handlePasswordChange(w http.ResponseWriter, r *http.Request) {
	var pw model.PwdChange

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &pw); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}
	pw.UserID = userID.String()

	_, err = govalidator.ValidateStruct(pw)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("[LOG] OldPassword: [" + pw.OldPassword + "] NewPassword: [" + pw.NewPassword + "]")

	err = s.UsrSvc.ChangeUserPassword(r.Context(), &pw)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "Incorrect current password", http.StatusNotFound)
			return
		}
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsg(w, "Password updated successfuly")
}

func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	usr, err := s.UsrSvc.FindUserByID(r.Context(), id)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "User not found", http.StatusNotFound)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, usr)
}

func retrievePhotos(photoData string) ([]string, error) {
	var photos []string
	if photoData != "" {
		err := json.Unmarshal([]byte(photoData), &photos)
		if err != nil {
			return nil, err
		}
	}
	return photos, nil
}
