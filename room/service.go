package room

import (
	"github.com/iamsayantan/messagerooms"
	"github.com/iamsayantan/messagerooms/pubsub"
	"github.com/pkg/errors"
)

var (
	// ErrUserAlreadyInRoom is returned when user tries to join same room twice.
	ErrUserAlreadyInRoom = errors.New("user is already in room")

	// ErrUserNotInRoom is returned when user tries to do something that reburies him to be a member of the room
	ErrUserNotInRoom = errors.New("user is not a member of the room")
)

// Service provides methods for room management.
type Service interface {
	// CreateNewRoom creates a new room.
	CreateNewRoom(roomName string, user messagerooms.User) (*messagerooms.Room, error)

	// RoomDetails returns details for the room with the id. Error is returned in case of invalid id.
	RoomDetails(id string) (*messagerooms.Room, error)

	// AllRooms returns list of all rooms available
	AllRooms() ([]*messagerooms.Room, error)

	// AddUserToRoom adds an user to a room.
	AddUserToRoom(room messagerooms.Room, user messagerooms.User) error

	// CheckUserExistsInRoom checks if an user is member of a room.
	CheckUserExistsInRoom(room messagerooms.Room, user messagerooms.User) bool

	// GetAllRoomMessages returns all the messages posted in a room.
	GetAllRoomMessages(room messagerooms.Room) ([]*messagerooms.Message, error)

	// PostMessage posts a message in a room.
	PostMessage(room messagerooms.Room, user messagerooms.User, messageText string) (*messagerooms.Message, error)
}

type roomService struct {
	room      messagerooms.RoomRepository
	message   messagerooms.MessageRepository
	publisher pubsub.Service
}

func (s *roomService) AllRooms() ([]*messagerooms.Room, error) {
	return s.room.FindAll()
}

func (s *roomService) GetAllRoomMessages(room messagerooms.Room) ([]*messagerooms.Message, error) {
	return s.message.GetMessagesByRoom(room)
}

func (s *roomService) CreateNewRoom(roomName string, user messagerooms.User) (*messagerooms.Room, error) {
	room, err := s.room.Create(roomName, user)
	if err != nil {
		return nil, err
	}

	return s.RoomDetails(room.ID)
}

func (s *roomService) RoomDetails(id string) (*messagerooms.Room, error) {
	room, err := s.room.Find(id)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) AddUserToRoom(room messagerooms.Room, user messagerooms.User) error {
	if exists := s.room.CheckUserExistsInRoom(room, user); exists {
		return ErrUserAlreadyInRoom
	}

	err := s.room.AddUserToRoom(room, user)
	return err
}

func (s *roomService) CheckUserExistsInRoom(room messagerooms.Room, user messagerooms.User) bool {
	return s.room.CheckUserExistsInRoom(room, user)
}

func (s *roomService) PostMessage(room messagerooms.Room, user messagerooms.User, messageText string) (*messagerooms.Message, error) {
	if exists := s.room.CheckUserExistsInRoom(room, user); !exists {
		return nil, ErrUserNotInRoom
	}

	message, err := s.message.PostMessage(room, user, messageText)
	if err != nil {
		return nil, err
	}

	// publishing the new message into the pubsub system.
	go func() {
		users, _ := s.room.GetRoomMembers(room)
		for _, user := range users {
			messageEvent := messagerooms.NewMessageEvent{User: *user, Room: room, Message: *message}
			s.publisher.Publish(&messageEvent)
		}
	}()

	return message, nil
}

// NewService returns a new room service with associated dependency.
func NewService(rs messagerooms.RoomRepository, ms messagerooms.MessageRepository, pub pubsub.Service) Service {
	service := &roomService{
		room:      rs,
		message:   ms,
		publisher: pub,
	}

	return service
}
