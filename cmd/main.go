package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gomodule/redigo/redis"
	"github.com/iamsayantan/messagerooms"
	"github.com/iamsayantan/messagerooms/mysql"
	"github.com/iamsayantan/messagerooms/pubsub"
	"github.com/iamsayantan/messagerooms/room"
	"github.com/iamsayantan/messagerooms/server"
	"github.com/iamsayantan/messagerooms/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	defaultDBHost     = getFromEnv("MYSQL_HOST", "localhost")
	defaultDBPort     = getFromEnv("MYSQL_PORT", "3306")
	defaultDBUsername = getFromEnv("MYSQL_USERNAME", "root")
	defaultDBPassword = getFromEnv("MYSQL_PASSWORD", "12345")
	defaultDBName     = getFromEnv("DATABASE_NAME", "rooms")

	defaultServerPort = "9050"
)

func getFromEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func main() {
	dbHost := flag.String("db.host", defaultDBHost, "Database host url")
	dbPort := flag.String("db.port", defaultDBPort, "Database port")
	dbUsername := flag.String("db.username", defaultDBUsername, "Database username")
	dbPassword := flag.String("db.password", defaultDBPassword, "Database password")
	serverPort := flag.String("server.port", defaultServerPort, "Server port where the server runs")

	mysqlHost := os.Getenv("MYSQL_HOST")
	log.Printf("ENV Database host: %s", mysqlHost)

	flag.Parse()

	// connect to the database
	// format: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	dbCred := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", *dbUsername, *dbPassword, *dbHost, *dbPort, defaultDBName)
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
		userService   user.Service
		roomService   room.Service
		pubsubService pubsub.Service
	)

	// can not use same underlying connection instance for publishing and subscribing.
	rSubConn, err := redisConn()
	if err != nil {
		panic(err)
	}
	defer rSubConn.Close()

	rPubConn, err := redisConn()
	if err != nil {
		panic(err)
	}
	defer rPubConn.Close()

	userRepo = mysql.NewUserRepository(db)
	roomRepo = mysql.NewRoomRepository(db)
	messageRepo = mysql.NewMessageRepository(db)

	labelNames := []string{"method"}

	pubsubService = pubsub.NewRedisPubsubService(rPubConn)
	userService = user.NewService(userRepo)
	roomService = room.NewService(roomRepo, messageRepo, pubsubService)
	roomService = room.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "messagerooms_api",
			Subsystem: "room_service",
			Name:      "request_count",
			Help:      "Number of requests received",
		}, labelNames),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "messagerooms_api",
			Subsystem: "room_service",
			Name:      "request_latency",
		}, labelNames),
		roomService,
	)

	pubsubConn = &redis.PubSubConn{Conn: rSubConn}

	hub := server.NewSSEHub(pubsubConn, pubsubService)
	srv := server.NewServer(userService, roomService, hub)

	log.Printf("Server starting on port %s", *dbPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), srv))
}
