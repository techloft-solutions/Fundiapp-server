package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) handleBookingByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID from path.
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, "Id is not a valid UUID", http.StatusBadRequest)
		return
	}

	resp, err := s.BkSvc.FindBookingByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			handleError(w, "Booking not found", http.StatusNotFound)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}

func (s *Server) handleBookingList(w http.ResponseWriter, r *http.Request) {
	booking, err := s.BkSvc.FindMyBookings(r.Context())
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, booking)
}

func (s *Server) handleProviderBookings(w http.ResponseWriter, r *http.Request) {
	providerID := mux.Vars(r)["id"]
	booking, err := s.BkSvc.FindBookings(r.Context(), providerID)
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, booking)
}

func (s *Server) handleProviderBooking(w http.ResponseWriter, r *http.Request) {
	// Parse ID from path.
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, "Id is not a valid UUID", http.StatusBadRequest)
		return
	}

	userId, err := middlewares.UserIDFromContext(r.Context())
	if err != nil {
		handleUnathorised(w)
		return
	}

	resp, err := s.BkSvc.FindProviderBookingByID(r.Context(), id, userId.String())
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			handleError(w, "Booking not found", http.StatusNotFound)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}

func (s *Server) handleBookingCreate(w http.ResponseWriter, r *http.Request) {
	var booking model.Booking

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &booking); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	booking.ID = uuid.New()

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	var photos []string
	photoData := r.PostFormValue("photos")
	if photoData != "" {
		err = json.Unmarshal([]byte(photoData), &photos)
		if err != nil {
			handleError(w, "photos: invalid json array value", http.StatusBadRequest)
			return
		}
		booking.Photos = photos
	}

	booking.ClientID = userID.String()
	// validate
	err = booking.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	time, err := dateparse.ParseStrict(booking.StartDate)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "start_date: invalid date format", http.StatusBadRequest)
		return
	}

	booking.StartDate = time.Format("2006-01-02T15:04:05")

	err = s.BkSvc.CreateBooking(r.Context(), &booking)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	handleSuccessMsgWithRes(w, "Booking created successfully", booking)
}

func (s *Server) handleBookingComplete(w http.ResponseWriter, r *http.Request) {
	bookingId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, "Id is not a valid UUID", http.StatusBadRequest)
		return
	}
	err = s.BkSvc.CompleteBooking(r.Context(), bookingId)
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, "Booking marked completed successfully")
}

func (s *Server) handleBookingCancel(w http.ResponseWriter, r *http.Request) {
	bookingId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, "Id is not a valid UUID", http.StatusBadRequest)
		return
	}
	err = s.BkSvc.CancelBooking(r.Context(), bookingId)
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, "Booking marked cancelled successfully")
}

/*
func retrieveFirebaseUserData(ctx context.Context, uid string) *auth.UserRecord {
	opt := option.WithCredentialsFile("keys/hudumaapp-firebase-adminsdk-jtet8-7370576c3f.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("error creating firebase app %s: %v\n", uid, err)
	}
	// Get an auth client from the firebase.App
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	user, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Printf("error getting firebase user %s: %v\n", uid, err)
	}
	return user
}
*/
func (s *Server) handleRequestCreate(w http.ResponseWriter, r *http.Request) {
	var request model.Request
	request.ID = uuid.New()
	request.Status = "bidding"

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	request.ClientID = userID.String()

	var photos []string
	photoData := r.PostFormValue("photos")
	if photoData != "" {
		err = json.Unmarshal([]byte(photoData), &photos)
		if err != nil {
			handleError(w, "photos: invalid json array value", http.StatusBadRequest)
			return
		}
		request.Photos = photos
	}

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("Marshall error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &request); err != nil {
		log.Printf("Unmarshal error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	// if start_date is not set, infer as urgent and set StartDate to 24 hours
	if request.StartDate == "" {
		if err == nil {
			request.Urgent = true
			request.StartDate = time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05")
		}
	} else {
		// validate start date
		time, err := dateparse.ParseStrict(request.StartDate)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			handleError(w, "start_date: invalid date format", http.StatusBadRequest)
			return
		}

		request.StartDate = time.Format("2006-01-02T15:04:05")
	}

	err = request.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.ReqSvc.CreateRequest(r.Context(), &request)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	handleSuccessMsgWithRes(w, "Request created successfully", request)
}

func (s *Server) handleRequestList(w http.ResponseWriter, r *http.Request) {
	var err error
	var requests []app.Request

	userId, err := middlewares.UserIDFromContext(r.Context())
	if err != nil {
		handleUnathorised(w)
		return
	}

	filter := model.RequestFilter{
		ClientID: userId.String(),
	}

	filter.Status = r.URL.Query().Get("status")

	requests, err = s.ReqSvc.FilterRequests(r.Context(), filter)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, requests)
}

