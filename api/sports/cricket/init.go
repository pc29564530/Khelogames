package cricket

import (
	shared "khelogames/api/shared"
	"khelogames/api/transactions"
	db "khelogames/database"
	"khelogames/logger"

	ampq "github.com/rabbitmq/amqp091-go"
)

type CricketServer struct {
	store            *db.Store
	logger           *logger.Logger
	rabbitChan       *ampq.Channel
	scoreBroadcaster shared.ScoreBroadcaster
	txStore          *transactions.SQLStore
}

func NewCricketServer(store *db.Store, logger *logger.Logger, rabbitChan *ampq.Channel, scoreBroadcaster shared.ScoreBroadcaster, txStore *transactions.SQLStore) *CricketServer {
	server := &CricketServer{
		store:            store,
		logger:           logger,
		rabbitChan:       rabbitChan,
		scoreBroadcaster: scoreBroadcaster,
		txStore:          txStore,
	}

	return server
}

func (s *CricketServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

// Ensure CricketServer implements the shared interfaces
var _ shared.CricketScoreUpdater = (*CricketServer)(nil)
