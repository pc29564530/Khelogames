package cricket

import (
	"context"
	"khelogames/database/models"
)

func (s *CricketServer) UpdateMatchStatusAndResult(ctx context.Context, inningScore *models.CricketScore, matchData *models.Match, matchID int64) error {
	if !inningScore.IsInningCompleted {
		return nil
	}

	matchInningScore, err := s.store.GetCricketScores(ctx, int32(matchID))
	if err != nil {
		s.logger.Error("Failed to get both inning scores: ", err)
		return err
	}

	if len(matchInningScore) == 2 {
		// ODI / T20 format: 2 innings total
		if !matchInningScore[0].IsInningCompleted || !matchInningScore[1].IsInningCompleted {
			return nil
		}
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
		} else {
			// Scores are equal - it's a tie (result=0 means no winner)
			updateMatchStatusResponse, err := s.store.UpdateMatchResult(ctx, int32(matchID), 0)
			if err != nil {
				s.logger.Error("Failed to update match result (tie): ", err)
				return err
			}
			matchData.StatusCode = updateMatchStatusResponse.StatusCode
			matchData.Result = updateMatchStatusResponse.Result
		}
	} else if len(matchInningScore) == 4 {
		// Adding the test functionality in future
		return nil
	}

	return nil
}
