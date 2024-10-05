package players

import (
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
)

type PlayerServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

func NewPlayerServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *PlayerServer {
	return &PlayerServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
