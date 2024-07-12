package clubs

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
)

type ClubServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

func NewClubsServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *ClubServer {
	return &ClubServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
