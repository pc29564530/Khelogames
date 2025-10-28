package tournaments

import (
	"khelogames/api/shared"
	"khelogames/api/transactions"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/logger"
	"khelogames/util"
)

type TournamentServer struct {
	store            *db.Store
	logger           *logger.Logger
	tokenMaker       token.Maker
	config           util.Config
	scoreBroadcaster shared.ScoreBroadcaster
	txStore          *transactions.SQLStore
}

func NewTournamentServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config, scoreBroadcaster shared.ScoreBroadcaster, txStore *transactions.SQLStore) *TournamentServer {
	return &TournamentServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config, scoreBroadcaster: scoreBroadcaster, txStore: txStore}
}

func (s *TournamentServer) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

var _ shared.ScoreBroadcaster
