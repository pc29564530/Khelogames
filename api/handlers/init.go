package handlers

import (
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
)

type HandlersServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

func NewHandlerServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *HandlersServer {
	return &HandlersServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
