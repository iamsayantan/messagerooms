package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	messagerooms "github.com/iamsayantan/MessageRooms"
)

type SSEHub struct {
	NewConnection   chan messagerooms.EventsourceConnection       // NewConnection is the channel for any new client connection
	CloseConnection chan messagerooms.EventsourceConnection       // CloseConnection is channel for any closing connection
	OpenConnections map[string]messagerooms.EventsourceConnection // OpenConnections holds all the active open connections to the server
}

func (s *Server) HandleSSE(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Make sure the connection supports SSE
	flusher, ok := w.(http.Flusher)
	if !ok {
		_ = render.Render(w, r, ErrInternalServer(errors.New("streaming unsupported")))
		return
	}

	// We need to verify the connecting user has authorization to connect to the this
	authUser, ok := ctx.Value(KeyAuthUser).(*messagerooms.User)
	if !ok {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("could not get user")))
		return
	}

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for now allowing cross origin requests.
}
