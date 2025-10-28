package auth

import (
	"khelogames/api/transactions"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type AuthServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
	txStore    *transactions.SQLStore
}

func NewAuthServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, txStore *transactions.SQLStore) *AuthServer {
	return &AuthServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, txStore: txStore}
}
