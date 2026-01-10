package cricket

// import (
// 	"context"
// 	"fmt"
// 	shared "khelogames/api/shared"
// 	db "khelogames/database"
// 	"khelogames/database/models"
// 	"khelogames/logger"

// 	"github.com/google/uuid"
// )

// // PlayerHelper handles player-related operations
// type PlayerHelper struct {
// 	store  *db.Store
// 	logger *logger.Logger
// }

// func NewPlayerHelper(store *db.Store, logger *logger.Logger) *PlayerHelper {
// 	return &PlayerHelper{store: store, logger: logger}
// }

// // GetPlayerInfo retrieves player information and converts to PlayerInfo struct
// func (ph *PlayerHelper) GetPlayerInfo(ctx context.Context, playerID int64) (*PlayerInfo, error) {
// 	player, err := ph.store.GetPlayerByID(ctx, playerID)
// 	if err != nil {
// 		ph.logger.Error("Failed to get player", map[string]interface{}{
// 			"player_id": playerID,
// 			"error":     err.Error(),
// 		})
// 		return nil, fmt.Errorf("failed to get player: %w", err)
// 	}

// 	return &PlayerInfo{
// 		ID:        player.ID,
// 		PublicID:  player.PublicID,
// 		Name:      player.Name,
// 		Slug:      player.Slug,
// 		ShortName: player.ShortName,
// 		Position:  player.Positions,
// 	}, nil
// }

// // GetPlayerInfoByPublicID retrieves player information by public ID
// func (ph *PlayerHelper) GetPlayerInfoByPublicID(ctx context.Context, publicID uuid.UUID) (*PlayerInfo, error) {
// 	player, err := ph.store.GetPlayerByPublicID(ctx, publicID)
// 	if err != nil {
// 		ph.logger.Error("Failed to get player by public ID", map[string]interface{}{
// 			"public_id": publicID,
// 			"error":     err.Error(),
// 		})
// 		return nil, fmt.Errorf("failed to get player: %w", err)
// 	}

// 	return &PlayerInfo{
// 		ID:        player.ID,
// 		PublicID:  player.PublicID,
// 		Name:      player.Name,
// 		Slug:      player.Slug,
// 		ShortName: player.ShortName,
// 		Position:  player.Positions,
// 	}, nil
// }

// // ConvertBatsmanScore converts database model to response struct
// func (ph *PlayerHelper) ConvertBatsmanScore(ctx context.Context, score models.BatsmanScore) (*BatsmanScore, error) {
// 	playerInfo, err := ph.GetPlayerInfo(ctx, int64(score.BatsmanID))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &BatsmanScore{
// 		Player:             *playerInfo,
// 		ID:                 score.ID,
// 		PublicID:           score.PublicID,
// 		MatchID:            score.MatchID,
// 		TeamID:             score.TeamID,
// 		BatsmanID:          score.BatsmanID,
// 		RunsScored:         score.RunsScored,
// 		BallsFaced:         score.BallsFaced,
// 		Fours:              score.Fours,
// 		Sixes:              score.Sixes,
// 		BattingStatus:      score.BattingStatus,
// 		IsStriker:          score.IsStriker,
// 		IsCurrentlyBatting: score.IsCurrentlyBatting,
// 		InningNumber:       score.InningNumber,
// 	}, nil
// }

// // ConvertBowlerScore converts database model to response struct
// func (ph *PlayerHelper) ConvertBowlerScore(ctx context.Context, score models.BowlerScore) (*BowlerScore, error) {
// 	playerInfo, err := ph.GetPlayerInfo(ctx, int64(score.BowlerID))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &BowlerScore{
// 		Player:          *playerInfo,
// 		ID:              score.ID,
// 		PublicID:        score.PublicID,
// 		MatchID:         score.MatchID,
// 		TeamID:          score.TeamID,
// 		BowlerID:        score.BowlerID,
// 		BallNumber:      score.BallNumber,
// 		Runs:            score.Runs,
// 		Wide:            score.Wide,
// 		NoBall:          score.NoBall,
// 		Wickets:         score.Wickets,
// 		BowlingStatus:   score.BowlingStatus,
// 		IsCurrentBowler: score.IsCurrentBowler,
// 		InningNumber:    score.InningNumber,
// 	}, nil
// }

