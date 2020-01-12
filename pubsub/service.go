package pubsub

import (
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/iamsayantan/messagerooms"
)

// Service interface defines methods for interacting with the pubsub system.
type Service interface {
	Publish(data messagerooms.Publishable)

	// Subscribe adds the given connection id to the topic's subscribed connection list. During publish
	// we fetch all the connectionIDs and dispatch it to the client.
	Subscribe(topic, connectionID string)
	Unsubscribe(topic, connectionID string)
}

type redisPubsubService struct {
	redisConn redis.Conn
}

func (rs *redisPubsubService) Subscribe(topic, connectionID string) {
	_, err := rs.redisConn.Do("LPUSH", topic, connectionID)
	if err != nil {
		log.Printf("Error: %s, subscribing to topic: %s, connectionID: %s", err.Error(), topic, connectionID)
	}
}

func (rs *redisPubsubService) Unsubscribe(topic, connectionID string) {
	// connection id is stored as a list in redis against the topic name. so to unsubscribe we just need to
	// do LREM the connection id from the list.
	_, err := rs.redisConn.Do("LREM", topic, 1, connectionID)
	if err != nil {
		log.Printf("Error: %s, unsubscribing to topic: %s, connectionID: %s", err.Error(), topic, connectionID)
	}
}

func (rs *redisPubsubService) Publish(data messagerooms.Publishable) {
	// for publishing data we find all the connection id that is subscribed to the given topic and prepare
	// event for all of those connection ids and publish
	topic := data.GetTopic()
	connIDs, err := redis.Strings(rs.redisConn.Do("LRANGE", topic, 0, -1))
	if err != nil {
		log.Printf("Error: %s, fetching topic: %s, ", err.Error(), topic)
	}

	for _, connID := range connIDs {
		publishEvent := data.ToPublish(connID)
		jsonEvent, err := publishEvent.ToJSON()
		if err != nil {
			return // silently ignoring the error
		}
		_, _ = rs.redisConn.Do("PUBLISH", messagerooms.HubChannel, jsonEvent)
	}
}

// NewRedisPubsubService returns an new instance of redis pubsub service
func NewRedisPubsubService(conn redis.Conn) Service {
	return &redisPubsubService{redisConn: conn}
}
