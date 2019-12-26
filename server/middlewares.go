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
	// KeyAuthUser holds the currently authenticatd user to context.
	KeyAuthUser contextKey = 0

	// AuthorizationHeader is the key from where we extract the authentication token.
	AuthorizationHeader = "Authorization"
)

// ErrUnauthorized error for authentication failure.
var ErrUnauthorized = errors.New("Unauthorized")

type authMiddleware struct {
	us user.Service
}

func (am *authMiddleware) Register(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthorizationHeader)
		ctx := r.Context()

		if token == "" {
			_ = render.Render(w, r, ErrUnAuthorized(ErrUnauthorized))
			return
		}

		user, err := am.us.VerifyAuthToken(token)
		if err != nil {
			_ = render.Render(w, r, ErrUnAuthorized(err))
			return
		}

		// set the user in request context.
		ctx = context.WithValue(ctx, KeyAuthUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func newAuthMiddleware(us user.Service) Middleware {
	return &authMiddleware{us: us}
}
