package football

import (
	"khelogames/api/shared"
	db "khelogames/database"
	"khelogames/logger"
)

type FootballServer struct {
	store            *db.Store
	logger           *logger.Logger
	scoreBroadcaster shared.ScoreBroadcaster
}

func NewFootballServer(store *db.Store, logger *logger.Logger, scoreBroadcaster shared.ScoreBroadcaster) *FootballServer {
	return &FootballServer{store: store, logger: logger, scoreBroadcaster: scoreBroadcaster}
}

func (s *FootballServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

var _ shared.ScoreBroadcaster
