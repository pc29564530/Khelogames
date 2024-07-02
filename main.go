package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"khelogames/api/auth"
	"khelogames/api/cricket"
	"khelogames/api/football"
	"khelogames/api/handlers"
	"khelogames/api/messenger"
	"khelogames/api/server"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

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
		log.Errorf("cannot create token maker: %v", err)
		os.Exit(1)
	}

	rabbitConn, rabbitChan, err := messenger.StartRabbitMQ(config)
	if err != nil {
		log.Fatal("cannot start RabbitMQ:", err)
	}
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
	broadcast := make(chan []byte)

	// Initialize HTTP servers and handlers
	otpServer := auth.NewOtpServer(store, log)
	loginServer := auth.NewLoginServer(store, log, tokenMaker, config)
	signupServer := auth.NewSignupServer(store, log)

	sessionServer := auth.NewSessionServer(store, log)
	threadServer := handlers.NewThreadServer(store, log)
	tokenServer := auth.NewTokenServer(store, log, tokenMaker)
	profileServer := handlers.NewProfileServer(store, log)
	likeThread := handlers.NewLikeThreadServer(store, log)
	clubServer := handlers.NewClubServer(store, log)
	userServer := handlers.NewUserServer(store, log, tokenMaker, config)
	followServer := handlers.NewFollowServer(store, log)
	communityServer := handlers.NewCommunityServer(store, log)
	joinCommunityServer := handlers.NewJoinCommunityServer(store, log)
	commentServer := handlers.NewCommentServer(store, log)
	clubMemberServer := handlers.NewClubMemberServer(store, log)
	groupTeamServer := handlers.NewGroupTeamServer(store, log)
	playerProfileServer := handlers.NewPlayerProfileServer(store, log)
	tournamentGroupServer := handlers.NewTournamentGroup(store, log)
	tournamentMatchServer := handlers.NewTournamentMatch(store, log)
	tournamentOrganizerServer := handlers.NewTournamentOrganizerServer(store, log)
	tournamentStanding := handlers.NewTournamentStanding(store, log)
	footballMatchServer := football.NewFootballMatches(store, log)
	cricketMatchServer := cricket.NewCricketMatch(store, log)
	tournamentServer := handlers.NewTournamentServer(store, log)
	cricketMatchTossServer := cricket.NewCricketMatchToss(store, log)
	cricketMatchPlayerScoreServer := cricket.NewCricketMatchScore(store, log)
	ClubTournamentServer := handlers.NewClubTournamentServer(store, log)
	footballUpdateServer := football.NewFootballUpdate(store, log)

	// Initialize WebSocket handler
	webSocketHandlerImpl := messenger.NewWebSocketHandler(store, tokenMaker, clients, broadcast, upgrader, rabbitChan, log)
	messageServer := messenger.NewMessageServer(store, log, broadcast)

	// Initialize Gin router
	router := gin.Default()
	server, err := server.NewServer(config,
		store,
		tokenMaker,
		log,
		otpServer,
		signupServer,
		loginServer,
		tokenServer,
		sessionServer,
		threadServer,
		profileServer,
		likeThread,
		clubServer,
		userServer,
		followServer,
		communityServer,
		joinCommunityServer,
		commentServer,
		clubMemberServer,
		groupTeamServer,
		playerProfileServer,
		tournamentGroupServer,
		tournamentOrganizerServer,
		tournamentMatchServer,
		tournamentStanding,
		footballMatchServer,
		cricketMatchServer,
		tournamentServer,
		cricketMatchTossServer,
		cricketMatchPlayerScoreServer,
		ClubTournamentServer,
		footballUpdateServer,
		webSocketHandlerImpl,
		messageServer,
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
