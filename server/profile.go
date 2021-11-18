package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/go-sql-driver/mysql"
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

	profile, err := s.UsrSvc.FindProfileByUserID(ctx, userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "Profile not found", http.StatusNotFound)
			return
		}
		handleError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	profile.Phone = middlewares.PhoneFromContext(ctx)

	handleSuccess(w, profile)
}

func (s *Server) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	var profile model.Profile
	ctx := r.Context()

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &profile); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	// Check if user is currently logged in.
	userID, err := middlewares.UserIDFromContext(ctx)
	if err != nil {
		handleError(w, "unauthorized", http.StatusUnauthorized)
		handleUnathorised(w)
		return
	}

	// retrieve existing profile from database
	_, err = s.UsrSvc.FindProfileByUserID(ctx, userID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("profile does not exist")
			handleError(w, "profile does not exist", http.StatusNotFound)
			//s.handleProfileCreate(w, r)
			return
		}
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	//profile.ID = currProfile.ID
	//profile.UserID = currProfile.UserID
	profile.UserID = userID.String()
	/*
		err = updateFirebaseUserData(ctx, profile)
		if err != nil {
			log.Printf("error updating firebase user %s: %v\n", userID, err)
			handleError(w, "error updating firebase user data", http.StatusInternalServerError)
			return
		}
	*/
	err = s.UsrSvc.UpdateProfile(ctx, &profile)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)

		if err == sql.ErrNoRows {
			handleError(w, "profile update failed", http.StatusNotFound)
			return
		}
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Profile updated successfuly", profile)
	//json.NewEncoder(w).Encode(profile)
}

func (s *Server) handleProfileCreate(w http.ResponseWriter, r *http.Request) {
	var profile model.Profile

	ctx := r.Context()
	userID, err := middlewares.UserIDFromContext(ctx)
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}

	//userData := RetrieveFirebaseUserData(ctx, userID.String())

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &profile); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	profile.ID = uuid.New()
	profile.UserID = userID.String()
	profile.Type = "client"

	err = s.UsrSvc.CreateProfile(r.Context(), &profile)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleDuplicateEntry(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	handleSuccessMsgWithRes(w, "Profile created successfuly", profile)
}

// Proivder

func (s *Server) handleProviderByID(w http.ResponseWriter, r *http.Request) {
	providerId := mux.Vars(r)["id"]

	provider, err := s.UsrSvc.FindProviderByID(r.Context(), providerId)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "profile not found", http.StatusNotFound)
			return
		}
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	//saveFirebaseUserToProfile(r.Context(), &provider.Profile)

	handleSuccess(w, provider)
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

	jsonStr, err := json.Marshal(allFormValues(r))
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing form values", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(jsonStr, &provider); err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "error parsing json string", http.StatusInternalServerError)
		return
	}

	provider.UserID = userID.String()
	provider.Type = "provider"

	err = s.UsrSvc.CreateProvider(r.Context(), &provider)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleDuplicateEntry(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	handleSuccessMsgWithRes(w, "Provider created successfuly", provider)
}

func handleDuplicateEntry(w http.ResponseWriter, err error) error {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}
	if me.Number == 1062 {
		handleError(w, "record already exists", http.StatusConflict)
	}
	return nil
}

func (s *Server) handleProviderList(w http.ResponseWriter, r *http.Request) {
	providers, err := s.UsrSvc.ListProviders(r.Context())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, providers)
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
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, user)
}

/*
func updateFirebaseUserData(ctx context.Context, profile model.Profile) error {
	var update bool
	params := &auth.UserToUpdate{}

	if profile.DisplayName != nil {
		params.DisplayName(*profile.DisplayName)
		update = true
	}
	if profile.PhotoUrl != nil {
		params.PhotoURL(*profile.PhotoUrl)
		update = true
	}

	if !update {
		return nil
	}

	uid := profile.UserID
	opt := option.WithCredentialsFile("keys/hudumaapp-firebase-adminsdk-jtet8-7370576c3f.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("error creating firebase app %s: %v\n", uid, err)
		return err
	}
	// Get an auth client from the firebase.App
	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		return err
	}

	_, err = client.UpdateUser(ctx, uid, params)
	if err != nil {
		log.Printf("error updating firebase user %s: %v\n", uid, err)
	}
	return err
}
*/
func saveFirebaseUserToProfile(ctx context.Context, profile *app.Profile) {
	firebaseUser := retrieveFirebaseUserData(ctx, profile.UserID)
	if firebaseUser != nil {
		profile.Email = strOrNil(firebaseUser.Email)
		profile.Phone = strOrNil(firebaseUser.PhoneNumber)
		profile.Username = strOrNil(firebaseUser.DisplayName)
		profile.PhotoUrl = strOrNil(firebaseUser.PhotoURL)
		profile.EmailVerified = firebaseUser.EmailVerified
	}
}