// // ConvertInningScore converts database model to response struct
// func (ph *PlayerHelper) ConvertInningScore(score models.CricketScore) *InningScore {
// 	return &InningScore{
// 		ID:                score.ID,
// 		PublicID:          score.PublicID,
// 		MatchID:           score.MatchID,
// 		TeamID:            score.TeamID,
// 		InningNumber:      score.InningNumber,
// 		Score:             score.Score,
// 		Wickets:           score.Wickets,
// 		Overs:             score.Overs,
// 		RunRate:           score.RunRate,
// 		TargetRunRate:     score.TargetRunRate,
// 		FollowOn:          score.FollowOn,
// 		IsInningCompleted: score.IsInningCompleted,
// 		Declared:          score.Declared,
// 	}
// }

// // ScoreCalculator handles score-related calculations
// type ScoreCalculator struct {
// 	logger *logger.Logger
// }

// func NewScoreCalculator(logger *logger.Logger) *ScoreCalculator {
// 	return &ScoreCalculator{logger: logger}
// }

// // CalculateRunRate calculates the run rate for an inning
// func (sc *ScoreCalculator) CalculateRunRate(runs, overs int32) string {
// 	if overs == 0 {
// 		return "0.00"
// 	}

// 	rate := float64(runs) / float64(overs)
// 	return fmt.Sprintf("%.2f", rate)
// }

// // CalculateStrikeRate calculates the strike rate for a batsman
// func (sc *ScoreCalculator) CalculateStrikeRate(runs, balls int32) string {
// 	if balls == 0 {
// 		return "0.00"
// 	}

// 	rate := float64(runs) / float64(balls) * 100
// 	return fmt.Sprintf("%.2f", rate)
// }

// // ShouldRotateStriker determines if the striker should be rotated
// func (sc *ScoreCalculator) ShouldRotateStriker(ballNumber, runsScored int32) bool {
// 	// Rotate if last ball of over (ball 6) and even runs
// 	if ballNumber%6 == 0 && runsScored%2 == 0 {
// 		return true
// 	}
// 	// Rotate if not last ball and odd runs
// 	if ballNumber%6 != 0 && runsScored%2 != 0 {
// 		return true
// 	}
// 	return false
// }

// // IsInningCompleted checks if an inning should be completed
// func (sc *ScoreCalculator) IsInningCompleted(wickets, overs int32, matchFormat string) bool {
// 	// 10 wickets fallen
// 	if wickets >= 10 {
// 		return true
// 	}

// 	// Overs completed based on match format
// 	switch matchFormat {
// 	case "T20":
// 		return overs >= 120 // 20 overs * 6 balls
// 	case "ODI":
// 		return overs >= 300 // 50 overs * 6 balls
// 	case "Test":
// 		// Test matches don't have over limits, only wicket limits
// 		return wickets >= 10
// 	default:
// 		// Default to T20 format
// 		return overs >= 120
// 	}
// }

// // BroadcastHelper handles broadcasting operations
// type BroadcastHelper struct {
// 	broadcaster shared.ScoreBroadcaster
// 	logger      *logger.Logger
// }

// func NewBroadcastHelper(broadcaster shared.ScoreBroadcaster, logger *logger.Logger) *BroadcastHelper {
// 	return &BroadcastHelper{
// 		broadcaster: broadcaster,
// 		logger:      logger,
// 	}
// }

// // BroadcastScoreUpdate broadcasts a score update event
// func (bh *BroadcastHelper) BroadcastScoreUpdate(ctx context.Context, eventType string, payload interface{}) error {
// 	if bh.broadcaster == nil {
// 		bh.logger.Warn("No broadcaster available, skipping broadcast")
// 		return nil
// 	}

// 	err := bh.broadcaster.BroadcastCricketEvent(ctx, eventType, payload.(map[string]interface{}))
// 	if err != nil {
// 		bh.logger.Error("Failed to broadcast score update", map[string]interface{}{
// 			"event_type": eventType,
// 			"error":      err.Error(),
// 		})
// 		return fmt.Errorf("failed to broadcast score update: %w", err)
// 	}

// 	bh.logger.Debug("Successfully broadcasted score update", map[string]interface{}{
// 		"event_type": eventType,
// 	})
// 	return nil
// }



