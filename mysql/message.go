package mysql

import (
	"github.com/iamsayantan/MessageRooms"
	"github.com/jinzhu/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func (m *messageRepository) PostMessage(room messagerooms.Room, user messagerooms.User, messageText string) (*messagerooms.Message, error) {
	panic("implement me")
}

// NewMessageRepository returns implementation of MessageRepository interface.
func NewMessageRepository(db *gorm.DB) messagerooms.MessageRepository {
	return &messageRepository{db: db}
}
