package mysql

import (
	"errors"

	messagerooms "github.com/iamsayantan/MessageRooms"
	"github.com/jinzhu/gorm"
)

var (
	// ErrRoomNotFound is returned when we don't get the room in our data store.
	ErrRoomNotFound = errors.New("room not found")
)

type roomRepository struct {
	db *gorm.DB
}

func (r *roomRepository) Create(name string, user messagerooms.User) (*messagerooms.Room, error) {
	panic("implement me")
}

func (r *roomRepository) Find(id string) (*messagerooms.Room, error) {
	room := messagerooms.Room{}
	if notFound := r.db.Preload("CreatedBy").Where("id = ?", id).First(&room).RecordNotFound(); notFound {
		return nil, ErrRoomNotFound
	}

	return &room, nil
}

func (r *roomRepository) AddUserToRoom(room messagerooms.Room, user messagerooms.User) error {
	r.db.Model(&room).Association("Users").Find(&user)
	return nil
}

func (r *roomRepository) CheckUserExistsInRoom(room messagerooms.Room, user messagerooms.User) bool {
	r.db.Model(&room).Association("Users").Find(&user)
	return true
}

// NewRoomRepository returns implementation of RoomRepository interface.
func NewRoomRepository(db *gorm.DB) messagerooms.RoomRepository {
	return &roomRepository{db: db}
}
