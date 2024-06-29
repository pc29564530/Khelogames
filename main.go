package main

import (
	"database/sql"
	"fmt"
	"khelogames/api/auth"
	"khelogames/api/cricket"
	"khelogames/api/football"
	"khelogames/api/handlers"
	"khelogames/api/server"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	newLogger := logger.NewLogger()

	config, err := util.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	log := logger.NewLogger()

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		fmt.Errorf("cannot create token maker: %w", err)
	}

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

	//handlers.NewCommunityMessageSever(store, log)
	//handlers.NewClubMemberServer(store, log)
	playerProfileServer := handlers.NewPlayerProfileServer(store, log)
	tournamentGroupServer := handlers.NewTournamentGroup(store, log)
	tournamentMatchServer := handlers.NewTournamentMatch(store, log)
	tournamentOrganizerServer := handlers.NewTournamentOrganizerServer(store, log)
	tournamentStanding := handlers.NewTournamentStanding(store, log)
	//handlers.NewClubMemberServer(store, log)
	footballMatchServer := football.NewFootballMatches(store, log)
	cricketMatchServer := cricket.NewCricketMatch(store, log)
	tournamentServer := handlers.NewTournamentServer(store, log)
	cricketMatchTossServer := cricket.NewCricketMatchToss(store, log)
	cricketMatchPlayerScoreServer := cricket.NewCricketMatchScore(store, log)
	ClubTournamentServer := handlers.NewClubTournamentServer(store, log)

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
	)
	if err != nil {
		newLogger.Error("Server does not created", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		newLogger.Error("cannot start server", err)
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
