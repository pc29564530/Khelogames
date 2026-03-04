package sports

import (
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type SportsServer struct {
	store          *db.Store
	logger         *logger.Logger
	tokenMaker     token.Maker
	config         util.Config
	footballServer *football.FootballServer
	cricketServer  *cricket.CricketServer
}

func NewSportsServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, footballServer *football.FootballServer, cricketServer *cricket.CricketServer) *SportsServer {
	return &SportsServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, footballServer: footballServer, cricketServer: cricketServer}
}
