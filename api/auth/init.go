package auth

import (
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
)

type AuthServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

func NewAuthServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *AuthServer {
	return &AuthServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
