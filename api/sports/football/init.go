package football

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
)

type FootballServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewFootballServer(store *db.Store, logger *logger.Logger) *FootballServer {
	return &FootballServer{store: store, logger: logger}
}
