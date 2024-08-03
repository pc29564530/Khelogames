package football

import (
	"context"
	db "khelogames/db/sqlc"
)

func (s *FootballServer) GetFootballMatchScore(matches []db.Match, matchDetails []map[string]interface{}) []map[string]interface{} {
	s.logger.Info("Get football match score")
	ctx := context.Background()
	for i, match := range matches {
		arg1 := db.GetFootballScoreParams{MatchID: match.ID, TeamID: match.HomeTeamID}
		arg2 := db.GetFootballScoreParams{MatchID: match.ID, TeamID: match.AwayTeamID}
		s.logger.Debug("football match arg ", arg1)
		s.logger.Debug("football match arg ", arg1)
		matchScoreData1, err := s.store.GetFootballScore(ctx, arg1)
		if err != nil {
			s.logger.Error("Failed to get football match score for team 1:", err)
			return nil
		}
		matchScoreData2, err := s.store.GetFootballScore(ctx, arg2)
		if err != nil {
			s.logger.Error("Failed to get football match score for team 2:", err)
			return nil
		}
		matchDetails[i]["home_score"] = matchScoreData1.Goals
		matchDetails[i]["away_score"] = matchScoreData2.Goals

	}
	return matchDetails
}
