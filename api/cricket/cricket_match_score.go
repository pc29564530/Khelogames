package cricket

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CricketMatchScoreServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCricketMatchScore(store *db.Store, logger *logger.Logger) *CricketMatchScoreServer {
	return &CricketMatchScoreServer{store: store, logger: logger}
}

type addCricketMatchScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
	Score        int64 `json:"score"`
	Wickets      int64 `json:"wickets"`
	Overs        int64 `json:"overs"`
	Extras       int64 `json:"extras"`
	Innings      int64 `json:"innings"`
}

func (s *CricketMatchScoreServer) AddCricketMatchScoreFunc(ctx *gin.Context) {

	var req addCricketMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateCricketMatchScoreParams{
		MatchID:      req.MatchID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
		Score:        req.Score,
		Wickets:      req.Wickets,
		Overs:        req.Overs,
		Extras:       req.Extras,
		Innings:      req.Innings,
	}

	response, err := s.store.CreateCricketMatchScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

// func (s *CricketMatchScoreServer) getCricketMatchScore(ctx *gin.Context) {
// 	matchIdStr := ctx.Query("match_id")
// 	matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
// 	if err != nil {
// 		ctx.Error(err)
// 		return
// 	}

// 	teamIdStr := ctx.Query("team_id")
// 	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
// 	if err != nil {
// 		ctx.Error(err)
// 		return
// 	}

// 	arg := db.GetCricketMatchScoreParams{
// 		MatchID: matchId,
// 		TeamID:  teamID,
// 	}

// 	response, err := s.store.GetCricketMatchScore(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return

// }

type updateCricketMatchWicketRequest struct {
	Wickets int64 `json:"wickets"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *CricketMatchScoreServer) UpdateCricketMatchWicketFunc(ctx *gin.Context) {

	var req updateCricketMatchWicketRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind wicket: %v", err)
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

func (s *CricketMatchScoreServer) UpdateCricketMatchRunsScoreFunc(ctx *gin.Context) {

	var req updateCricketMatchRunsScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind the cricket match score: %v", err)
		return
	}

	arg := db.UpdateCricketMatchRunsScoreParams{
		Score:   req.Score,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketMatchRunsScore(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update runs score : %v", err)
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

func (s *CricketMatchScoreServer) UpdateCricketMatchExtrasFunc(ctx *gin.Context) {

	var req updateCricketMatchExtrasRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind the cricket extras: %v", err)
		return
	}

	arg := db.UpdateCricketMatchExtrasParams{
		Extras:  req.Extras,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketMatchExtras(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update cricket extras: %v", err)
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

func (s *CricketMatchScoreServer) UpdateCricketMatchInningsFunc(ctx *gin.Context) {

	var req updateCricketMatchInningsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind cricket match innings ", err)
		return
	}

	arg := db.UpdateCricketMatchInningsParams{
		Innings: req.Innings,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketMatchInnings(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update cricket innings : %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

//cricket player score

type addCricketTeamPlayerScoreRequest struct {
	MatchID          int64  `json:"match_id"`
	TournamentID     int64  `json:"tournament_id"`
	TeamID           int64  `json:"team_id"`
	BattingOrBowling string `json:"batting_or_bowling, omit_empty"`
	Position         string `json:"position"`
	PlayerID         int64  `json:"player_id"`
	RunsScored       string `json:"runs_scored"`
	BallsFaced       string `json:"balls_faced"`
	Fours            string `json:"fours"`
	Sixes            string `json:"sixes"`
	WicketsTaken     string `json:"wickets_taken"`
	OversBowled      string `json:"overs_bowled"`
	RunsConceded     string `json:"runs_conceded"`
	WicketTakenBy    string `json:"wicket_taken_by"`
	WicketOf         string `json:"wicket_of"`
}

func (s *CricketMatchScoreServer) AddCricketPlayerScoreFunc(ctx *gin.Context) {
	var req addCricketTeamPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var position int64
	if req.Position == "" {
		position = 0
	} else {
		position, err = strconv.ParseInt(req.Position, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse position: %v", err)
			return
		}
	}

	var runsScored int64
	if req.RunsScored == "" {
		runsScored = 0
	} else {
		runsScored, err = strconv.ParseInt(req.RunsScored, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse runs scored: %v", err)
			return
		}
	}

	var ballsFaced int64
	if req.BallsFaced == "" {
		ballsFaced = 0
	} else {
		ballsFaced, err = strconv.ParseInt(req.BallsFaced, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse ball faced: %v", err)
			return
		}
	}

	var fours int64
	if req.Fours == "" {
		fours = 0
	} else {
		fours, err = strconv.ParseInt(req.Fours, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse fours: %v", err)
			return
		}
	}

	var sixes int64
	if req.Sixes == "" {
		runsScored = 0
	} else {
		sixes, err = strconv.ParseInt(req.Sixes, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse sixes: %v", err)
			return
		}
	}

	var wicketsTaken int64
	if req.WicketsTaken == "" {
		wicketsTaken = 0
	} else {
		wicketsTaken, err = strconv.ParseInt(req.WicketsTaken, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse wicket taken: %v", err)
			return
		}
	}

	var runsConceded int64
	if req.RunsConceded == "" {
		runsConceded = 0
	} else {
		runsConceded, err = strconv.ParseInt(req.RunsConceded, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse runs conceded: %v", err)
			return
		}
	}

	var wicketTakenBy int64
	if req.WicketTakenBy == "" {
		wicketTakenBy = 0
	} else {
		wicketTakenBy, err = strconv.ParseInt(req.WicketTakenBy, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse wicket taken by: %v", err)
			return
		}
	}

	var wicketOf int64
	if req.WicketOf == "" {
		runsScored = 0
	} else {
		wicketOf, err = strconv.ParseInt(req.WicketOf, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse wicket of: %v", err)
			return
		}
	}
	arg := db.AddCricketTeamPlayerScoreParams{
		MatchID:          req.MatchID,
		TournamentID:     req.TournamentID,
		TeamID:           req.TeamID,
		BattingOrBowling: req.BattingOrBowling,
		Position:         position,
		PlayerID:         req.PlayerID,
		RunsScored:       runsScored,
		BallsFaced:       ballsFaced,
		Fours:            fours,
		Sixes:            sixes,
		WicketsTaken:     wicketsTaken,
		OversBowled:      req.OversBowled,
		RunsConceded:     runsConceded,
		WicketTakenBy:    wicketTakenBy,
		WicketOf:         wicketOf,
	}

	response, err := s.store.AddCricketTeamPlayerScore(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to add the cricket player score: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *CricketMatchScoreServer) GetCricketMatchScoreFunc(ctx *gin.Context) {

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse match id: %v", err)
		return
	}

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse team id: %v", err)
		return
	}

	playerIdStr := ctx.Query("player_id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse player id: %v", err)
		return
	}

	arg := db.GetCricketPlayerScoreParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
		TeamID:       teamID,
		PlayerID:     playerID,
	}

	response, err := s.store.GetCricketPlayerScore(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to cricket player score : %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getCricketTeamPlayerScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *CricketMatchScoreServer) GetCricketPlayerScoreFunc(ctx *gin.Context) {

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse match id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to team id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketTeamPlayerScoreParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
		TeamID:       teamID,
	}

	response, err := s.store.GetCricketTeamPlayerScore(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to get cricket team player score : %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

func (s *CricketMatchScoreServer) UpdateCricketMatchScoreBattingFunc(ctx *gin.Context) {
	var req updateCricketMatchScoreBattingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var position int64
	if req.Position == "" {
		position = 0
	} else {
		position, err = strconv.ParseInt(req.Position, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse the position: %v", err)
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
			fmt.Errorf("Failed to parse runs scored %v", err)
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
			fmt.Errorf("Failed to parse balls faced: %v", err)
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
			fmt.Errorf("Failed to parse fours: %v", err)
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
			fmt.Errorf("Failed to parse sixes: %v", err)
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
			fmt.Errorf("Failed to parse wicket taken by: %v", err)
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
		fmt.Errorf("Failed to update cricket team score batting: %v", err)
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

func (s *CricketMatchScoreServer) UpdateCricketMatchScoreBowlingFunc(ctx *gin.Context) {
	var req updateCricketMatchScoreBowlingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var wicketsTaken int64
	if req.WicketsTaken == "" {
		wicketsTaken = 0
	} else {
		wicketsTaken, err = strconv.ParseInt(req.WicketsTaken, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse wicket taken: %v", err)
			return
		}
	}

	var runsConceded int64
	if req.RunsConceded == "" {
		runsConceded = 0
	} else {
		runsConceded, err = strconv.ParseInt(req.RunsConceded, 10, 64)
		if err != nil {
			fmt.Errorf("Failed to parse runs conceded: %v", err)
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
		fmt.Errorf("Failed to update cricket score bowling: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
