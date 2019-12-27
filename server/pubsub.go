package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/render"
	messagerooms "github.com/iamsayantan/MessageRooms"
)

type SSEHub struct {
	NewConnection   chan messagerooms.EventsourceConnection       // NewConnection is the channel for any new client connection
	CloseConnection chan messagerooms.EventsourceConnection       // CloseConnection is channel for any closing connection
	OpenConnections map[string]messagerooms.EventsourceConnection // OpenConnections holds all the active open connections to the server
}

// HandleSSE handles incoming persistent connection.
func (s *SSEHub) HandleSSE(w http.ResponseWriter, r *http.Request) {

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

	eventSourceConn := messagerooms.NewEventsourceConnection(authUser)

	// Signal the SSEHub that we have a new client connection.
	s.NewConnection <- *eventSourceConn

	// We need to notify the hub if somehow the connection dies and the handler exits.
	defer func() {
		log.Printf("Connection interrupted: %s", eventSourceConn.ConnectionID)
		s.CloseConnection <- *eventSourceConn
	}()

	// If the client closes the connection, we receive that in ctx.Done() channel for which we need to listen
	// and inform the hub about the event.
	go func() {
		log.Printf("Connection closed by client: %s", eventSourceConn.ConnectionID)
		<-ctx.Done()
		s.CloseConnection <- *eventSourceConn
	}()

	// block waiting for messages broadcast on this connection's SendCh
	for {
		evt := <-eventSourceConn.SendCh
		log.Printf("Got event for connectionID: %s Event:\n%s", evt.DestinationID, evt.String())
		_, _ = fmt.Fprintf(w, evt.String())
		flusher.Flush()
	}
}

// Listen spawns a goroutine that listens for any incoming or closing client connections.
func (s *SSEHub) Listen() {
	go func() {
		for {
			select {
			case sseConn := <-s.NewConnection:
				s.OpenConnections[sseConn.ConnectionID] = sseConn

				// send an initial event with the connection id
				msg := messagerooms.EventMessage{Event: messagerooms.ConnectionEvent, DestinationID: sseConn.ConnectionID, Data: fmt.Sprintf("connectionID: %s", sseConn.ConnectionID)}
				sseConn.PublishEvent(msg)
				sseConn.Heartbeat()

				log.Printf("New client connected. ConnectionID: %s Number of registered clients %d", sseConn.ConnectionID, len(s.OpenConnections))
			case sseConn := <-s.CloseConnection:
				sseConn.Closing()
				delete(s.OpenConnections, sseConn.ConnectionID)
				log.Printf("Removed client. ConnectionID %s Number of registered clients %d", sseConn.ConnectionID, len(s.OpenConnections))
			}
		}
	}()
}

// NewSSEHub returns a new hub instance.
func NewSSEHub() *SSEHub {
	sseHub := SSEHub{
		NewConnection:   make(chan messagerooms.EventsourceConnection),
		CloseConnection: make(chan messagerooms.EventsourceConnection),
		OpenConnections: make(map[string]messagerooms.EventsourceConnection),
	}
	return &sseHub
}
