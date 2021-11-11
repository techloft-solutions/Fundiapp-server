package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

func (s *Server) handleBookingByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID from path.
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, app.Errorf(app.INVALID_ERR, "Invalid ID format"), 400)
		return
	}

	resp, err := s.BkSvc.FindBookingByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
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
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
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
			http.Error(w, "photos: must input a valid photo value", http.StatusBadRequest)
			return
		}
	}
	booking.Photos = photos
	booking.Title = r.PostFormValue("title")
	booking.StartDate = r.PostFormValue("start_date")
	booking.Description = strOrNil(r.PostFormValue("description"))
	booking.LocationID = r.PostFormValue("location_id")
	booking.Status = r.PostFormValue("status")
	booking.ProviderID = strOrNil(r.PostFormValue("provider_id"))
	booking.ServiceID = r.PostFormValue("service_id")

	err = booking.Validate()
	if err != nil {
		//handleError(w, err, http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.BkSvc.CreateBooking(r.Context(), &booking)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}

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
			ID: booking.LocationID,
		},
	}

	handleSuccess(w, res)
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

func strOrNil(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func strOrNull(ptr *string) string {
	switch ptr {
	case nil:
		return "null"
	default:
		return *ptr
	}
}

// newOrCurr returns the new value if it is not empty, otherwise the current value.
func PtrNewOrCurr(new string, curr *string) string {
	if new == "" {
		return *curr
	}
	return new
}

func NewOrCurr(new string, curr string) string {
	if new == "" {
		return curr
	}
	return new
}

func ptrToStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (s *Server) handleRequestCreate(w http.ResponseWriter, r *http.Request) {
	var request model.Request
	request.ID = uuid.New()

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	request.ClientID = userID.String()

	var photos []string
	photoData := r.PostFormValue("photos")
	fmt.Println("photodata:", photoData)
	if photoData != "" {
		err = json.Unmarshal([]byte(photoData), &photos)
		if err != nil {
			http.Error(w, "photos: must input a valid photo value", http.StatusBadRequest)
			return
		}
	}
	request.Photos = photos
	request.Title = r.PostFormValue("title")
	request.StartDate = r.PostFormValue("start_date")
	request.Note = r.PostFormValue("note")
	//urgent := r.PostFormValue("urgent")
	request.LocationID = r.PostFormValue("location_id")
	request.Type = "REQUEST"

	err = request.Validate()
	if err != nil {
		//handleError(w, err, http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.ReqSvc.CreateRequest(r.Context(), &request)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleUnathorised(w)
		return
	}

	res := app.RequestDetail{
		ID:      request.ID,
		Title:   request.Title,
		Status:  request.Status,
		Created: request.CreatedAt.String(),
		Start:   request.StartDate,
		Note:    request.Note,
		Photos:  request.Photos,
	}

	handleSuccess(w, res)
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
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
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

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handleError(w, app.Errorf(app.INVALID_ERR, "Invalid ID format"), 400)
		return
	}
	// Fetch dials from database.
	resp, err := s.ReqSvc.FindRequestByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
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
