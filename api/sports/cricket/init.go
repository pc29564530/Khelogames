package cricket

import (
	shared "khelogames/api/shared"
	"khelogames/api/transactions"
	db "khelogames/database"
	"khelogames/logger"
)

type CricketServer struct {
	store            *db.Store
	logger           *logger.Logger
	scoreBroadcaster shared.ScoreBroadcaster
	txStore          *transactions.SQLStore
}

func NewCricketServer(store *db.Store, logger *logger.Logger, scoreBroadcaster shared.ScoreBroadcaster, txStore *transactions.SQLStore) *CricketServer {
	server := &CricketServer{
		store:            store,
		logger:           logger,
		scoreBroadcaster: scoreBroadcaster,
		txStore:          txStore,
	}

	return server
}

func (s *CricketServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

// GetScoreBroadcaster returns the assigned ScoreBroadcaster
func (s *CricketServer) GetScoreBroadcaster() shared.ScoreBroadcaster {
	return s.scoreBroadcaster
}

// Ensure CricketServer implements the shared interfaces
// var _ shared.CricketScoreUpdater

var _ shared.ScoreBroadcaster
