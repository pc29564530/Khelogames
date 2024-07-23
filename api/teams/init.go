package teams

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
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
