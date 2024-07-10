package football

import (
	"context"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
)

type FootballMatchServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewFootballMatchServer(store *db.Store, logger *logger.Logger) *FootballMatchServer {
	return &FootballMatchServer{store: store, logger: logger}
}

func (s *FootballMatchServer) GetFootballMatchScore(matches []db.TournamentMatch, matchDetails []map[string]interface{}) []map[string]interface{} {
	s.logger.Info("Get football match score")
	ctx := context.Background()
	for _, match := range matches {
		arg1 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team1ID, TournamentID: match.TournamentID}
		arg2 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team2ID, TournamentID: match.TournamentID}
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

		fmt.Println("Score2: ", matchScoreData2)
		matchDetail := map[string]interface{}{
			"team1_score": matchScoreData1.GoalFor,
			"team2_score": matchScoreData2.GoalFor,
		}

		matchDetails = append(matchDetails, matchDetail)

	}
	return matchDetails
}
