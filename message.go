package messagerooms

import "time"

// Message struct represents a single message
type Message struct {
	ID          string
	UserID      string
	MessageText string
	CreatedAt   time.Time
	CreatedBy   User
}

// MessageRepository provides interface to access message storage.
type MessageRepository interface {
	PostMessage(room Room, user User, messageText string) (*Message, error)
}
