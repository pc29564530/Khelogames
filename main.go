package main

import (
	"database/sql"
	"fmt"
	"khelogames/api/auth"
	"khelogames/api/handlers"
	"khelogames/api/players"
	"khelogames/api/sports"
	"khelogames/api/teams"
	apiToken "khelogames/api/token"
	"khelogames/api/transactions"
	coreToken "khelogames/core/token"
	"khelogames/hub"
	"net/http"
	"os"

	"khelogames/api/messenger"
	"khelogames/api/server"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	"khelogames/api/tournaments"
	db "khelogames/database"
	"khelogames/logger"
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
	log := logger.NewLogger()

	config, _ := util.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)

	tokenMaker, err := coreToken.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Errorf("cannot create token maker: ", err)
		os.Exit(1)
	}

	txStore := transactions.NewStore(conn, &tokenMaker, log, nil)

	rabbitConn, rabbitChan, err := hub.StartRabbitMQ(config)
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
	subscriber := make(map[string]map[*hub.Client]bool)
	hub := hub.NewHub(store, log, upgrader, rabbitChan, tokenMaker, nil, nil, subscriber)

	// Channel for broadcasting messages to WebSocket clients
	messageBroadCast := make(chan []byte)
	scoredBroadCast := make(chan []byte)

	cricketServer := cricket.NewCricketServer(store, log, nil, txStore)

	// Initialize HTTP servers and handlers
	authServer := auth.NewAuthServer(store, log, tokenMaker, config, txStore)
	handlerServer := handlers.NewHandlerServer(store, log, tokenMaker, config, txStore)
	footballServer := football.NewFootballServer(store, log, nil, txStore)

	teamsServer := teams.NewTeamsServer(store, log, tokenMaker, config)
	tournamentServer := tournaments.NewTournamentServer(store, log, tokenMaker, config, nil, txStore)
	tokenServer := apiToken.NewTokenServer(store, log, tokenMaker, config)

	// Create messenger server with cricket server as both updater and broadcaster
	messengerServer := messenger.NewMessageServer(store, tokenMaker, clients, messageBroadCast, scoredBroadCast, upgrader, rabbitChan, log, nil)
	playerServer := players.NewPlayerServer(store, log, tokenMaker, config)
	sportsServer := sports.NewSportsServer(store, log, tokenMaker, config)
	tournamentServer.SetScoreBroadcaster(hub)
	cricketServer.SetScoreBroadcaster(hub)
	footballServer.SetScoreBroadcaster(hub)

	fmt.Printf("Cricket broadcaster pointer: %p\n", cricketServer.GetScoreBroadcaster())
	fmt.Printf("Football broadcaster pointer: %p\n", footballServer.GetScoreBroadcaster())
	fmt.Printf("Messenger server pointer:    %p\n", messengerServer)
	fmt.Printf("Tournament broadcaster pointer: %p\n", tournamentServer.GetScoreBroadcaster())

	go hub.StartRabbitMQConsumer("scoreHub")
	go hub.StartRabbitMQConsumer("chatHub")
	// go hub.StartMessageHub()

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
		tokenServer,
		hub,
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
