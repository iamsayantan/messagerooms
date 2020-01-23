package room

import (
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/iamsayantan/messagerooms"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func (s *instrumentingService) CreateNewRoom(roomName string, user messagerooms.User) (*messagerooms.Room, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create_room").Add(1)
		s.requestLatency.With("method", "create_room").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.CreateNewRoom(roomName, user)
}

func (s *instrumentingService) RoomDetails(id string) (*messagerooms.Room, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "room_details").Add(1)
		s.requestLatency.With("method", "room_details").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.RoomDetails(id)
}

func (s *instrumentingService) AllRooms() ([]*messagerooms.Room, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "all_rooms").Add(1)
		s.requestLatency.With("method", "all_rooms").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.AllRooms()
}

func (s *instrumentingService) AddUserToRoom(room messagerooms.Room, user messagerooms.User) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "add_user_to_room").Add(1)
		s.requestLatency.With("method", "add_user_to_room").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.AddUserToRoom(room, user)
}

func (s *instrumentingService) CheckUserExistsInRoom(room messagerooms.Room, user messagerooms.User) bool {
	defer func(begin time.Time) {
		s.requestCount.With("method", "check_user_exists_in_room").Add(1)
		s.requestLatency.With("method", "check_user_exists_in_room").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.CheckUserExistsInRoom(room, user)
}

func (s *instrumentingService) GetAllRoomMessages(room messagerooms.Room) ([]*messagerooms.Message, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get_all_room_messages").Add(1)
		s.requestLatency.With("method", "get_all_room_messages").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetAllRoomMessages(room)
}

func (s *instrumentingService) PostMessage(room messagerooms.Room, user messagerooms.User, messageText string) (*messagerooms.Message, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "post_message").Add(1)
		s.requestLatency.With("method", "post_message").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.PostMessage(room, user, messageText)
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, summary metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: summary,
		next:           s,
	}
}
