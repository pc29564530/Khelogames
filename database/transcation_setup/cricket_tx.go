package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (store *SQLStore) AddCricketTossTx(ctx context.Context, matchPublicID uuid.UUID, tossDescision string, tossWinPublicID uuid.UUID) (models.CricketToss, error) {
	var cricketToss models.CricketToss

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		cricketToss, err = q.AddCricketToss(ctx, matchPublicID, tossDescision, tossWinPublicID)
		if err != nil {
			return err
		}
		return err
	})
	return cricketToss, err
}

func (store *SQLStore) AddCricketScoreTx(ctx context.Context, arg database.NewCricketScoreParams) (models.CricketScore, error) {
	var inningScore models.CricketScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		inningScore, err = q.NewCricketScore(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return inningScore, err
}

func (store *SQLStore) UpdateCricketNoBallTx(ctx context.Context, matchPublicID, bowlerPublicID, battingTeamPublicID uuid.UUID, runsScored int32, inningNumber int) (*models.BatsmanScore, *models.BowlerScore, *models.CricketScore, error) {
	var bowlerScore *models.BowlerScore
	var inningScore *models.CricketScore
	var batsmanScore *models.BatsmanScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		batsmanScore, bowlerScore, inningScore, err = q.UpdateNoBallsRuns(ctx, matchPublicID, bowlerPublicID, battingTeamPublicID, runsScored, inningNumber)
		if err != nil {
			return err
		}
		return err
	})
	return batsmanScore, bowlerScore, inningScore, err
}

func (store *SQLStore) UpdateWideRuns(ctx context.Context, matchPublicID, bowlerPublicID, battingTeamPublicID uuid.UUID, runsScored int32, inningNumber int) (models.BatsmanScore, models.BowlerScore, models.CricketScore, error) {
	var inningScore *models.CricketScore
	var batsmanScore *models.BatsmanScore
	var bowlerScore *models.BowlerScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		batsmanScore, bowlerScore, inningScore, err = q.UpdateWideRuns(ctx, matchPublicID, bowlerPublicID, battingTeamPublicID, runsScored, inningNumber)
		if err != nil {
			return err
		}
		return err
	})
	return *batsmanScore, *bowlerScore, *inningScore, err
}

func (store *SQLStore) ToggleCricketStrikerTX(ctx context.Context, matchPublicID uuid.UUID, inningNumber int) ([]models.BatsmanScore, error) {
	var batsmanScore []models.BatsmanScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		batsmanScore, err = q.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			return err
		}
		return err
	})
	return batsmanScore, err
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

func (store *SQLStore) AddCricketBatsScoreTx(ctx context.Context, arg database.AddCricketBatsScoreParams) (models.BatsmanScore, error) {
	var batsmanScore models.BatsmanScore

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		batsmanScore, err = q.AddCricketBatsScore(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return batsmanScore, err
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
			MatchPublicID:   matchPublicID,
			TeamPublicID:    teamPublicID,
			BowlerPublicID:  bowlerPublicID,
			BallNumber:      ballNumber,
			Runs:            runs,
			Wickets:         wickets,
			Wide:            wide,
			NoBall:          noBall,
			BowlingStatus:   bowlingStatus,
			IsCurrentBowler: isCurrentBowler,
			InningNumber:    inningNumber,
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

func (store *SQLStore) AddCricketWicketTx(ctx context.Context, matchPublicID, battingTeamID, batsmanPublicID, bowlerPublicID uuid.UUID, wicketNumber int, wicketType string, ballNumber int, fielderPublicID uuid.UUID, inningScore models.CricketScore, bowlType *string, runsScored int32, inningNumber int, toggleStriker bool) (models.BatsmanScore, *models.BatsmanScore, *models.BatsmanScore, *models.BowlerScore, *models.CricketScore, *models.Wicket, error) {
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
		// this function is undefined now fix it
		err = UpdateMatchStatusAndResult(ctx, inningScoreResponse, matchData, matchData.ID)
		if err != nil {
			store.logger.Error("Failed to update match status and result: ", err)
			return err
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
	return *outBatsmanResponse, notOutBatsmanResponse, currentBatsman, bowlerResponse, inningScoreResponse, wicketResponse, err
}

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
