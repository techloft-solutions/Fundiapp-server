package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
	return
}

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

func (s *Server) handleLocationList(w http.ResponseWriter, r *http.Request) {
	var authUser = &app.AuthUser{}
	ctx := r.Context()
	authUser = middlewares.UserFromContext(ctx)
	resp, err := s.LcSvc.ListMyLocations(ctx, authUser)
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
	return
}

func (s *Server) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	location.Title = r.PostFormValue("title")
	//location.Latitude = r.PostFormValue("latitude")
	//location.Longitude = r.PostFormValue("longitude")
	location.Address = r.PostFormValue("address")

	err := s.LcSvc.CreateLocation(r.Context(), &location)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleUnathorised(w)
		return
	}

	handleSuccess(w, nil)
}
