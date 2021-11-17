package middlewares

import (
	"context"
	"log"

	app "github.com/andrwkng/hudumaapp"
)

// UserFromContext returns the current logged in user.
func UserFromContext(ctx context.Context) *app.AuthUser {
	user, ok := ctx.Value(accessKeyAuthToken).(app.AuthUser)
	if !ok {
		log.Println("Failed getting user from context", ctx.Value(accessKeyAuthToken), user)
		return nil
	}
	return &user
}

// UserIDFromContext is a helper function that returns the ID of the current
// logged in user. Returns zero if no user is logged in.
func UserIDFromContext(ctx context.Context) (app.UserID, error) {
	if user := UserFromContext(ctx); user != nil {
		return app.UserID(user.ID), nil
	}
	return app.UserIDNil, app.Errorf(app.UNAUTHORIZED_ERR, "Not logged in.")
}

func PhoneFromContext(ctx context.Context) *string {
	if user := UserFromContext(ctx); user != nil {
		return &user.PhoneNumber
	}
	return nil
}
