package main

import (
	"context"
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
	"os/signal"
	"syscall"
	"time"

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
	// Load app.env if it exists (local dev). In production, env vars are set directly.
	if err := godotenv.Load("./app.env"); err != nil {
		fmt.Println("No app.env file found, using environment variables")
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
		log.Warnf("RabbitMQ not available, chat features will be disabled: %v", err)
	}
	if rabbitChan != nil {
		_, _ = rabbitChan.QueueDeclare("chatHub", true, false, false, false, nil)
		_, _ = rabbitChan.QueueDeclare("scoreHub", true, false, false, false, nil)
	}
	if rabbitConn != nil {
		defer rabbitConn.Close()
	}

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

	teamsServer := teams.NewTeamsServer(store, log, tokenMaker, config, txStore)
	tournamentServer := tournaments.NewTournamentServer(store, log, tokenMaker, config, nil, txStore)
	tokenServer := apiToken.NewTokenServer(store, log, tokenMaker, config)

	// Create messenger server with cricket server as both updater and broadcaster
	messengerServer := messenger.NewMessageServer(store, tokenMaker, clients, messageBroadCast, scoredBroadCast, upgrader, rabbitChan, log, nil)
	playerServer := players.NewPlayerServer(store, log, tokenMaker, config)
	sportsServer := sports.NewSportsServer(store, log, tokenMaker, config)
	tournamentServer.SetScoreBroadcaster(hub)
	cricketServer.SetScoreBroadcaster(hub)
	footballServer.SetScoreBroadcaster(hub)
	txStore.SetScoreBroadcaster(hub)

	log.Info("Broadcasters initialized for cricket, football, tournament, and messenger")

	if rabbitChan != nil {
		go hub.StartRabbitMQConsumer("scoreHub")
		go hub.StartRabbitMQConsumer("chatHub")
	}
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
		conn,
		rabbitChan,
		nil, // httpServer will be created in Start()
	)
	if err != nil {
		newLogger.Error("Server creation failed", err)
		os.Exit(1)
	}

	// Start server
	go func() {
		err = server.Start(config.ServerAddress)
		if err != nil && err != http.ErrServerClosed {
			newLogger.Error("Server start failed", err)
			os.Exit(1)
		}
	}()

	//Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", err)
	}

	// Close database connection
	if conn != nil {
		conn.Close()
		log.Info("Database connection closed")
	}

	// Close RabbitMQ connection
	if rabbitChan != nil {
		rabbitChan.Close()
		log.Info("RabbitMQ connection closed")
	}

	log.Info("Server exited gracefully")
}
