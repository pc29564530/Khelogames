package tournaments

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
)

type TournamentServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

func NewTournamentServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *TournamentServer {
	return &TournamentServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}
