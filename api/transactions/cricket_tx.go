package transactions

import (
	"context"
	"fmt"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Add cricket toss and new inning score for batting team
func (store *SQLStore) AddCricketTossTx(ctx context.Context, matchPublicID uuid.UUID, tossDescision string, tossWinPublicID uuid.UUID) (models.CricketScore, models.Team, string, error) {
	var newScore models.CricketScore
	var team models.Team

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user
		response, err := q.AddCricketToss(ctx, matchPublicID, tossDescision, tossWinPublicID)
		if err != nil {
			store.logger.Error("Failed to add cricket match toss : ", err)
			return err
		}

		match, err := q.GetTournamentMatchByMatchID(ctx, matchPublicID)
		if err != nil {
			store.logger.Error("Failed to get the match by id: ", err)
			return err
		}

		var teamID int32
		if tossDescision == "batting" {
			teamID = response.TossWin
		} else {
			if match.HomeTeamID != response.TossWin {
				teamID = match.AwayTeamID
			} else {
				teamID = match.HomeTeamID
			}
		}
		team, err := q.GetTeamByID(ctx, int64(teamID))
		if err != nil {
			store.logger.Error("Failed to get team by id: ", err)
			return err
		}

		inningR := database.NewCricketScoreParams{
			MatchPublicID:     matchPublicID,
			TeamPublicID:      team.PublicID,
			InningNumber:      1,
			Score:             0,
			Wickets:           0,
			Overs:             0,
			RunRate:           "0.00",
			TargetRunRate:     "0.00",
			FollowOn:          false,
			IsInningCompleted: false,
			Declared:          false,
			InningStatus:      "not_started",
		}

		newScore, err = q.NewCricketScore(ctx, inningR)
		if err != nil {
			store.logger.Error("Failed to add the team score: ", err)
			return err
		}
		return err
	})
	return newScore, team, tossDescision, err
}

func (store *SQLStore) AddCricketBlowerTx(ctx context.Context,
	matchPublicID, teamPublicID, bowlerPublicID, prevBowlerPublicID uuid.UUID,
	inningNumber int,
) (models.BowlerScore, map[string]interface{}, error) {
	var prevBowlerID uuid.UUID
	var currentBowlerResponse models.BowlerScore
	var prevBowler map[string]interface{}
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		if prevBowlerPublicID != prevBowlerID {
			prevBowlerResponse, err := q.UpdateBowlingBowlerStatus(ctx, matchPublicID, bowlerPublicID, prevBowlerPublicID, inningNumber)
			if err != nil {
				store.logger.Error("Failed to update current bowler status: ", err)
				return err
			}

			playerData, err := q.GetPlayerByPublicID(ctx, bowlerPublicID)
			if err != nil {
				store.logger.Error("Failed to get Player: ", err)
			}
			prevBowler = map[string]interface{}{
				"player":            map[string]interface{}{"id": playerData.ID, "public_id": playerData.PublicID, "name": playerData.Name, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions},
				"id":                prevBowlerResponse.ID,
				"public_id":         prevBowlerResponse.PublicID,
				"match_id":          prevBowlerResponse.MatchID,
				"team_id":           prevBowlerResponse.TeamID,
				"bowler_id":         prevBowlerResponse.BowlerID,
				"runs":              prevBowlerResponse.Runs,
				"inning_number":     prevBowlerResponse.InningNumber,
				"ball_number":       prevBowlerResponse.BallNumber,
				"wide":              prevBowlerResponse.Wide,
				"no_ball":           prevBowlerResponse.NoBall,
				"wickets":           prevBowlerResponse.Wickets,
				"bowling_status":    prevBowlerResponse.BowlingStatus,
				"is_current_bowler": prevBowlerResponse.IsCurrentBowler,
			}
		}

		arg := database.AddCricketBallParams{
			MatchPublicID:  matchPublicID,
			TeamPublicID:   teamPublicID,
			BowlerPublicID: bowlerPublicID,
			InningNumber:   0,
			BallNumber:     0,
			Runs:           0,
			Wickets:        0,
			Wide:           0,
			NoBall:         0,
		}

		currentBowlerResponse, err = q.AddCricketBall(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to add the cricket bowler data: ", gin.H{"error": err.Error()})
			return err
		}
		return err
	})
	return currentBowlerResponse, prevBowler, err
}

func (store *SQLStore) UpdateCricketNoBallTx(ctx context.Context, matchPublicID, bowlerPublicID, battingTeamPublicID, batsmanPublicID uuid.UUID, runsScored int32, inningNumber int) (*models.BatsmanScore, []models.BatsmanScore, *models.BowlerScore, *models.CricketScore, error) {
	var bowlerScore *models.BowlerScore
	var inningScore *models.CricketScore
	var batsmanScore *models.BatsmanScore
	var currentBatsman []models.BatsmanScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		batsmanScore, bowlerScore, inningScore, err = q.UpdateNoBallsRuns(ctx, matchPublicID, battingTeamPublicID, bowlerPublicID, runsScored, inningNumber)
		if err != nil {
			store.logger.Error("Failed to update no_ball: ", err)
			return err
		}

		if bowlerScore.BallNumber%6 == 0 && runsScored%2 == 0 {
			currentBatsman, err = q.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
			if err != nil {
				store.logger.Error("Failed to update stricker: ", err)
			}
		} else if bowlerScore.BallNumber%6 != 0 && runsScored%2 != 0 {
			currentBatsman, err = q.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
			if err != nil {
				store.logger.Error("Failed to update stricker: ", err)
			}
		}
		return err
	})
	return batsmanScore, currentBatsman, bowlerScore, inningScore, err
}

