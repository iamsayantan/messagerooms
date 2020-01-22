package room

import (
	"github.com/go-kit/kit/metrics"
	"github.com/iamsayantan/messagerooms"
)

type instrumentingService struct {
	requestCount metrics.Counter
	next         Service
}

func (s *instrumentingService) CreateNewRoom(roomName string, user messagerooms.User) (*messagerooms.Room, error) {
	defer func() {
		s.requestCount.With("method", "create_room").Add(1)
	}()

	return s.next.CreateNewRoom(roomName, user)
}

func (s *instrumentingService) RoomDetails(id string) (*messagerooms.Room, error) {
	defer func() {
		s.requestCount.With("method", "room_details").Add(1)
	}()

	return s.next.RoomDetails(id)
}

func (s *instrumentingService) AllRooms() ([]*messagerooms.Room, error) {
	defer func() {
		s.requestCount.With("method", "all_rooms").Add(1)
	}()

	return s.next.AllRooms()
}

func (s *instrumentingService) AddUserToRoom(room messagerooms.Room, user messagerooms.User) error {
	defer func() {
		s.requestCount.With("method", "add_user_to_room").Add(1)
	}()

	return s.next.AddUserToRoom(room, user)
}

func (s *instrumentingService) CheckUserExistsInRoom(room messagerooms.Room, user messagerooms.User) bool {
	defer func() {
		s.requestCount.With("method", "check_user_exists_in_room").Add(1)
	}()

	return s.next.CheckUserExistsInRoom(room, user)
}

func (s *instrumentingService) GetAllRoomMessages(room messagerooms.Room) ([]*messagerooms.Message, error) {
	defer func() {
		s.requestCount.With("method", "get_all_room_messages").Add(1)
	}()

	return s.next.GetAllRoomMessages(room)
}

func (s *instrumentingService) PostMessage(room messagerooms.Room, user messagerooms.User, messageText string) (*messagerooms.Message, error) {
	defer func() {
		s.requestCount.With("method", "post_message").Add(1)
	}()

	return s.next.PostMessage(room, user, messageText)
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, s Service) Service {
	return &instrumentingService{
		requestCount: counter,
		next:         s,
	}
}
