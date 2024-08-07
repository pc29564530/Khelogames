package football

import (
	"context"
	db "khelogames/db/sqlc"
)

func (s *FootballServer) GetFootballMatchScore(matches []db.TournamentMatch, matchDetails []map[string]interface{}) []map[string]interface{} {
	s.logger.Info("Get football match score")
	ctx := context.Background()
	for i, match := range matches {
		arg1 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team1ID, TournamentID: match.TournamentID}
		arg2 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team2ID, TournamentID: match.TournamentID}
		s.logger.Debug("football match arg ", arg1)
		s.logger.Debug("football match arg ", arg1)
		matchScoreData1, err := s.store.GetFootballMatchScore(ctx, arg1)
		if err != nil {
			s.logger.Error("Failed to get football match score for team 1:", err)
			return nil
		}
		matchScoreData2, err := s.store.GetFootballMatchScore(ctx, arg2)
		if err != nil {
			s.logger.Error("Failed to get football match score for team 2:", err)
			return nil
		}
		matchDetails[i]["team1_score"] = matchScoreData1.GoalFor
		matchDetails[i]["team2_score"] = matchScoreData2.GoalFor

	}
	return matchDetails
}
