package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
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
	booking, err := s.BkSvc.FindBookings(r.Context())
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, booking)
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

	booking.ClientID = userID.String()
	// validate
	err = booking.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		log.Println("[DEBUG] User not logged in", err)
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
	}

	err = request.Validate()
	if err != nil {
		//handleError(w, err, http.StatusBadRequest)
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.ReqSvc.CreateRequest(r.Context(), &request)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	/*
		res := app.RequestDetail{
			ID:      request.ID,
			Title:   request.Title,
			Status:  request.Status,
			Created: request.CreatedAt.String(),
			Start:   request.StartDate,
			Note:    request.Note,
			Photos:  request.Photos,
		}
	*/
	handleSuccessMsgWithRes(w, "Request created successfully", request)
}

func (s *Server) handleRequestList(w http.ResponseWriter, r *http.Request) {
	userId, err := middlewares.UserIDFromContext(r.Context())
	if err != nil {
		handleUnathorised(w)
		return
	}
	resp, err := s.ReqSvc.ListRequests(r.Context(), userId)
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
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

	handleSuccessMsgWithRes(w, "Service created successfully", bid)
}
