package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/render"
	"github.com/gomodule/redigo/redis"
	messagerooms "github.com/iamsayantan/MessageRooms"
)

// SSEHub maintains persistent eventsource connection to server.
type SSEHub struct {
	mu              sync.Mutex
	NewConnection   chan messagerooms.EventsourceConnection       // NewConnection is the channel for any new client connection
	CloseConnection chan messagerooms.EventsourceConnection       // CloseConnection is channel for any closing connection
	OpenConnections map[string]messagerooms.EventsourceConnection // OpenConnections holds all the active open connections to the server

	pubsubConn *redis.PubSubConn
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
		<-ctx.Done()
		s.CloseConnection <- *eventSourceConn
		log.Printf("Connection closed by client: %s", eventSourceConn.ConnectionID)
	}()

	// block waiting for messages broadcast on this connection's SendCh
	for {
		evt := <-eventSourceConn.SendCh
		//log.Printf("Got event for connectionID: %s Event:\n%s", evt.DestinationID, evt.String())
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
				s.mu.Lock()
				s.OpenConnections[sseConn.ConnectionID] = sseConn
				s.mu.Unlock()

				// send an initial event with the connection id
				connectionEvt := struct {
					ConnectionID string `json:"connection_id"`
					ServerTime   int64  `json:"server_time"`
				}{
					ConnectionID: sseConn.ConnectionID,
					ServerTime:   time.Now().Unix(),
				}

				msg := messagerooms.EventMessage{Event: messagerooms.ConnectionEvent, DestinationID: sseConn.ConnectionID, Data: connectionEvt}
				sseConn.PublishEvent(msg)
				sseConn.Heartbeat()

				log.Printf("New client connected. ConnectionID: %s Number of registered clients %d", sseConn.ConnectionID, len(s.OpenConnections))
			case sseConn := <-s.CloseConnection:
				s.mu.Lock()
				delete(s.OpenConnections, sseConn.ConnectionID)
				s.mu.Unlock()

				sseConn.Closing()
				log.Printf("Removed client. ConnectionID %s Number of registered clients %d", sseConn.ConnectionID, len(s.OpenConnections))
			}
		}
	}()
}

// ReceiveHubEvents spawns a goroutine which listens
func (s *SSEHub) ReceiveHubEvents() {
	go func() {
		for {
			switch v := s.pubsubConn.Receive().(type) {
			case redis.Message:
				// We expect that data should be of type PublishEvent. Otherwise its an error and we don't process it.
				payload := v.Data
				var eventMessage *messagerooms.PublishEvent
				if err := json.Unmarshal(payload, &eventMessage); err != nil {
					log.Printf("Invalid Event Received")
					break
				}

				log.Printf("[Redis Message] Channel: %s, Message: %s\n", v.Channel, string(v.Data))
			case redis.Subscription:
				log.Printf("[Redis Subscription] Channel: %s, Kind: %s, Count: %d\n", v.Channel, v.Kind, v.Count)
			case error:
				log.Printf("Error pub/sub on connection, delivery has stopped %s\n", v.Error())
				return
			}
		}
	}()
}

// NewSSEHub returns a new hub instance.
func NewSSEHub(conn *redis.PubSubConn) *SSEHub {
	sseHub := &SSEHub{
		NewConnection:   make(chan messagerooms.EventsourceConnection),
		CloseConnection: make(chan messagerooms.EventsourceConnection),
		OpenConnections: make(map[string]messagerooms.EventsourceConnection),
		pubsubConn:      conn,
	}

	// subscribe to the hub channel.
	err := sseHub.pubsubConn.Subscribe(messagerooms.HubChannel)
	if err != nil {
		log.Printf("Could not connect to the hub: %s", err.Error())
		panic(err)
	}

	sseHub.Listen()
	sseHub.ReceiveHubEvents()

	return sseHub
}
