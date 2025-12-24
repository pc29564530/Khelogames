package teams

import (
	"khelogames/api/transactions"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type TeamsServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
	txStore    *transactions.SQLStore
}

func NewTeamsServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, txStore *transactions.SQLStore) *TeamsServer {
	return &TeamsServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, txStore: txStore}
}
