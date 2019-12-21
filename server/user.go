package server

import (
	"errors"
	"net/http"

	"gopkg.in/go-playground/validator.v10"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	"github.com/iamsayantan/MessageRooms/user"
)

type authRequest struct {
	Nickname string `json:"nickname" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type userHandler struct {
	service  user.Service
	validate *validator.Validate
}

func (h *userHandler) Route() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.login)
	r.Post("/register", h.register)

	return r
}

func (h *userHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginReq authRequest

	err := decodeJSONBody(w, r, &loginReq)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			render.Render(w, r, ErrInvalidRequest(mr))
		} else {
			render.Render(w, r, ErrInternalServer(err))
		}
		return
	}

	usr, er := h.service.Login(loginReq.Nickname, loginReq.Password)
	if er != nil {
		render.Render(w, r, ErrInvalidRequest(er))
		return
	}

	sendResponse(w, http.StatusOK, usr)
}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var registerRequest authRequest
	err := decodeJSONBody(w, r, &registerRequest)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			render.Render(w, r, ErrInvalidRequest(mr))
		} else {
			render.Render(w, r, ErrInternalServer(err))
		}
		return
	}

	usr, er := h.service.NewUser(registerRequest.Nickname, registerRequest.Password)
	if er != nil {
		render.Render(w, r, ErrInvalidRequest(er))
		return
	}

	sendResponse(w, http.StatusCreated, usr)
}

// NewUserHandler returns new user handler.
func NewUserHandler(s user.Service, v *validator.Validate) WebHandler {
	h := &userHandler{service: s, validate: v}
	return h
}
