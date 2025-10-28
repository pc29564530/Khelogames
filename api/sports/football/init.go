package football

import (
	"khelogames/api/shared"
	"khelogames/api/transactions"

	db "khelogames/database"
	"khelogames/logger"
)

type FootballServer struct {
	store            *db.Store
	logger           *logger.Logger
	scoreBroadcaster shared.ScoreBroadcaster
	txStore          *transactions.SQLStore
}

func NewFootballServer(store *db.Store, logger *logger.Logger, scoreBroadcaster shared.ScoreBroadcaster, txStore *transactions.SQLStore) *FootballServer {
	return &FootballServer{store: store, logger: logger, scoreBroadcaster: scoreBroadcaster, txStore: txStore}
}

func (s *FootballServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

var _ shared.ScoreBroadcaster
