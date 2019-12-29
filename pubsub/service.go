package pubsub

import (
	"github.com/gomodule/redigo/redis"
	messagerooms "github.com/iamsayantan/MessageRooms"
)

// Service interface defines methods for interacting with the pubsub system.
type Service interface {
	Publish(data messagerooms.Publishable)
}

type redisPubsubService struct {
	redisConn redis.Conn
}

func (rs *redisPubsubService) Publish(data messagerooms.Publishable) {
	publishEvent := data.ToPublish()
	jsonEvent, err := publishEvent.ToJSON()
	if err != nil {
		return // silently ignoring the error
	}
	_, _ = rs.redisConn.Do("PUBLISH", messagerooms.HubChannel, jsonEvent)
}

// NewRedisPubsubService returns an new instance of redis pubsub service
func NewRedisPubsubService(conn redis.Conn) Service {
	return &redisPubsubService{redisConn: conn}
}
