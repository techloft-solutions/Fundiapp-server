package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
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
/*
func (s *Server) handlePortfolioCreate(w http.ResponseWriter, r *http.Request) {
	var portfolio model.Portfolio
	portfolio.Title = r.PostFormValue("title")
	bookingId, err := uuid.Parse(r.PostFormValue("booking_id"))
	if err != nil {
		handleError(w, app.Errorf(app.INVALID_ERR, "Invalid ID format"), 400)
		return
	}
	portfolio.BookingID = bookingId

	photoData := r.PostFormValue("photos")
	var photos []string
	err = json.Unmarshal([]byte(photoData), &photos)
	if err != nil {
		panic(err)
	}
	portfolio.Photos = photos

	err = s.PfoSvc.CreatePortfolio(r.Context(), &portfolio)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleUnathorised(w)
		return
	}

	handleSuccess(w, nil)
}
*/
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
		log.Println(err)
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

	err = user.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.UsrSvc.CreateUser(r.Context(), &user)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	// obfuscate the password
	user.Password = strings.Repeat("*", len(user.Password))

	handleSuccessMsgWithRes(w, "User created successfuly", user)
}

func (s *Server) handleUserGet(w http.ResponseWriter, r *http.Request) {
	var res user
	urlquery := r.URL.Query()
	res.Username = urlquery.Get("displayname")
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
		Password string `valid:"required" json:"password"`
		Phone    string `valid:"required" json:"phone"`
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

	_, err = govalidator.ValidateStruct(usr)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.UsrSvc.ValidateUser(r.Context(), usr.Phone, usr.Password)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "Incorrect phone number or password", http.StatusUnauthorized)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, nil)
}
