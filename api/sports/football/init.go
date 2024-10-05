package football

import (
	db "khelogames/database"
	"khelogames/logger"
)

type FootballServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewFootballServer(store *db.Store, logger *logger.Logger) *FootballServer {
	return &FootballServer{store: store, logger: logger}
}
