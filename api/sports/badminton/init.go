package badminton

import (
	shared "khelogames/api/shared"
	"khelogames/api/transactions"
	db "khelogames/database"
	"khelogames/logger"
)

type BadmintonServer struct {
	store            *db.Store
	logger           *logger.Logger
	scoreBroadcaster shared.ScoreBroadcaster
	txStore          *transactions.SQLStore
}

func NewBadmintonServer(store *db.Store, logger *logger.Logger, scoreBroadcaster shared.ScoreBroadcaster, txStore *transactions.SQLStore) *BadmintonServer {
	server := &BadmintonServer{
		store:            store,
		logger:           logger,
		scoreBroadcaster: scoreBroadcaster,
		txStore:          txStore,
	}

	return server
}

func (s *BadmintonServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

// GetScoreBroadcaster returns the assigned ScoreBroadcaster
func (s *BadmintonServer) GetScoreBroadcaster() shared.ScoreBroadcaster {
	return s.scoreBroadcaster
}

var _ shared.ScoreBroadcaster
