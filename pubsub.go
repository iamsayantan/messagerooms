package messagerooms

import (
	"bytes"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// ServerEvent is type for any SSE events
type ServerEvent string

var (
	ConnectionEvent ServerEvent = "ClientConnection"
	HeartbeatEvent  ServerEvent = "Heartbeat"
)

// EventsourceConnection represents a single persistent connection.
type EventsourceConnection struct {
	ConnectionID string            // ConnectionID is an unique connection id for the connection
	User         *User             // User for whom the connection is opened
	SendCh       chan EventMessage // SendCh channel is used to send messages to the particular connection

	ticker  *time.Ticker // ticker is used for sending heartbeats to the client
	closing chan bool    // closing channel receives data if the client closes
}

// PublishEvent publishes an event to connection's SendChannel
func (ec *EventsourceConnection) PublishEvent(evt EventMessage) {
	ec.SendCh <- evt
}

// Closing is for housekeeping works.
func (ec *EventsourceConnection) Closing() {
	ec.closing <- true
}

// Heartbeat registers a goroutine that periodically pings over the persistent connection so that the client
// does not close the connection.
func (ec *EventsourceConnection) Heartbeat() {
	go func() {
		for {
			select {
			case <-ec.closing:
				ec.ticker.Stop()
				return
			case <-ec.ticker.C:
				msg := EventMessage{Event: HeartbeatEvent, DestinationID: ec.ConnectionID, Data: "Heartbeat"}
				ec.PublishEvent(msg)
			}
		}
	}()
}

// EventMessage represents a single SSE Event.
type EventMessage struct {
	Event         ServerEvent // Event is the name of the event.
	DestinationID string      `json:"destination_id"` // ConnectionID of the EventsourceConnection where this message should be delivered
	Data          string      `json:"data"`           // Data is what we send in the response
}

// String converts the event to a string eligible for publishing to SSE connection.
func (evt *EventMessage) String() string {
	var buff bytes.Buffer

	if len(evt.DestinationID) > 0 {
		buff.WriteString(fmt.Sprintf("id: %s\n", evt.DestinationID))
	}

	if len(evt.Event) > 0 {
		buff.WriteString(fmt.Sprintf("event: %s\n", evt.Event))
	}

	if len(evt.Data) > 0 {
		buff.WriteString(fmt.Sprintf("data: %s\n", evt.Data))
	}

	buff.WriteString("\n")
	return buff.String()
}

// NewEventsourceConnection returns a new EventsourceConnection
func NewEventsourceConnection(user *User) *EventsourceConnection {
	id, _ := uuid.NewV4()
	ticker := time.NewTicker(time.Second * 5)
	closingChan := make(chan bool)

	eventsourceConnection := &EventsourceConnection{
		ConnectionID: id.String(),
		User:         user,
		SendCh:       make(chan EventMessage),
		ticker:       ticker,
		closing:      closingChan,
	}
	return eventsourceConnection
}
