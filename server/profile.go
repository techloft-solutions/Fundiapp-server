package server

import (
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
		handleUnathorised(w)
		return
	}

	// retrieve existing profile from database
	_, err = s.UsrSvc.FindProfileByUserID(ctx, userID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("profile not found")
			handleError(w, "profile not found", http.StatusNotFound)
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

func (s *Server) handleProfileLocationUpdate(w http.ResponseWriter, r *http.Request) {
	location := r.FormValue("location_id")
	handleSuccessMsg(w, "Location updated successfuly:"+location)
}

func (s *Server) handleProviderUpdate(w http.ResponseWriter, r *http.Request) {
	var provider model.Provider
	ctx := r.Context()

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

	// Check if user is currently logged in.
	userID, err := middlewares.UserIDFromContext(ctx)
	if err != nil {
		handleUnathorised(w)
		return
	}

	provider.UserID = userID.String()

	// Check if user is a service provider.
	user, err := s.UsrSvc.FindUserByID(ctx, userID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Provider does not exist")
			handleError(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	if !user.IsProvider {
		handleError(w, "User is not a provider", http.StatusUnauthorized)
		return
	}

	err = s.UsrSvc.UpdateProvider(ctx, &provider)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)

		if err == sql.ErrNoRows {
			handleError(w, "provider update failed", http.StatusNotFound)
			return
		}
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccessMsgWithRes(w, "Profile updated successfuly", provider)
}

// Provider

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

func (s *Server) handleProviderGet(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.UserIDFromContext(r.Context())
	if err != nil {
		handleUnathorised(w)
		return
	}

	providers, err := s.UsrSvc.FindProviderByUserID(r.Context(), userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err == sql.ErrNoRows {
			handleError(w, "provider not found", http.StatusNotFound)
			return
		}
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, providers)
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

	err = s.UsrSvc.CreateProvider(r.Context(), &provider)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleMysqlErrors(w, err); err != nil {
			handleError(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	handleSuccessMsgWithRes(w, "Provider created successfuly", provider)
}

func (s *Server) handleProviderReviews(w http.ResponseWriter, r *http.Request) {
	providerId := mux.Vars(r)["id"]

	reviews, err := s.RevSvc.ListReviewsByProviderID(r.Context(), providerId)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, reviews)
}

func (s *Server) handleProviderServices(w http.ResponseWriter, r *http.Request) {
	providerId := mux.Vars(r)["id"]

	services, err := s.UsrSvc.ListServicesByProviderID(r.Context(), providerId)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, services)
}

func (s *Server) handleProviderPortfolios(w http.ResponseWriter, r *http.Request) {
	providerId := mux.Vars(r)["id"]

	portfolios, err := s.PfoSvc.ListPortfoliosByProviderId(r.Context(), providerId)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		handleError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, portfolios)
}

func handleMysqlErrors(w http.ResponseWriter, err error) error {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}
	if me.Number == 1062 {
		handleError(w, "duplicate entry", http.StatusConflict)
		return nil
	}
	if me.Number == 1452 {
		handleError(w, "foreign key constraint fail", http.StatusNotFound)
		return nil
	}
	return err
}

func (s *Server) handleProviderList(w http.ResponseWriter, r *http.Request) {
	var err error
	var providers []*app.ProviderBrief

	if r.URL.Query().Get("filter") == "true" {
		var filter model.ProviderFilter
		filter.CategoryID = r.URL.Query().Get("category_id")
		filter.IndustryID = r.URL.Query().Get("industry_id")

		providers, err = s.UsrSvc.FilterProviders(r.Context(), filter)
	} else {
		providers, err = s.UsrSvc.ListProviders(r.Context())
	}
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
/*
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
*/
