package sports

import (
	"khelogames/api/sports/badminton"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type SportsServer struct {
	store           *db.Store
	logger          *logger.Logger
	tokenMaker      token.Maker
	config          util.Config
	footballServer  *football.FootballServer
	cricketServer   *cricket.CricketServer
	badmintonServer *badminton.BadmintonServer
}

func NewSportsServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, footballServer *football.FootballServer, cricketServer *cricket.CricketServer, badmintonServer *badminton.BadmintonServer) *SportsServer {
	return &SportsServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, footballServer: footballServer, cricketServer: cricketServer, badmintonServer: badmintonServer}
}
