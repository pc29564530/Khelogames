package util

import (
	"khelogames/api/shared"
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	db "khelogames/database"
	"khelogames/logger"

	"github.com/google/uuid"
)

type CheckSportServer struct {
	store            *db.Store
	logger           *logger.Logger
	scoreBroadcaster shared.ScoreBroadcaster
}

func NewCheckSport(store *db.Store, logger *logger.Logger, scoreBroadcaster shared.ScoreBroadcaster) *CheckSportServer {
	return &CheckSportServer{store: store, logger: logger}
}

func (s *CheckSportServer) CheckSport(sports string, matches []db.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
	footballServer := football.NewFootballServer(s.store, s.logger, s.scoreBroadcaster)
	cricketServer := cricket.NewCricketServer(s.store, s.logger, nil, nil)
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
