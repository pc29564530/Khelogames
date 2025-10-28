package orchestrator

import (
	"khelogames/api/shared"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	"khelogames/api/transactions"
	db "khelogames/database"
	"khelogames/logger"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type CheckSportServer struct {
	store            *db.Store
	logger           *logger.Logger
	rabbitChan       *amqp091.Channel
	scoreBroadcaster shared.ScoreBroadcaster
	txStore          *transactions.SQLStore
}

func NewCheckSport(store *db.Store, logger *logger.Logger, scoreBroadcaster shared.ScoreBroadcaster) *CheckSportServer {
	return &CheckSportServer{store: store, logger: logger}
}

func (s *CheckSportServer) CheckSport(sports string, matches []db.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
	footballServer := football.NewFootballServer(s.store, s.logger, s.scoreBroadcaster, s.txStore)
	cricketServer := cricket.NewCricketServer(s.store, s.logger, s.rabbitChan, s.scoreBroadcaster, s.txStore)
	switch sports {
	case "cricket":
		return cricketServer.GetCricketScore(matches, tournamentPublicID)
	case "football":
		return footballServer.GetFootballScore(matches, tournamentPublicID)
	default:
		s.logger.Error("Unsupported sport type:", sports)
		return nil
	}
}
