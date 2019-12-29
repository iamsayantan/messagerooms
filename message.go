package messagerooms

import "time"

// Message struct represents a single message
type Message struct {
	ID          string    `json:"id"`
	UserID      string    `json:"-"`
	RoomID      string    `json:"-"`
	MessageText string    `json:"message_text"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   *User     `json:"created_by,omitempty" gorm:"foreignkey:UserID"`
	RoomDetails *Room     `json:"room_details,omitempty" gorm:"foreignkey:RoomID"`
}

func (m Message) GetTopic() string {
	return TopicMessageRoom + ":" + m.RoomID
}

func (m Message) ToPublish() *PublishEvent {
	return NewPublishEvent(m.GetTopic(), m)
}

// MessageRepository provides interface to access message storage.
type MessageRepository interface {
	PostMessage(room Room, user User, messageText string) (*Message, error)
	GetMessage(messageID string) (*Message, error)
	GetMessagesByRoom(room Room) ([]*Message, error)
}