func (store *SQLStore) UpdateWideRunsTx(ctx context.Context, matchPublicID, bowlerPublicID, battingTeamPublicID uuid.UUID, runsScored int32, inningNumber int) (models.BatsmanScore, []models.BatsmanScore, models.BowlerScore, models.CricketScore, error) {
	var inningScore *models.CricketScore
	var batsmanScore *models.BatsmanScore
	var bowlerScore *models.BowlerScore
	var currentBatsman []models.BatsmanScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		batsmanScore, bowlerScore, inningScore, err = q.UpdateWideRuns(ctx, matchPublicID, battingTeamPublicID, bowlerPublicID, runsScored, inningNumber)
		if err != nil {
			return fmt.Errorf("Failed to update wide runs: ", err)
		}

		if bowlerScore.BallNumber%6 == 0 && runsScored%2 == 0 {
			currentBatsman, err = q.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
			if err != nil {
				store.logger.Error("Failed to update stricker: ", err)
			}
		} else if bowlerScore.BallNumber%6 != 0 && runsScored%2 != 0 {
			currentBatsman, err = q.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
			if err != nil {
				store.logger.Error("Failed to update stricker: ", err)
			}
		}

		return err
	})
	return *batsmanScore, currentBatsman, *bowlerScore, *inningScore, err
}

func (store *SQLStore) UpdateCricketEndInningTx(ctx context.Context, matchID, batsmanTeamID int32, inningNumber int) (*models.CricketScore, *models.BatsmanScore, *models.BowlerScore, error) {
	var batsmanScore *models.BatsmanScore
	var bowlerScore *models.BowlerScore
	var inningScore *models.CricketScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		inningScore, batsmanScore, bowlerScore, err = q.UpdateInningEndStatus(ctx, matchID, batsmanTeamID, inningNumber)
		if err != nil {
			return err
		}
		return err
	})
	return inningScore, batsmanScore, bowlerScore, err
}

func (store *SQLStore) AddCricketBallTx(ctx context.Context, matchPublicID, teamPublicID, bowlerPublicID, prevBowlerPulbicID uuid.UUID, ballNumber, runs, wickets, wide, noBall int32, inningNumber int, bowlingStatus, isCurrentBowler bool) (models.BowlerScore, *models.BowlerScore, error) {
	var currentbowler models.BowlerScore
	var prevBowler *models.BowlerScore
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		prevBowlerPublicIDString := prevBowlerPulbicID.String()
		var prevBowlerEmptyString string

		if prevBowlerEmptyString != prevBowlerPublicIDString {
			prevBowler, err = q.UpdateBowlingBowlerStatus(ctx, matchPublicID, teamPublicID, bowlerPublicID, inningNumber)
			if err != nil {
				return err
			}
		}

		arg := database.AddCricketBallParams{
			MatchPublicID:  matchPublicID,
			TeamPublicID:   teamPublicID,
			BowlerPublicID: bowlerPublicID,
			InningNumber:   inningNumber,
			BallNumber:     ballNumber,
			Runs:           runs,
			Wickets:        wickets,
			Wide:           wide,
			NoBall:         noBall,
		}

		currentbowler, err = q.AddCricketBall(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to add the cricket bowler data: ", gin.H{"error": err.Error()})
			return err
		}
		return err
	})
	return currentbowler, prevBowler, err
}

