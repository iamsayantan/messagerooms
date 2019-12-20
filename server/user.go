package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/iamsayantan/MessageRooms/user"
	"net/http"
)

type loginRequest struct {
	Nickname string
	Password string
}

type userHandler struct {
	service user.Service
}

func (h *userHandler) Route() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.login)

	return r
}

func (h *userHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginReq loginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		encodeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	usr, er := h.service.Login(loginReq.Nickname, loginReq.Password)
	if er != nil {
		encodeError(w, http.StatusBadRequest, er.Error())
		return
	}

	sendResponse(w, http.StatusOK, usr)
}

// NewUserHandler returns new user handler.
func NewUserHandler(s user.Service) WebHandler {
	h := &userHandler{service: s}
	return h
}
