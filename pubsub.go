package messagerooms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// ServerEvent is type for any SSE events
type ServerEvent string

// HubChannel is the main channel where any hub related events should be published. It is upto the hub then to
// decide what to do with the events.
const HubChannel = "HubChannel"

var (
	ConnectionEvent ServerEvent = "ClientConnection"
	HeartbeatEvent  ServerEvent = "Heartbeat"
	MessageEvent    ServerEvent = "MessageEvent"
)

// Publishable is the interface that all types must implement that wish to be published into the pubsub system.
type Publishable interface {
	// GetTopic returns an topic identifier
	GetTopic() string

	// ToPublish() method returns an PublishEvent that can be then pushed into the redis pubsub to be broadcast.
	ToPublish() *PublishEvent
}

// EventsourceConnection represents a single persistent connection.
type EventsourceConnection struct {
	ConnectionID string            // ConnectionID is an unique connection id for the connection
	User         *User             // User for whom the connection is opened
	SendCh       chan EventMessage // SendCh channel is used to send messages to the particular connection

	ticker  *time.Ticker // ticker is used for sending heartbeats to the client
	closing chan bool    // closing channel receives data if the client closes
}

// PublishEvent is the container for publishing events.
type PublishEvent struct {
	Topic     string      `json:"topic"`      // Topic of the event.
	CreatedAt int64       `json:"created_at"` // CreatedAt when the event was created. Can be used to track how much time it takes form creation to delivery.
	Payload   interface{} `json:"payload"`    // Payload actual payload for the event.
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
				data := struct {
					Heartbeat  string `json:"heartbeat"`
					ServerTime int64  `json:"server_time"`
				}{Heartbeat: "OK", ServerTime: time.Now().Unix()}

				msg := EventMessage{Event: HeartbeatEvent, DestinationID: ec.ConnectionID, Data: data}
				ec.PublishEvent(msg)
			}
		}
	}()
}

// ToJSON converts the PublishEvent to a JSON.
func (pe *PublishEvent) ToJSON() (string, error) {
	byts, err := json.Marshal(pe)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

// EventMessage represents a single SSE Event.
type EventMessage struct {
	Event         ServerEvent // Event is the name of the event.
	DestinationID string      `json:"destination_id"` // ConnectionID of the EventsourceConnection where this message should be delivered
	Data          interface{} `json:"data"`           // Data is what we send in the response
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

	byts, err := json.Marshal(evt.Data)
	if err != nil {
		byts = []byte("Not a valid JSON")
	}

	data := string(byts)

	if len(data) > 0 {
		buff.WriteString(fmt.Sprintf("data: %s\n", data))
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

// NewPublishEvent returns a new PublishEvent
func NewPublishEvent(topic string, payload interface{}) *PublishEvent {
	return &PublishEvent{
		Topic:     topic,
		CreatedAt: time.Now().UnixNano(),
		Payload:   payload,
	}
}
