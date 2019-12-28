package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	messagerooms "github.com/iamsayantan/MessageRooms"
	"github.com/iamsayantan/MessageRooms/mysql"
	"github.com/iamsayantan/MessageRooms/room"
	"github.com/iamsayantan/MessageRooms/server"
	"github.com/iamsayantan/MessageRooms/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	defaultDBHost     = "127.0.0.1"
	defaultDBPort     = "3306"
	defaultDBUsername = "root"
	defaultDBPassword = ""
	defaultDBName     = "rooms"

	defaultServerPort = "9050"
)

func main() {
	// connect to the database
	// format: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	dbCred := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", defaultDBUsername, defaultDBPassword, defaultDBHost, defaultDBPort, defaultDBName)
	log.Printf("Database Credential: %s", dbCred)

	db, err := gorm.Open("mysql", dbCred)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Automatically migrate the schemas.
	db.AutoMigrate(&messagerooms.User{}, &messagerooms.Room{}, &messagerooms.Message{})

	// initialize application dependencies
	var (
		// Pubsub related initialization
		pubsubConn *redis.PubSubConn
		redisConn  = func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		}

		// Repositories
		userRepo    messagerooms.UserRepository
		roomRepo    messagerooms.RoomRepository
		messageRepo messagerooms.MessageRepository

		// Services
		userService user.Service
		roomService room.Service
	)

	userRepo = mysql.NewUserRepository(db)
	roomRepo = mysql.NewRoomRepository(db)
	messageRepo = mysql.NewMessageRepository(db)

	userService = user.NewService(userRepo)
	roomService = room.NewService(roomRepo, messageRepo)

	rConn, err := redisConn()
	if err != nil {
		panic(err)
	}

	pubsubConn = &redis.PubSubConn{Conn: rConn}

	hub := server.NewSSEHub(pubsubConn)
	srv := server.NewServer(userService, roomService, hub)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", defaultServerPort), srv))
}
