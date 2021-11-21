package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
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
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (s *Server) handleBookingList(w http.ResponseWriter, r *http.Request) {
	// Fetch dials from database.
	resp, err := s.BkSvc.FindBookings(r.Context())
	if err != nil {
		log.Println(err)
		handleError(w, "Something went wrong", http.StatusInternalServerError)
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

	var photos []string
	photoData := r.PostFormValue("photos")
	fmt.Println("photodata:", photoData)
	if photoData != "" {
		err = json.Unmarshal([]byte(photoData), &photos)
		if err != nil {
			handleError(w, "photos: must input a valid photo value", http.StatusBadRequest)
			return
		}
	}
	/*
		booking.Photos = photos
		booking.Title = r.PostFormValue("title")
		booking.StartDate = r.PostFormValue("start_date")
		booking.Description = strOrNil(r.PostFormValue("description"))
		booking.LocationID = r.PostFormValue("location_id")
		booking.Status = r.PostFormValue("status")
		booking.ProviderID = strOrNil(r.PostFormValue("provider_id"))
		booking.ServiceID = r.PostFormValue("service_id")
	*/

	err = booking.Validate()
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*
		err = s.BkSvc.CreateBooking(r.Context(), &booking)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			handleError(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	*/
	/*
		res := app.Booking{
			ID:          booking.ID,
			Title:       booking.Title,
			Status:      booking.Status,
			BookedAt:    booking.CreatedAt.String(),
			StartAt:     booking.StartDate,
			Description: booking.Description,
			Photos:      booking.Photos,
			Client:      app.Client{},
			Location: app.Location{
				ID: &booking.LocationID,
			},
		}

		handleSuccess(w, res)
	*/
	handleSuccessMsgWithRes(w, "Booking:", booking)
}

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
	/*
		var photos []string
		photoData := r.PostFormValue("photos")
		fmt.Println("photodata:", photoData)
		if photoData != "" {
			err = json.Unmarshal([]byte(photoData), &photos)
			if err != nil {
				handleError(w, "photos: must input a valid photo value", http.StatusBadRequest)
				return
			}
		}
		request.Photos = photos

		request.Title = r.PostFormValue("title")
		request.StartDate = r.PostFormValue("start_date")
		request.Note = r.PostFormValue("note")
		//urgent := r.PostFormValue("urgent")
		request.LocationID = r.PostFormValue("location_id")
	*/
	request.Type = "REQUEST"

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &request); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	urgent, err := govalidator.ToBoolean(r.PostFormValue("urgent"))
	if err == nil {
		request.Urgent = urgent
	}

	if request.Urgent {
		request.StartDate = time.Now().Add(time.Hour * 24).Format(time.RFC3339)
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
		handleUnathorised(w)
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
	// Fetch dials from database.
	request, err := s.ReqSvc.FindRequestByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			handleError(w, "Profile not found", http.StatusNotFound)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	request.Status = "pending"

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
