package mysql

import (
	"github.com/iamsayantan/MessageRooms"
	"github.com/jinzhu/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func (r *roomRepository) Create(name string, user messagerooms.User) (*messagerooms.Room, error) {
	panic("implement me")
}

func (r *roomRepository) Find(id string) (*messagerooms.Room, error) {
	panic("implement me")
}

func (r *roomRepository) AddUserToRoom(room messagerooms.Room, user messagerooms.User) error {
	panic("implement me")
}

func (r *roomRepository) CheckUserExistsInRoom(room messagerooms.Room, user messagerooms.User) bool {
	panic("implement me")
}

// NewRoomRepository returns implementation of RoomRepository interface.
func NewRoomRepository(db *gorm.DB) messagerooms.RoomRepository {
	return &roomRepository{db: db}
}
