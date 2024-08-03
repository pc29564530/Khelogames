package util

import (
	"khelogames/api/sports/cricket"
	"khelogames/api/sports/football"
	db "khelogames/db/sqlc"
	"khelogames/logger"
)

type CheckSportServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCheckSport(store *db.Store, logger *logger.Logger) *CheckSportServer {
	return &CheckSportServer{store: store, logger: logger}
}

func (s *CheckSportServer) CheckSport(sports string, matches []db.Match, matchDetails []map[string]interface{}) []map[string]interface{} {
	footballServer := football.NewFootballServer(s.store, s.logger)
	cricketServer := cricket.NewCricketServer(s.store, s.logger)
	switch sports {
	case "Cricket":
		return cricketServer.GetCricketScore(matches, matchDetails)
	case "Football":
		return footballServer.GetFootballScore(matches, matchDetails)
	default:
		s.logger.Error("Unsupported sport type:", sports)
		return nil
	}
}
