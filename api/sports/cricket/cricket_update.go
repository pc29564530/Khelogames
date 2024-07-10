package cricket

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CricketUpdateServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCricketUpdateServer(store *db.Store, logger *logger.Logger) *CricketUpdateServer {
	return &CricketUpdateServer{store: store, logger: logger}
}

type updateCricketMatchWicketRequest struct {
	Wickets int64 `json:"wickets"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketUpdateServer) UpdateCricketMatchWicketFunc(ctx *gin.Context) {

	var req updateCricketMatchWicketRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind wicket: %v", err)
		return
	}

	arg := db.UpdateCricketMatchWicketsParams{
		Wickets: req.Wickets,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketMatchWickets(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchRunsScoreRequest struct {
	Score   int64 `json:"score"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketUpdateServer) UpdateCricketMatchRunsScoreFunc(ctx *gin.Context) {

	var req updateCricketMatchRunsScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the cricket match score: %v", err)
		return
	}

	arg := db.UpdateCricketMatchRunsScoreParams{
		Score:   req.Score,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketMatchRunsScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update runs score : %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchExtrasRequest struct {
	Extras  int64 `json:"extras"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketUpdateServer) UpdateCricketMatchExtrasFunc(ctx *gin.Context) {

	var req updateCricketMatchExtrasRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the cricket extras: %v", err)
		return
	}

	arg := db.UpdateCricketMatchExtrasParams{
		Extras:  req.Extras,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketMatchExtras(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update cricket extras: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchInningsRequest struct {
	Innings int64 `json:"innings"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketUpdateServer) UpdateCricketMatchInningsFunc(ctx *gin.Context) {

	var req updateCricketMatchInningsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind cricket match innings ", err)
		return
	}

	arg := db.UpdateCricketMatchInningsParams{
		Innings: req.Innings,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketMatchInnings(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update cricket innings : %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchScoreBattingRequest struct {
	Position      string `json:"position"`
	RunsScored    string `json:"runs_scored"`
	BallsFaced    string `json:"balls_faced"`
	Fours         string `json:"fours"`
	Sixes         string `json:"sixes"`
	WicketTakenBy string `json:"wicket_taken_by"`
	TournamentID  int64  `json:"tournament_id"`
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	PlayerID      int64  `json:"player_id"`
}

func (s *CricketUpdateServer) UpdateCricketMatchScoreBattingFunc(ctx *gin.Context) {
	var req updateCricketMatchScoreBattingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var position int64
	if req.Position == "" {
		position = 0
	} else {
		position, err = strconv.ParseInt(req.Position, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse the position: %v", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var runsScored int64
	if req.RunsScored == " " {
		runsScored = 0
	} else {

		runsScored, err = strconv.ParseInt(req.RunsScored, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse runs scored %v", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var ballsFaced int64
	if req.BallsFaced == "" {
		ballsFaced = 0
	} else {
		ballsFaced, err = strconv.ParseInt(req.BallsFaced, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse balls faced: %v", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var fours int64
	if req.Fours == "" {
		fours = 0
	} else {
		fours, err = strconv.ParseInt(req.Fours, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse fours: %v", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var sixes int64
	if req.Sixes == "" {
		sixes = 0
	} else {
		sixes, err = strconv.ParseInt(req.Sixes, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse sixes: %v", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var wicketTakenBy int64
	if req.WicketTakenBy == "" {
		wicketTakenBy = 0
	} else {
		wicketTakenBy, err = strconv.ParseInt(req.WicketTakenBy, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse wicket taken by: %v", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	arg := db.UpdateCricketTeamPlayerScoreBattingParams{
		Position:      position,
		RunsScored:    runsScored,
		BallsFaced:    ballsFaced,
		Fours:         fours,
		Sixes:         sixes,
		WicketTakenBy: wicketTakenBy,
		TournamentID:  req.TournamentID,
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		PlayerID:      req.PlayerID,
	}

	response, err := s.store.UpdateCricketTeamPlayerScoreBatting(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update cricket team score batting: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchScoreBowlingRequest struct {
	OversBowled  string `json:"overs_bowled"`
	RunsConceded string `json:"runs_conceded"`
	WicketsTaken string `json:"wickets_taken"`
	TournamentID int64  `json:"tournament_id"`
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
	PlayerID     int64  `json:"player_id"`
}

func (s *CricketUpdateServer) UpdateCricketMatchScoreBowlingFunc(ctx *gin.Context) {
	var req updateCricketMatchScoreBowlingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var wicketsTaken int64
	if req.WicketsTaken == "" {
		wicketsTaken = 0
	} else {
		wicketsTaken, err = strconv.ParseInt(req.WicketsTaken, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse wicket taken: %v", err)
			return
		}
	}

	var runsConceded int64
	if req.RunsConceded == "" {
		runsConceded = 0
	} else {
		runsConceded, err = strconv.ParseInt(req.RunsConceded, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse runs conceded: %v", err)
			return
		}
	}

	arg := db.UpdateCricketTeamPlayerScoreBowlingParams{
		OversBowled:  req.OversBowled,
		RunsConceded: runsConceded,
		WicketsTaken: wicketsTaken,
		TournamentID: req.TournamentID,
		MatchID:      req.MatchID,
		TeamID:       req.TeamID,
		PlayerID:     req.PlayerID,
	}

	response, err := s.store.UpdateCricketTeamPlayerScoreBowling(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update cricket score bowling: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
