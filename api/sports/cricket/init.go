package cricket

import (
	shared "khelogames/api/shared"
	db "khelogames/database"
	"khelogames/logger"

	ampq "github.com/rabbitmq/amqp091-go"
)

type CricketServer struct {
	store            *db.Store
	logger           *logger.Logger
	rabbitChan       *ampq.Channel
	scoreBroadcaster shared.ScoreBroadcaster
}

func NewCricketServer(store *db.Store, logger *logger.Logger, rabbitChan *ampq.Channel, scoreBroadcaster shared.ScoreBroadcaster) *CricketServer {
	server := &CricketServer{
		store:            store,
		logger:           logger,
		rabbitChan:       rabbitChan,
		scoreBroadcaster: scoreBroadcaster,
	}

	return server
}

func (s *CricketServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

// Ensure CricketServer implements the shared interfaces
var _ shared.CricketScoreUpdater = (*CricketServer)(nil)
