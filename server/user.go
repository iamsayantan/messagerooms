package server

import (
	"encoding/json"
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

func (h *userHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginReq loginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		encodeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, er := h.service.Login(loginReq.Nickname, loginReq.Password)
	if er != nil {
		encodeError(w, http.StatusBadRequest, er.Error())
		return
	}

}
