package teams

import (
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
}

func NewTeamsServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *TeamsServer {
	return &TeamsServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
