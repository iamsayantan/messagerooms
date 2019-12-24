package mysql

import (
	"errors"
	"time"

	messagerooms "github.com/iamsayantan/MessageRooms"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrMessageNotFound is returned when no message is found with the given message id.
	ErrMessageNotFound = errors.New("message not found")
)

type messageRepository struct {
	db *gorm.DB
}

func (m *messageRepository) PostMessage(room messagerooms.Room, user messagerooms.User, messageText string) (*messagerooms.Message, error) {
	id, _ := uuid.NewV4()
	msg := messagerooms.Message{
		ID:          id.String(),
		MessageText: messageText,
		RoomID:      room.ID,
		UserID:      user.ID,
		CreatedAt:   time.Now(),
	}

	if err := m.db.Preload("CreatedBy").Create(&msg).Error; err != nil {
		return nil, err
	}

	// refetching the message for loading the relations.
	return m.GetMessage(msg.ID)
}

func (m *messageRepository) GetMessage(messageID string) (*messagerooms.Message, error) {
	var msg messagerooms.Message
	if notFound := m.db.Preload("CreatedBy").Where("id = ?", messageID).First(&msg).RecordNotFound(); notFound {
		return nil, ErrMessageNotFound
	}

	return &msg, nil
}

// NewMessageRepository returns implementation of MessageRepository interface.
func NewMessageRepository(db *gorm.DB) messagerooms.MessageRepository {
	return &messageRepository{db: db}
}
