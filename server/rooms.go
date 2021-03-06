package server

import (
	"errors"
	"net/http"

	"github.com/iamsayantan/messagerooms"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"

	"github.com/iamsayantan/messagerooms/room"
)

var (
	// ErrInvalidRoomID is returned in case an invalid room id is provided.
	ErrInvalidRoomID = errors.New("the room id must be a valid room id")

	// ErrMessageTextEmpty is returned when user posts messgae with blank string
	ErrMessageTextEmpty = errors.New("message text should not be empty")

	// ErrRoomNameEmpty is returned when user tries to create a room with empty room name
	ErrRoomNameEmpty = errors.New("room name can not be empty")
)

// newMessageRequest request payload for posting new messages.
type newMessageRequest struct {
	MessageText string `json:"message_text"`
}

// createRoomRequest creates a new room.
type createRoomRequest struct {
	RoomName string `json:"room_name"`
}

type roomHandler struct {
	service room.Service
}

func (h *roomHandler) Route() chi.Router {
	router := chi.NewRouter()
	router.Get("/", h.allRooms)
	router.Post("/create", h.createRoom)
	router.Get("/{roomID}", h.getRoomDetails)
	router.Put("/{roomID}/join", h.joinRoom)
	router.Get("/{roomID}/messages", h.getAllMessages)
	router.Post("/{roomID}/messages", h.postMessage)
	return router
}

func (h *roomHandler) createRoom(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := r.Context().Value(KeyAuthUser).(*messagerooms.User)
	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("could not get user")))
		return
	}

	var req createRoomRequest
	err := decodeJSONBody(w, r, &req)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			_ = render.Render(w, r, ErrInvalidRequest(mr))
		} else {
			_ = render.Render(w, r, ErrInternalServer(err))
		}
		return
	}

	if req.RoomName == "" {
		_ = render.Render(w, r, ErrInvalidRequest(ErrRoomNameEmpty))
		return
	}

	createdRoom, err := h.service.CreateNewRoom(req.RoomName, *authenticatedUser)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	resp := struct {
		Room messagerooms.Room `json:"room"`
	}{Room: *createdRoom}
	sendResponse(w, http.StatusOK, resp)

}

func (h *roomHandler) allRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.AllRooms()
	if err != nil {
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}

	resp := struct {
		Rooms []*messagerooms.Room `json:"rooms"`
	}{Rooms: rooms}

	sendResponse(w, http.StatusOK, resp)
}

func (h *roomHandler) getRoomDetails(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		_ = render.Render(w, r, ErrInvalidRequest(ErrInvalidRoomID))
		return
	}

	authUser, ok := r.Context().Value(KeyAuthUser).(*messagerooms.User)

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

	resp := struct {
		Room     messagerooms.Room `json:"room_details"`
		IsMember bool              `json:"is_member"`
	}{
		Room:     *roomDetails,
		IsMember: exists,
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

	authUser, ok := r.Context().Value(KeyAuthUser).(*messagerooms.User)

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

	authUser, ok := r.Context().Value(KeyAuthUser).(*messagerooms.User)

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

func (h *roomHandler) getAllMessages(w http.ResponseWriter, r *http.Request) {
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

	authUser, ok := r.Context().Value(KeyAuthUser).(*messagerooms.User)

	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("could not get user")))
		return
	}

	if exists := h.service.CheckUserExistsInRoom(*roomDetails, *authUser); !exists {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("you do not have permission to view messages in the room")))
		return
	}

	messages, err := h.service.GetAllRoomMessages(*roomDetails)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	resp := struct {
		Messages []*messagerooms.Message `json:"messages"`
	}{Messages: messages}

	sendResponse(w, http.StatusOK, resp)
}

// newRoomHandler returns a new roomHandler instance.
func newRoomHandler(rs room.Service) WebHandler {
	rh := &roomHandler{service: rs}
	return rh
}
