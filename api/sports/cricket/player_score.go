package cricket

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addCricketBatScore struct {
	BatsmanID  int64 `json:"batsman_id"`
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
	Position   int32 `json:"position"`
	RunsScored int32 `json:"runs_scored"`
	BallsFaced int32 `json:"balls_faced"`
	Fours      int32 `json:"fours"`
	Sixes      int32 `json:"sixes"`
}

func (s *CricketServer) AddCricketBatScoreFunc(ctx *gin.Context) {
	var req addCricketBatScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := db.AddCricketBatsScoreParams{
		BatsmanID:  req.BatsmanID,
		MatchID:    req.MatchID,
		TeamID:     req.TeamID,
		Position:   req.Position,
		RunsScored: req.RunsScored,
		BallsFaced: req.BallsFaced,
		Fours:      req.Fours,
		Sixes:      req.Sixes,
	}

	response, err := s.store.AddCricketBatsScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket player score: %v", gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

type addCricketBallScore struct {
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
	BowlerID   int64 `json:"bowler_id"`
	OverNumber int32 `json:"over_number"`
	BallNumber int32 `json:"ball_number"`
	Runs       int32 `json:"runs"`
	Wickets    int32 `json:"wickets"`
}

func (s *CricketServer) AddCricketBallFunc(ctx *gin.Context) {
	var req addCricketBallScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := db.AddCricketBallParams{
		MatchID:    req.MatchID,
		TeamID:     req.TeamID,
		BowlerID:   req.BowlerID,
		OverNumber: req.OverNumber,
		BallNumber: req.BallNumber,
		Runs:       req.Runs,
		Wickets:    req.Wickets,
	}

	response, err := s.store.AddCricketBall(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket player score: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type addCricketWicketScore struct {
	MatchID    int64  `json:"match_id"`
	TeamID     int64  `json:"team_id"`
	BatsmanID  int64  `json:"batsman_id"`
	BowlerID   int64  `json:"bowler_id"`
	FielderID  int64  `json:"fielder_id"`
	WicketType string `json:"wicket_type"`
}

func (s *CricketServer) AddCricketWicketFunc(ctx *gin.Context) {
	var req addCricketWicketScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.AddCricketWicketsParams{
		MatchID:    req.MatchID,
		BatsmanID:  req.BatsmanID,
		BowlerID:   req.BowlerID,
		FielderID:  req.FielderID,
		WicketType: req.WicketType,
	}

	response, err := s.store.AddCricketWickets(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket player score: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketBatRequest struct {
	BatsmanID  int64 `json:"batsman_id"`
	TeamID     int64 `json:"team_id"`
	MatchID    int64 `json:"match_id"`
	Position   int32 `json:"position"`
	RunsScored int32 `json:"runs_scored"`
	BallsFaced int32 `json:"balls_faced"`
	Fours      int32 `json:"fours"`
	Sixes      int32 `json:"sixes"`
}

func (s *CricketServer) UpdateCricketBatScoreFunc(ctx *gin.Context) {
	var req updateCricketBatRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.UpdateCricketRunsScoredParams{
		RunsScored: req.RunsScored,
		BallsFaced: req.BallsFaced,
		Fours:      req.Fours,
		Sixes:      req.Sixes,
		MatchID:    req.MatchID,
		BatsmanID:  req.BatsmanID,
	}

	response, err := s.store.UpdateCricketRunsScored(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket player score: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketBallRequest struct {
	OverNumber int32 `json:"over_number"`
	BallNumber int32 `json:"ball_number"`
	Runs       int32 `json:"runs"`
	Wickets    int32 `json:"wickets"`
	MatchID    int64 `json:"match_id"`
	BowlerID   int64 `json:"bowler_id"`
}

func (s *CricketServer) UpdateCricketBallFunc(ctx *gin.Context) {
	var req updateCricketBallRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.UpdateCricketBowlerParams{
		OverNumber: req.OverNumber,
		BallNumber: req.BallNumber,
		Runs:       req.Runs,
		Wickets:    req.Wickets,
		MatchID:    req.MatchID,
		BowlerID:   req.BowlerID,
	}

	response, err := s.store.UpdateCricketBowler(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket player score: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
