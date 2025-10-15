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

	// Helper components
	// errorHandler    *ErrorHandler
	// playerHelper    *PlayerHelper
	// scoreCalculator *ScoreCalculator
	// broadcastHelper *BroadcastHelper
}

func NewCricketServer(store *db.Store, logger *logger.Logger, rabbitChan *ampq.Channel, scoreBroadcaster shared.ScoreBroadcaster) *CricketServer {
	server := &CricketServer{
		store:            store,
		logger:           logger,
		rabbitChan:       rabbitChan,
		scoreBroadcaster: scoreBroadcaster,
	}

	// Initialize helper components
	// server.errorHandler = NewErrorHandler(logger)
	// server.playerHelper = NewPlayerHelper(store, logger)
	// server.scoreCalculator = NewScoreCalculator(logger)
	// server.broadcastHelper = NewBroadcastHelper(scoreBroadcaster, logger)

	return server
}

// Ensure CricketServer implements the shared interfaces
// var _ shared.CricketScoreUpdater = (*CricketServer)(nil)
