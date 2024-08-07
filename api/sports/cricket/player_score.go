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
	MatchID  int64 `json:"match_id"`
	TeamID   int64 `json:"team_id"`
	BowlerID int64 `json:"bowler_id"`
	Ball     int32 `json:"ball"`
	Runs     int32 `json:"runs"`
	Wickets  int32 `json:"wickets"`
	Wide     int32 `json:"wide"`
	NoBall   int32 `json:"no_ball"`
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
		MatchID:  req.MatchID,
		TeamID:   req.TeamID,
		BowlerID: req.BowlerID,
		Ball:     req.Ball,
		Runs:     req.Runs,
		Wickets:  req.Wickets,
		Wide:     req.Wide,
		NoBall:   req.NoBall,
	}

	response, err := s.store.AddCricketBall(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket bowler data: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type addCricketWicketScore struct {
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	BatsmanID     int64  `json:"batsman_id"`
	BowlerID      int64  `json:"bowler_id"`
	WicketsNumber int32  `json:"wickets_number"`
	WicketType    string `json:"wicket_type"`
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
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		BatsmanID:     req.BatsmanID,
		BowlerID:      req.BowlerID,
		WicketsNumber: req.WicketsNumber,
		WicketType:    req.WicketType,
	}

	response, err := s.store.AddCricketWickets(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket wicket: %v", gin.H{"error": err.Error()})
		return
	}

	var updageCricketWickets *db.Wicket

	if updageCricketWickets != nil {
		arg := db.UpdateCricketWicketsParams{
			MatchID: req.MatchID,
			TeamID:  req.TeamID,
		}

		_, err := s.store.UpdateCricketWickets(ctx, arg)
		if err != nil {
			s.logger.Error("Failed to update the cricket wicket: %v", gin.H{"error": err.Error()})
			return
		}
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

	s.logger.Debug("successfully bind :", req)

	arg := db.UpdateCricketRunsScoredParams{
		RunsScored: req.RunsScored,
		BallsFaced: req.BallsFaced,
		Fours:      req.Fours,
		Sixes:      req.Sixes,
		MatchID:    req.MatchID,
		BatsmanID:  req.BatsmanID,
		TeamID:     req.TeamID,
	}

	response, err := s.store.UpdateCricketRunsScored(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the cricket player runs: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketBallRequest struct {
	Ball     int32 `json:"ball"`
	Runs     int32 `json:"runs"`
	Wickets  int32 `json:"wickets"`
	Wide     int32 `json:"wide"`
	NoBall   int32 `json:""no_ball`
	MatchID  int64 `json:"match_id"`
	BowlerID int64 `json:"bowler_id"`
	TeamID   int64 `json:"team_id"`
}

func (s *CricketServer) UpdateCricketBallFunc(ctx *gin.Context) {
	var req updateCricketBallRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Debug("successfully bind: ", req)

	arg := db.UpdateCricketBowlerParams{
		Ball:     req.Ball,
		Runs:     req.Runs,
		Wickets:  req.Wickets,
		Wide:     req.Wide,
		NoBall:   req.NoBall,
		MatchID:  req.MatchID,
		BowlerID: req.BowlerID,
		TeamID:   req.TeamID,
	}

	response, err := s.store.UpdateCricketBowler(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the cricket player bowler: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getPlayerScoreRequest struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketServer) GetPlayerScoreFunc(ctx *gin.Context) {
	var req getPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketPlayersScoreParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.GetCricketPlayersScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get players score : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	argCricketScore := db.UpdateCricketScoreParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	_, err = s.store.UpdateCricketScore(ctx, argCricketScore)
	if err != nil {
		s.logger.Error("Failed to update the cricket score: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type getCricketBowlerRequest struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketServer) GetCricketBowlerFunc(ctx *gin.Context) {
	var req getCricketBowlerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketBallsParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}
	response, err := s.store.GetCricketBalls(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get cricket bowler score : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	matchResponse, err := s.store.GetMatchByMatchID(ctx, arg.MatchID)
	if err != nil {
		s.logger.Error("Failed to update the cricket player bowler: %v", gin.H{"error": err.Error()})
		return
	}

	var currentID int64
	if matchResponse.AwayTeamID != arg.TeamID {
		currentID = matchResponse.AwayTeamID
	} else {
		currentID = matchResponse.HomeTeamID
	}
	arg1 := db.UpdateCricketOversParams{
		MatchID: req.MatchID,
		TeamID:  currentID,
	}

	_, err = s.store.UpdateCricketOvers(ctx, arg1)
	if err != nil {
		s.logger.Error("Failed to add the cricket overs: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type getCricketWicketsRequest struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketServer) GetCricketWicketsFunc(ctx *gin.Context) {
	var req getCricketWicketsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketWicketsParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}
	s.logger.Debug("cricket wicket arg: ", arg)
	response, err := s.store.GetCricketWickets(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get cricket bowler score : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Successfully get the wickets: ", response)

	argCricketTeamWicket := db.UpdateCricketWicketsParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	updateResponse, err := s.store.UpdateCricketWickets(ctx, argCricketTeamWicket)
	if err != nil {
		s.logger.Error("Failed to upate cricket wicket : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Debug("Successfully update the wickets: ", updateResponse)

	ctx.JSON(http.StatusAccepted, response)
}
