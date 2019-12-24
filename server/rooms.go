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

	// ErrMessageTextEmpty is returned when user posts messgae with blank string
	ErrMessageTextEmpty = errors.New("message text should not be empty")
)

// newMessageRequest request payload for posting new messages.
type newMessageRequest struct {
	MessageText string `json:"message_text"`
}

type roomHandler struct {
	service  room.Service
	validate *validator.Validate
}

func (h *roomHandler) Route() chi.Router {
	router := chi.NewRouter()
	router.Get("/{roomID}", h.getRoomDetails)
	router.Put("/{roomID}/join", h.joinRoom)
	router.Post("/{roomID}/message", h.postMessage)
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

func (h *roomHandler) postMessage(w http.ResponseWriter, r *http.Request) {
	var messageReq newMessageRequest

	err := decodeJSONBody(w, r, &messageReq)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			_ = render.Render(w, r, ErrInvalidRequest(mr))
		} else {
			_ = render.Render(w, r, ErrInternalServer(err))
		}
		return
	}

	if messageReq.MessageText == "" {
		_ = render.Render(w, r, ErrInvalidRequest(ErrMessageTextEmpty))
		return
	}

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

	msg, err := h.service.PostMessage(*roomDetails, *authUser, messageReq.MessageText)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	resp := struct {
		Message messagerooms.Message `json:"message"`
	}{
		Message: *msg,
	}

	sendResponse(w, http.StatusOK, resp)
}

// newRoomHandler returns a new roomHandler instance.
func newRoomHandler(rs room.Service, v *validator.Validate) WebHandler {
	rh := &roomHandler{service: rs, validate: v}
	return rh
}
