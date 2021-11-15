package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/google/uuid"
)

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
