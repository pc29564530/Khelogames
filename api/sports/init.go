package sports

import (
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type SportsServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

func NewSportsServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *SportsServer {
	return &SportsServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
