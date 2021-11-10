package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) handleProfileGet(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}
	ctx := r.Context()

	profile, err := s.UsrSvc.GetProfile(ctx, userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
		return
	}
	saveFirebaseUserToProfile(ctx, profile)

	handleSuccess(w, profile)
}

func (s *Server) handleProfileByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	user, err := s.UsrSvc.FindProfileByUserID(r.Context(), userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
		return
	}

	handleSuccess(w, user)
}

func (s *Server) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	var profile model.Profile
	profile.UserID = userID.String()
	profile.FirstName = strOrNil(r.PostFormValue("first_name"))
	profile.LastName = strOrNil(r.PostFormValue("last_name"))
	profile.Bio = r.PostFormValue("bio")
	profile.Phone = r.PostFormValue("phone")
	profile.Email = r.PostFormValue("email")

	err = s.UsrSvc.UpdateProfile(r.Context(), &profile)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		return
	}

	handleSuccess(w, profile)
}

func (s *Server) handleProfileCreate(w http.ResponseWriter, r *http.Request) {
	var profile model.Profile
	profile.ID = uuid.New()
	ctx := r.Context()
	userID, err := middlewares.UserIDFromContext(ctx)
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	//userData := RetrieveFirebaseUserData(ctx, userID.String())

	profile.UserID = userID.String()
	profile.FirstName = strOrNil(r.PostFormValue("first_name"))
	profile.LastName = strOrNil(r.PostFormValue("last_name"))
	profile.Bio = r.PostFormValue("bio")
	profile.LocationID = r.PostFormValue("location_id")
	profile.Status = r.PostFormValue("status")
	profile.Profession = r.PostFormValue("profession")
	profile.Type = r.PostFormValue("account_type")

	err = profile.Validate()
	if err != nil {
		//handleError(w, err, http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.UsrSvc.CreateProfile(r.Context(), &profile)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
		return
	}

	/*
		res := app.Profile{
			User: app.User{
				UserID:    profile.UserID,
				FirstName: profile.FirstName,
				LastName:  profile.LastName,
			},
			//Email:     profile.Email,
			//Phone:     profile.Phone,
			//Bio:       profile.Bio,
		}*/

	//handleSuccess(w)
}

func (s *Server) handleProviderByID(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	resp, err := s.UsrSvc.FindProviderByID(r.Context(), userId)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
		return
	}

	handleSuccess(w, resp)
}

func (s *Server) handleProviderCreate(w http.ResponseWriter, r *http.Request) {
	var provider model.Provider
	provider.ID = uuid.New()

	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}
	provider.UserID = userID.String()

	err = s.UsrSvc.CreateProvider(r.Context(), &provider)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
		return
	}

	handleSuccess(w, provider)
}

func (s *Server) handleProviderList(w http.ResponseWriter, r *http.Request) {
	providers, err := s.UsrSvc.ListProviders(r.Context())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, errors.New("something went wrong"), http.StatusInternalServerError)
		return
	}

	handleSuccess(w, providers)
}
