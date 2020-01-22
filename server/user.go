package server

import (
	"errors"
	"net/http"

	"github.com/iamsayantan/messagerooms"
	"github.com/iamsayantan/messagerooms/user"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
)

type authRequest struct {
	Nickname string `json:"nickname" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type userHandler struct {
	authMiddleware Middleware
	service        user.Service
}

func (h *userHandler) Route() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.login)
	r.Post("/register", h.register)

	r.Group(func(r chi.Router) {
		r.Use(h.authMiddleware.Register)
		r.Get("/me", h.me)
	})

	return r
}

func (h *userHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginReq authRequest

	err := decodeJSONBody(w, r, &loginReq)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			_ = render.Render(w, r, ErrInvalidRequest(mr))
		} else {
			_ = render.Render(w, r, ErrInternalServer(err))
		}
		return
	}

	usr, er := h.service.Login(loginReq.Nickname, loginReq.Password)
	if er != nil {
		// if error is invalid password, then just returning
		if er == user.ErrInvalidPassword {
			_ = render.Render(w, r, ErrInvalidRequest(er))
			return
		} else if er == user.ErrInvalidNickname {
			// so if user is trying to login with a nickname with does not exist then we will create a new
			// record with the given details and log the user in.
			usr, er = h.service.NewUser(loginReq.Nickname, loginReq.Password)
			if er != nil {
				_ = render.Render(w, r, ErrInvalidRequest(er))
				return
			}
		} else {
			_ = render.Render(w, r, ErrInternalServer(er))
		}
	}

	token, _ := h.service.GenerateAuthToken(*usr)
	resp := struct {
		User        *messagerooms.User `json:"user"`
		AccessToken string             `json:"access_token"`
	}{
		User:        usr,
		AccessToken: token,
	}

	sendResponse(w, http.StatusOK, resp)
}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var registerRequest authRequest
	err := decodeJSONBody(w, r, &registerRequest)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			_ = render.Render(w, r, ErrInvalidRequest(mr))
		} else {
			_ = render.Render(w, r, ErrInternalServer(err))
		}
		return
	}

	usr, er := h.service.NewUser(registerRequest.Nickname, registerRequest.Password)
	if er != nil {
		_ = render.Render(w, r, ErrInvalidRequest(er))
		return
	}

	sendResponse(w, http.StatusCreated, usr)
}

func (h *userHandler) me(w http.ResponseWriter, r *http.Request) {
	authUser, ok := r.Context().Value(KeyAuthUser).(*messagerooms.User)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("could not get user")))
		return
	}

	resp := struct {
		User *messagerooms.User `json:"user"`
	}{User: authUser}

	sendResponse(w, http.StatusOK, resp)
}

// NewUserHandler returns new user handler.
func NewUserHandler(s user.Service, am Middleware) WebHandler {
	h := &userHandler{service: s, authMiddleware: am}
	return h
}
