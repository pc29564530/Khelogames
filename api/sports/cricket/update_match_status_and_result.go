package cricket

import (
	"context"
	"khelogames/database/models"
)

func (s *CricketServer) UpdateMatchStatusAndResult(ctx context.Context, inningScore *models.CricketScore, matchData *models.Match, matchID int64) error {
	if inningScore.IsInningCompleted {
		matchInningScore, err := s.store.GetCricketScores(ctx, int32(matchID))
		if err != nil {
			s.logger.Error("Failed to both inning score: ", err)
			return err
		}
		if len(matchInningScore) == 2 {
			if matchInningScore[0].Score > matchInningScore[1].Score {
				updateMatchStatusResponse, err := s.store.UpdateMatchResult(ctx, int32(matchID), int32(matchInningScore[0].TeamID))
				if err != nil {
					s.logger.Error("Failed to update match result: ", err)
					return err
				}
				matchData.StatusCode = updateMatchStatusResponse.StatusCode
				matchData.Result = updateMatchStatusResponse.Result
			} else if matchInningScore[0].Score < matchInningScore[1].Score {
				updateMatchStatusResponse, err := s.store.UpdateMatchResult(ctx, int32(matchID), int32(matchInningScore[1].TeamID))
				if err != nil {
					s.logger.Error("Failed to update match result: ", err)
					return err
				}
				matchData.StatusCode = updateMatchStatusResponse.StatusCode
				matchData.Result = updateMatchStatusResponse.Result
			}
		} else if len(matchInningScore) == 4 {
			firstBatTeamScore := matchInningScore[0].Score + matchInningScore[2].Score
			secondBatTeamScore := matchInningScore[1].Score + matchInningScore[3].Score
			if firstBatTeamScore > secondBatTeamScore {
				updateMatchStatusResponse, err := s.store.UpdateMatchResult(ctx, int32(matchID), int32(matchInningScore[0].TeamID))
				if err != nil {
					s.logger.Error("Failed to update match result: ", err)
					return err
				}
				matchData.StatusCode = updateMatchStatusResponse.StatusCode
				matchData.Result = updateMatchStatusResponse.Result
			} else if firstBatTeamScore < secondBatTeamScore {
				updateMatchStatusResponse, err := s.store.UpdateMatchResult(ctx, int32(matchID), int32(matchInningScore[1].TeamID))
				if err != nil {
					s.logger.Error("Failed to update match result: ", err)
					return err
				}
				matchData.StatusCode = updateMatchStatusResponse.StatusCode
				matchData.Result = updateMatchStatusResponse.Result
			}
		}
	}
	return nil
}