func (s *Server) handleAllRequests(w http.ResponseWriter, r *http.Request) {

	filter := model.RequestFilter{
		Category: r.URL.Query().Get("category"),
		Distance: r.URL.Query().Get("distance"),
	}

	userId, err := middlewares.UserIDFromContext(r.Context())
	if err == nil {
		filter.UserID = userId.String()
	}

	requests, err := s.ReqSvc.AllRequests(r.Context(), filter)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, requests)
}

func (s *Server) handleRecommendedRequests(w http.ResponseWriter, r *http.Request) {

	filter := model.RequestFilter{}

	userId, err := middlewares.UserIDFromContext(r.Context())
	if err == nil {
		filter.UserID = userId.String()
	}

	requests, err := s.ReqSvc.AllRequests(r.Context(), filter)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, requests)
}

func (s *Server) handleRequestCategories(w http.ResponseWriter, r *http.Request) {

	requests, err := s.ReqSvc.ListRequestsCategories(r.Context())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, requests)
}

func (s *Server) handleRequestInstantSearch(w http.ResponseWriter, r *http.Request) {
	results := make([]app.RequestSearchResult, 0)
	var err error

	if r.URL.Query().Get("q") != "" {
		search := model.Search{
			Query:     r.URL.Query().Get("q"),
			Latitude:  r.URL.Query().Get("latitude"),
			Longitude: r.URL.Query().Get("longitude"),
			Distance:  r.URL.Query().Get("distance"),
		}

		// validate
		err = search.Validate()
		if err != nil {
			handleError(w, err.Error(), http.StatusBadRequest)
			return
		}
		results, err = s.SrchSvc.InstantSearchRequests(r.Context(), search)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	handleSuccess(w, results)
}

func (s *Server) handleRequestSearch(w http.ResponseWriter, r *http.Request) {
	results := make([]app.RequestSearchResult, 0)
	var err error

	if r.URL.Query().Get("q") != "" {
		search := model.Search{
			Query:     r.URL.Query().Get("q"),
			Latitude:  r.URL.Query().Get("latitude"),
			Longitude: r.URL.Query().Get("longitude"),
			Distance:  r.URL.Query().Get("distance"),
		}

		// validate
		err = search.Validate()
		if err != nil {
			handleError(w, err.Error(), http.StatusBadRequest)
			return
		}
		results, err = s.SrchSvc.InstantSearchRequests(r.Context(), search)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	handleSuccess(w, results)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	results := make([]app.SearchResult, 0)
	var err error

	if r.URL.Query().Get("q") != "" {
		search := model.Search{
			Query:     r.URL.Query().Get("q"),
			Latitude:  r.URL.Query().Get("latitude"),
			Longitude: r.URL.Query().Get("longitude"),
			Distance:  r.URL.Query().Get("distance"),
		}

		// validate
		err = search.Validate()
		if err != nil {
			handleError(w, err.Error(), http.StatusBadRequest)
			return
		}
		results, err = s.SrchSvc.SearchByQuery(r.Context(), search)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	handleSuccess(w, results)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, "Id is not a valid UUID", http.StatusBadRequest)
		return
	}

	request, err := s.ReqSvc.FindRequestByID(r.Context(), id)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "Request not found", http.StatusNotFound)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//request.Status = "pending"

	handleSuccess(w, request)
}

func (s *Server) handleBidCreate(w http.ResponseWriter, r *http.Request) {
	var bid model.Bid

	userId, err := middlewares.UserIDFromContext(r.Context())
	if err != nil {
		handleUnathorised(w)
		return
	}

	bid.BidderID = userId.String()

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &bid); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	if err := bid.Validate(); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.BidSvc.CreateBid(r.Context(), &bid)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Booking created successfully", bid)
}

func (s *Server) handleMyBids(w http.ResponseWriter, r *http.Request) {
	userId, err := middlewares.UserIDFromContext(r.Context())
	if err != nil {
		handleUnathorised(w)
		return
	}

	resp, err := s.BidSvc.ListMyBids(r.Context(), userId.String())
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}

func (s *Server) handleRequestBids(w http.ResponseWriter, r *http.Request) {
	requestId := mux.Vars(r)["id"]

	userId, err := middlewares.UserIDFromContext(r.Context())
	if err != nil {
		handleUnathorised(w)
		return
	}

	resp, err := s.BidSvc.FindBidsByRequestID(r.Context(), userId.String(), requestId)
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}

func (s *Server) handleAcceptBid(w http.ResponseWriter, r *http.Request) {
	bidId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, "Id is not a valid integer", http.StatusBadRequest)
		return
	}

	// TODO: check if bid exists

	err = s.BidSvc.AcceptBid(r.Context(), bidId)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsg(w, "Bid accepted successfully")
}

func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	startTime := r.PostFormValue("start_time")
	time, err := dateparse.ParseStrict(startTime)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "start_date: invalid date format", http.StatusBadRequest)
		return
	}
	startTime = time.String()
	err = s.BkSvc.InsertDate(r.Context(), startTime)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, time.Format("2021-12-11 14:12:03"))
}
