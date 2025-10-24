package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"khelogames/api/auth"
	"khelogames/api/handlers"
	"khelogames/api/players"
	"khelogames/api/sports"
	"khelogames/api/teams"

	"khelogames/api/messenger"
	"khelogames/api/server"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	"khelogames/api/tournaments"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load("./app.env")
	if err != nil {
		fmt.Errorf("Unable to read env file: ", err)
	}
}

func main() {
	newLogger := logger.NewLogger()

	config, _ := util.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	log := logger.NewLogger()

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Errorf("cannot create token maker: ", err)
		os.Exit(1)
	}

	rabbitConn, rabbitChan, err := messenger.StartRabbitMQ(config)
	if err != nil {
		log.Fatal("cannot start RabbitMQ:", err)
	}
	_, _ = rabbitChan.QueueDeclare("chatHub", true, false, false, false, nil)
	_, _ = rabbitChan.QueueDeclare("scoreHub", true, false, false, false, nil)
	defer rabbitConn.Close()

	// Define clients map for WebSocket connections
	clients := make(map[*websocket.Conn]bool)

	// WebSocket upgrader configuration
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Channel for broadcasting messages to WebSocket clients
	messageBroadCast := make(chan []byte)
	scoredBroadCast := make(chan []byte)

	cricketServer := cricket.NewCricketServer(store, log, rabbitChan, nil)

	// Initialize HTTP servers and handlers
	authServer := auth.NewAuthServer(store, log, tokenMaker, config)
	handlerServer := handlers.NewHandlerServer(store, log, tokenMaker, config)
	footballServer := football.NewFootballServer(store, log, nil)

	teamsServer := teams.NewTeamsServer(store, log, tokenMaker, config)
	tournamentServer := tournaments.NewTournamentServer(store, log, tokenMaker, config, nil)

	// Create messenger server with cricket server as both updater and broadcaster
	messengerServer := messenger.NewMessageServer(store, tokenMaker, clients, messageBroadCast, scoredBroadCast, upgrader, rabbitChan, log, cricketServer, nil)
	playerServer := players.NewPlayerServer(store, log, tokenMaker, config)
	sportsServer := sports.NewSportsServer(store, log, tokenMaker, config)
	cricketServer.SetScoreBroadcaster(messengerServer)
	footballServer.SetScoreBroadcaster(messengerServer)
	tournamentServer.SetScoreBroadcaster(messengerServer)

	go messengerServer.StartRabbitMQConsumer("scoreHub")
	go messengerServer.StartRabbitMQConsumer("chatHub")
	go messengerServer.StartMessageHub()
	go messengerServer.StartCricketScoreHub()

	// Initialize Gin router
	router := gin.Default()
	server, err := server.NewServer(config,
		store,
		tokenMaker,
		log,
		authServer,
		handlerServer,
		tournamentServer,
		footballServer,
		cricketServer,
		teamsServer,
		messengerServer,
		playerServer,
		sportsServer,
		router,
	)
	if err != nil {
		newLogger.Error("Server creation failed", err)
		os.Exit(1)
	}

	// Start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		newLogger.Error("Server start failed", err)
		os.Exit(1)
	}
}
