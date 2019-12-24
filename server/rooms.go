package server

import (
	"errors"
	"log"
	"net/http"

	messagerooms "github.com/iamsayantan/MessageRooms"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"

	"github.com/iamsayantan/MessageRooms/room"
	"gopkg.in/go-playground/validator.v10"
)

var (
	// ErrInvalidRoomID is returned in case an invalid room id is provided.
	ErrInvalidRoomID = errors.New("the room id must be a valid room id")
)

type roomHandler struct {
	service  room.Service
	validate *validator.Validate
}

func (h *roomHandler) Route() chi.Router {
	router := chi.NewRouter()
	router.Get("/{roomID}", h.getRoomDetails)
	router.Put("/{roomID}/join", h.joinRoom)
	return router
}

func (h *roomHandler) getRoomDetails(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		_ = render.Render(w, r, ErrInvalidRequest(ErrInvalidRoomID))
		return
	}

	authUser, ok := r.Context().Value(keyAuthUser).(*messagerooms.User)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("could not get user")))
		return
	}

	roomDetails, err := h.service.RoomDetails(roomID)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	exists := h.service.CheckUserExistsInRoom(*roomDetails, *authUser)
	log.Printf("Exists %v", exists)

	resp := struct {
		Room     messagerooms.Room `json:"roomDetails"`
		IsMember bool              `json:"is_member"`
	}{
		Room:     *roomDetails,
		IsMember: h.service.CheckUserExistsInRoom(*roomDetails, *authUser),
	}
	sendResponse(w, http.StatusOK, resp)
}

func (h *roomHandler) joinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		_ = render.Render(w, r, ErrInvalidRequest(ErrInvalidRoomID))
		return
	}

	roomDetails, err := h.service.RoomDetails(roomID)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	authUser, ok := r.Context().Value(keyAuthUser).(*messagerooms.User)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("could not get user")))
		return
	}

	if err := h.service.AddUserToRoom(*roomDetails, *authUser); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	resp := struct {
		OK bool `json:"ok"`
	}{OK: true}

	sendResponse(w, http.StatusOK, resp)
}

// newRoomHandler returns a new roomHandler instance.
func newRoomHandler(rs room.Service, v *validator.Validate) WebHandler {
	rh := &roomHandler{service: rs, validate: v}
	return rh
}
