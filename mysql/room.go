package mysql

import (
	"errors"
	messagerooms "github.com/iamsayantan/MessageRooms"
	"github.com/jinzhu/gorm"
)

var (
	// ErrRoomNotFound is returned when we don't get the room in our data store.
	ErrRoomNotFound = errors.New("room not found")

	// ErrUserAlreadyMember is returned when user tries to join an room he is already part of
	ErrUserAlreadyMember = errors.New("user is already part of the room")
)

type roomRepository struct {
	db *gorm.DB
}

func (r *roomRepository) Create(name string, user messagerooms.User) (*messagerooms.Room, error) {
	panic("implement me")
}

func (r *roomRepository) Find(id string) (*messagerooms.Room, error) {
	room := messagerooms.Room{}
	if notFound := r.db.Preload("CreatedBy").Preload("Users").Where("id = ?", id).First(&room).RecordNotFound(); notFound {
		return nil, ErrRoomNotFound
	}

	return &room, nil
}

func (r *roomRepository) AddUserToRoom(room messagerooms.Room, user messagerooms.User) error {
	if alreadyExistsInRoom := r.CheckUserExistsInRoom(room, user); alreadyExistsInRoom {
		return ErrUserAlreadyMember
	}

	if err := r.db.Model(&room).Association("Users").Append(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) CheckUserExistsInRoom(room messagerooms.Room, user messagerooms.User) bool {
	var existingUser messagerooms.User

	sql := `
		SELECT users.id, users.nickname FROM users INNER JOIN room_users ON room_users.user_id = users.id
		WHERE room_users.user_id = ? AND room_users.room_id = ? LIMIT 1 
	`
	r.db.Raw(sql, user.ID, room.ID).Scan(&existingUser)

	return existingUser.ID != ""
}

// NewRoomRepository returns implementation of RoomRepository interface.
func NewRoomRepository(db *gorm.DB) messagerooms.RoomRepository {
	return &roomRepository{db: db}
}
