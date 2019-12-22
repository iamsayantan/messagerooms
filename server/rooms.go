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
	return router
}

func (h *roomHandler) getRoomDetails(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		render.Render(w, r, ErrInvalidRequest(ErrInvalidRoomID))
		return
	}

	authUser, ok := r.Context().Value(keyAuthUser).(*messagerooms.User)

	if !ok {
		render.Render(w, r, ErrInvalidRequest(errors.New("could not get user")))
		return
	}

	room, err := h.service.RoomDetails(roomID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	exists := h.service.CheckUserExistsInRoom(*room, *authUser)
	log.Printf("Exists %v", exists)

	resp := struct {
		User messagerooms.User `json:"user"`
		Room messagerooms.Room `json:"room"`
	}{
		User: *authUser,
		Room: *room,
	}
	sendResponse(w, http.StatusOK, resp)
}

// newRoomHandler returns a new roomHandler instance.
func newRoomHandler(rs room.Service, v *validator.Validate) WebHandler {
	rh := &roomHandler{service: rs, validate: v}
	return rh
}