// Add cricket wicket
func (store *SQLStore) AddCricketWicketTx(ctx context.Context,
	matchPublicID, battingTeamID, batsmanPublicID, bowlerPublicID uuid.UUID,
	wicketNumber int,
	wicketType string,
	ballNumber int,
	fielderPublicID uuid.UUID,
	inningScore models.CricketScore,
	bowlType *string,
	runsScored int32,
	inningNumber int,
	toggleStriker bool,
) (*models.BatsmanScore, *models.BatsmanScore, *models.BowlerScore, *models.CricketScore, *models.Wicket, error) {
	var outBatsmanResponse *models.BatsmanScore
	var notOutBatsmanResponse *models.BatsmanScore
	var bowlerResponse *models.BowlerScore
	var inningScoreResponse *models.CricketScore
	var currentBatsman *models.BatsmanScore
	var wicketResponse *models.Wicket
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		if bowlType != nil {
			outBatsmanResponse, notOutBatsmanResponse, bowlerResponse, inningScoreResponse, wicketResponse, err = q.AddCricketWicketWithBowlType(ctx, matchPublicID, battingTeamID, batsmanPublicID, bowlerPublicID, wicketNumber, wicketType, ballNumber, &fielderPublicID, inningScore.Score, *bowlType, inningNumber)
			if err != nil {
				store.logger.Error("failed to add cricket wicket with bowl type: ", err)
			}
		} else {
			outBatsmanResponse, notOutBatsmanResponse, bowlerResponse, inningScoreResponse, wicketResponse, err = q.AddCricketWicket(ctx, matchPublicID, battingTeamID, batsmanPublicID, bowlerPublicID, int(inningScore.Wickets), wicketType, int(inningScore.Overs), fielderPublicID, inningScore.Score, runsScored, inningNumber)
			if err != nil {
				store.logger.Error("failed to add cricket wicket: ", err)
				return err
			}
		}

		matchData, err := q.GetMatchModelByPublicId(ctx, matchPublicID)
		if err != nil {
			store.logger.Error("failed to get match: ", err)
			return err
		}

		if inningScoreResponse.Wickets == 10 {
			inningScoreResponse, notOutBatsmanResponse, bowlerResponse, err = q.UpdateInningEndStatus(ctx, int32(matchData.ID), notOutBatsmanResponse.TeamID, inningNumber)
			if err != nil {
				store.logger.Error("failed to update inning_numberscore: ", err)
				return err
			}
		} else if *matchData.MatchFormat == "T20" && inningScoreResponse.Overs/6 == 20 {
			inningScoreResponse, notOutBatsmanResponse, bowlerResponse, err = q.UpdateInningEndStatus(ctx, int32(matchData.ID), notOutBatsmanResponse.TeamID, inningNumber)
			if err != nil {
				store.logger.Error("failed to update inning_numberscore: ", err)
				return err
			}
		} else if *matchData.MatchFormat == "ODI" && inningScoreResponse.Overs/6 == 50 {
			inningScoreResponse, notOutBatsmanResponse, bowlerResponse, err = q.UpdateInningEndStatus(ctx, int32(matchData.ID), notOutBatsmanResponse.TeamID, inningNumber)
			if err != nil {
				store.logger.Error("failed to update inning_numberscore: ", err)
				return err
			}
		}

		if toggleStriker {
			notOut, err := q.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
			if err != nil {
				store.logger.Error("failed to toggle batsman: ", err)
				return err
			}
			notOutBatsmanResponse = &notOut[0]
		}
		currentBatsman = notOutBatsmanResponse
		if bowlerResponse.BallNumber%6 == 0 {
			currentBatsmanResponse, err := q.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
			if err != nil {
				store.logger.Error("Failed to update stricker: ", err)
			}
			currentBatsman = &currentBatsmanResponse[0]
		}
		return err
	})
	return outBatsmanResponse, currentBatsman, bowlerResponse, inningScoreResponse, wicketResponse, err
}

// Update bowling bowler status
func (store *SQLStore) UpdateBowlingBowlerStatusTx(ctx context.Context, matchPublicID, teamPublicID, currentBowlerPublicID, nextBowlerPublicID uuid.UUID, inningNumber int) (*models.BowlerScore, *models.BowlerScore, error) {
	var currentBowlerResponse *models.BowlerScore
	var nextBowlerResponse *models.BowlerScore
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		currentBowlerResponse, err = q.UpdateBowlingBowlerStatus(ctx, matchPublicID, teamPublicID, currentBowlerPublicID, inningNumber)
		if err != nil {
			store.logger.Error("Failed to update current bowler status: ", err)
			return err
		}

		nextBowlerResponse, err = q.UpdateBowlingBowlerStatus(ctx, matchPublicID, teamPublicID, nextBowlerPublicID, inningNumber)
		if err != nil {
			store.logger.Error("Failed to update next bowler status: ", err)
			return err
		}

		return err
	})
	return currentBowlerResponse, nextBowlerResponse, err
}

func (store *SQLStore) AddCricketSquadTx(ctx context.Context, playerData []map[string]interface{}, matchPublicID uuid.UUID, teamPublicID uuid.UUID, playerPublicID uuid.UUID) ([]map[string]interface{}, error) {
	var cricketSquad []map[string]interface{}
	err := store.execTx(ctx, func(q *database.Queries) error {
		for _, player := range playerData {
			var err error
			playerPublicID, err := uuid.Parse(player["public_id"].(string))
			if err != nil {
				store.logger.Error("Invalid UUID format", err)
				return fmt.Errorf("Invalid UUID format: ", err)
			}

			squad, err := q.AddCricketSquad(ctx, matchPublicID, teamPublicID, playerPublicID, player["position"].(string), player["on_bench"].(bool), false)
			if err != nil {
				return fmt.Errorf("Failed to add cricket squad: ", err)
			}

			cricketSquad = append(cricketSquad, map[string]interface{}{
				"id":         squad.ID,
				"public_id":  squad.PublicID,
				"match_id":   squad.MatchID,
				"team_id":    squad.TeamID,
				"player":     player,
				"role":       squad.Role,
				"on_bench":   squad.OnBench,
				"is_captain": squad.IsCaptain,
			})
		}
		return nil
	})
	return cricketSquad, err
}
