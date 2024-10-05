package cricket

import (
	db "khelogames/database"
	"khelogames/logger"
)

type CricketServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCricketServer(store *db.Store, logger *logger.Logger) *CricketServer {
	return &CricketServer{store: store, logger: logger}
}
