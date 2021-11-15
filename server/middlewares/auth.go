package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	app "github.com/andrwkng/hudumaapp"
	"google.golang.org/api/option"
)

type accessKey int

const accessKeyAuthToken accessKey = iota

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := parseAuthorizationHeader(r.Context(), r.Header.Get("Authorization"))
		if err != nil {
			//handleError(w, "Invalid authorization token", http.StatusUnauthorized)
			handleInvalidToken(w)
			return
		}

		user := app.AuthUser{
			ID: token.UID,
			//Name: token.Claims["name"].(*string),
			//Email: token.Claims["email"].(string),
		}

		r = r.WithContext(context.WithValue(r.Context(), accessKeyAuthToken, user))
		next.ServeHTTP(w, r)
	},
	)

}

func parseAuthorizationHeader(ctx context.Context, tokenHeader string) (*auth.Token, error) {
	idToken := strings.Replace(tokenHeader, "Bearer ", "", 1)

	if idToken == "" {
		return nil, fmt.Errorf("token not set")
	}
	opt := option.WithCredentialsFile("keys/hudumaapp-firebase-adminsdk-jtet8-7370576c3f.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("error creating firebase app:", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}

	return token, nil
}

func handleInvalidToken(w http.ResponseWriter) {
	resp := make(map[string]string)
	resp["error"] = "Invalid Authorization token"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(jsonResp)
}
