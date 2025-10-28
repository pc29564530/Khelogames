package handlers

import (
	"khelogames/api/transactions"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type HandlersServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
	txStore    *transactions.SQLStore
}

func NewHandlerServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, txStore *transactions.SQLStore) *HandlersServer {
	return &HandlersServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, txStore: txStore}
}
