package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"

	"github.com/iamsayantan/MessageRooms/user"
)

// Middleware provides generic interface for registering different middlewares.
type Middleware interface {
	Register(http.Handler) http.Handler
}

// ====================================================//
//               Auth Middleware.                      //
//=====================================================//

// contextKey type is used hold values to http context.
type contextKey int

const (
	// keyAuthUser holds the currently authenticatd user to context.
	keyAuthUser contextKey = 0

	// authorizationHeader is the key from where we extract the authentication token.
	authorizationHeader = "Authorization"
)

// ErrUnauthorized error for authentication failure.
var ErrUnauthorized = errors.New("Unauthorized")

type authMiddleware struct {
	us user.Service
}

func (am *authMiddleware) Register(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(authorizationHeader)
		ctx := r.Context()

		if token == "" {
			render.Render(w, r, ErrUnAuthorized(ErrUnauthorized))
			return
		}

		user, err := am.us.VerifyAuthToken(token)
		if err != nil {
			render.Render(w, r, ErrUnAuthorized(err))
			return
		}

		// set the user in request context.
		ctx = context.WithValue(ctx, keyAuthUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func newAuthMiddleware(us user.Service) Middleware {
	return &authMiddleware{us: us}
}
