package tournaments

import (
	"khelogames/api/shared"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
)

type TournamentServer struct {
	store            *db.Store
	logger           *logger.Logger
	tokenMaker       token.Maker
	config           util.Config
	scoreBroadcaster shared.ScoreBroadcaster
}

func NewTournamentServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, scoreBroadcaster shared.ScoreBroadcaster) *TournamentServer {
	return &TournamentServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, scoreBroadcaster: scoreBroadcaster}
}

func (s *TournamentServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

var _ shared.ScoreBroadcaster
