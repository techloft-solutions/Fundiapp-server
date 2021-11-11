package server

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	app "github.com/andrwkng/hudumaapp"
	"github.com/andrwkng/hudumaapp/model"
	"github.com/andrwkng/hudumaapp/server/middlewares"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
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
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	saveFirebaseUserToProfile(ctx, profile)

	handleSuccess(w, profile)
}

func (s *Server) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	var profile model.Profile
	ctx := r.Context()
	userID, err := middlewares.UserIDFromContext(r.Context())
	// Return an error if the user is not currently logged in.
	if err != nil {
		handleUnathorised(w)
		return
	}
	profile.UserID = userID.String()

	currProfilePtr, err := s.UsrSvc.GetProfile(ctx, userID.String())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	currProfile := *currProfilePtr
	profile.FirstName = strOrNil(NewOrCurr(r.PostFormValue("first_name"), *currProfile.FirstName))
	profile.LastName = strOrNil(NewOrCurr(r.PostFormValue("last_name"), *currProfile.LastName))
	//email := PtrNewOrCurr(r.PostFormValue("email"), currProfile.Email)
	//profile.Email = &email
	//photoUrl := PtrNewOrCurr(r.PostFormValue("photo_url"), currProfile.PhotoUrl)
	//profile.PhotoUrl = &photoUrl
	//profile.LocationID = &(PtrNewOrCurr(r.PostFormValue("location_id"), &currProfile.Location.ID))
	fBuserData := *retrieveFirebaseUserData(ctx, userID.String())
	profile.DisplayName = strOrNil(NewOrCurr(r.PostFormValue("display_name"), fBuserData.DisplayName))
	profile.PhotoUrl = strOrNil(NewOrCurr(r.PostFormValue("photo_url"), fBuserData.PhotoURL))
	profile.Email = strOrNil(NewOrCurr(r.PostFormValue("email"), fBuserData.Email))
	/*
		err = s.UsrSvc.UpdateProfile(ctx, &profile)
		if err != nil {
			log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}*/

	updateFirebaseUserData(ctx, profile)

	handleSuccess(w, profile)
	//handleSuccessText(w, profile.ID)
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
	//profile.Email = r.PostFormValue("email")
	//profile.Phone = r.PostFormValue("phone")
	profile.PhotoUrl = strOrNil(r.PostFormValue("photo_url"))
	profile.LocationID = strOrNil(r.PostFormValue("location_id"))
	//profile.Status = strOrNil(r.PostFormValue("status"))
	profile.Type = "client"

	err = s.UsrSvc.CreateProfile(r.Context(), &profile)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleDuplicateEntry(w, err); err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
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

	handleSuccessText(w, profile.ID)
}

func (s *Server) handleProviderByID(w http.ResponseWriter, r *http.Request) {
	providerId := mux.Vars(r)["id"]

	provider, err := s.UsrSvc.FindProviderByID(r.Context(), providerId)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	saveFirebaseUserToProfile(r.Context(), &provider.Profile)

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
	provider.UserID = userID.String()
	provider.FirstName = strOrNil(r.PostFormValue("first_name"))
	provider.LastName = strOrNil(r.PostFormValue("last_name"))
	provider.Email = strOrNil(r.PostFormValue("email"))
	provider.PhotoUrl = strOrNil(r.PostFormValue("photo_url"))
	provider.Bio = strOrNil(r.PostFormValue("bio"))
	provider.LocationID = strOrNil(r.PostFormValue("location_id"))
	provider.Profession = strOrNil(r.PostFormValue("profession"))
	//profile.Status = strOrNil(r.PostFormValue("status"))
	provider.Type = "provider"

	err = s.UsrSvc.CreateProfile(r.Context(), &provider.Profile)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleDuplicateEntry(w, err); err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	err = s.UsrSvc.CreateProvider(r.Context(), &provider)
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		if err = handleDuplicateEntry(w, err); err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}

	handleSuccessText(w, provider.ID)
}

func handleDuplicateEntry(w http.ResponseWriter, err error) error {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}
	if me.Number == 1062 {
		http.Error(w, "record already exists", http.StatusConflict)
	}
	return nil
}

func (s *Server) handleProviderList(w http.ResponseWriter, r *http.Request) {
	providers, err := s.UsrSvc.ListProviders(r.Context())
	if err != nil {
		log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
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
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, user)
}

func updateFirebaseUserData(ctx context.Context, profile model.Profile) {
	uid := profile.UserID
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

	params := (&auth.UserToUpdate{}).
		Email(strOrNull(profile.Email)).
		DisplayName(strOrNull(profile.DisplayName)).
		PhotoURL(strOrNull(profile.PhotoUrl))

	_, err = client.UpdateUser(ctx, uid, params)
	if err != nil {
		log.Printf("error updating firebase user %s: %v\n", uid, err)
	}
}

func saveFirebaseUserToProfile(ctx context.Context, profile *app.Profile) {
	firebaseUser := retrieveFirebaseUserData(ctx, profile.UserID)
	if firebaseUser != nil {
		profile.Email = strOrNil(firebaseUser.Email)
		profile.Phone = strOrNil(firebaseUser.PhoneNumber)
		profile.DisplayName = strOrNil(firebaseUser.DisplayName)
		profile.PhotoUrl = strOrNil(firebaseUser.PhotoURL)
		profile.EmailVerified = firebaseUser.EmailVerified
	}
}
