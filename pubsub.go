package messagerooms

// EventsourceConnection represents a single persistent connection.
type EventsourceConnection struct {
	ConnectionID string            // ConnectionID is an unique connection id for the connection
	User         *User             // User for whom the connection is opened
	SendCh       chan EventMessage // SendCh channel is used to send messages to the particular connection
}

// PublishEvent publishes an event to connection's SendChannel
func (ec *EventsourceConnection) PublishEvent(evt EventMessage) {
	ec.SendCh <- evt
}

type EventMessage struct {
}
