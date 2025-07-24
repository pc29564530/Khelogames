package util

import (
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	db "khelogames/database"
	"khelogames/logger"

	"github.com/google/uuid"
)

type CheckSportServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCheckSport(store *db.Store, logger *logger.Logger) *CheckSportServer {
	return &CheckSportServer{store: store, logger: logger}
}

func (s *CheckSportServer) CheckSport(sports string, matches []db.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
	footballServer := football.NewFootballServer(s.store, s.logger)
	cricketServer := cricket.NewCricketServer(s.store, s.logger)
	switch sports {
	case "cricket":
		return cricketServer.GetCricketScore(matches, tournamentID)
	case "football":
		return footballServer.GetFootballScore(matches, tournamentID)
	default:
		s.logger.Error("Unsupported sport type:", sports)
		return nil
	}
}
