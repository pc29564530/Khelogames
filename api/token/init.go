package token

import (
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type TokenServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

func NewTokenServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *TokenServer {
	return &TokenServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
