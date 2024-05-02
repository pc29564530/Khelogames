package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (server *Server) addCricketMatchScore(ctx *gin.Context) {

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
	fmt.Println("Line no 42: ", arg)
	response, err := server.store.CreateCricketMatchScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (server *Server) getCricketMatchScore(ctx *gin.Context) {
	matchIdStr := ctx.Query("match_id")
	matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	arg := db.GetCricketMatchScoreParams{
		MatchID: matchId,
		TeamID:  teamID,
	}
	fmt.Println("Arg: ", arg)

	response, err := server.store.GetCricketMatchScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	fmt.Println("Score: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return

}

type updateCricketMatchScoreRequest struct {
	Score   int64 `json:"score"`
	Wickets int64 `json:"wickets"`
	Extras  int64 `json:"extras"`
	Innings int64 `json:"innings"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (server *Server) updateCricketMatchScore(ctx *gin.Context) {

	var req updateCricketMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("unable to get the correct data or order ", err)
		return
	}

	arg := db.UpdateCricketMatchScoreParams{
		Score:   req.Score,
		Wickets: req.Wickets,
		Extras:  req.Extras,
		Innings: req.Innings,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := server.store.UpdateCricketMatchScore(ctx, arg)
	if err != nil {
		fmt.Errorf("unable to get the response ")
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
	BattingOrBowling string `json:"batting_or_bowling"`
	Position         int64  `json:"position"`
	PlayerID         int64  `json:"player_id"`
	RunsScored       int64  `json:"runs_scored",validate:"min=0"`
	BallsFaced       int64  `json:"balls_faced",validate:"min=0"`
	Fours            int64  `json:"fours", validate:"min=0"`
	Sixes            int64  `json:"sixes", validate:"min=0"`
	WicketsTaken     int64  `json:"wickets_taken", validate:"min=0"`
	OversBowled      string `json:"overs_bowled", validate:"min=0"`
	RunsConceded     int64  `json:"runs_conceded", validate:"min=0"`
	WicketTakenBy    int64  `json:"wicket_taken_by"`
	WicketOf         int64  `json:"wicket_of"`
}

func (server *Server) addCricketTeamPlayerScore(ctx *gin.Context) {
	var req addCricketTeamPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.AddCricketTeamPlayerScoreParams{
		MatchID:          req.MatchID,
		TournamentID:     req.TournamentID,
		TeamID:           req.TeamID,
		BattingOrBowling: req.BattingOrBowling,
		Position:         req.Position,
		PlayerID:         req.PlayerID,
		RunsScored:       req.RunsScored,
		BallsFaced:       req.BallsFaced,
		Fours:            req.Fours,
		Sixes:            req.Sixes,
		WicketsTaken:     req.WicketsTaken,
		OversBowled:      req.OversBowled,
		RunsConceded:     req.RunsConceded,
		WicketTakenBy:    req.WicketTakenBy,
		WicketOf:         req.WicketOf,
	}

	response, err := server.store.AddCricketTeamPlayerScore(ctx, arg)
	if err != nil {
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

func (server *Server) getCricketTeamPlayerScore(ctx *gin.Context) {
	// var req getCricketTeamPlayerScoreRequest
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	fmt.Println("error: ", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		fmt.Println("error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Println("error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		fmt.Println("error: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketTeamPlayerScoreParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
		TeamID:       teamID,
	}

	response, err := server.store.GetCricketTeamPlayerScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchScoreBattingRequest struct {
	Position      int64 `json:"position"`
	RunsScored    int64 `json:"runs_scored"`
	BallsFaced    int64 `json:"balls_faced"`
	Fours         int64 `json:"fours"`
	Sixes         int64 `json:"sixes"`
	WicketTakenBy int64 `json:"wicket_taken_by"`
	TournamentID  int64 `json:"tournament_id"`
	MatchID       int64 `json:"match_id"`
	TeamID        int64 `json:"team_id"`
	PlayerID      int64 `json:"player_id"`
}

func (server *Server) updateCricketMatchScoreBatting(ctx *gin.Context) {
	var req updateCricketMatchScoreBattingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	arg := db.UpdateCricketTeamPlayerScoreBattingParams{
		Position:      req.Position,
		RunsScored:    req.RunsScored,
		BallsFaced:    req.BallsFaced,
		Fours:         req.Fours,
		Sixes:         req.Sixes,
		WicketTakenBy: req.WicketTakenBy,
		TournamentID:  req.TournamentID,
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		PlayerID:      req.PlayerID,
	}

	response, err := server.store.UpdateCricketTeamPlayerScoreBatting(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchScoreBowlingRequest struct {
	OversBowled  string `json:"overs_bowled"`
	RunsConceded int64  `json:"runs_conceded"`
	WicketsTaken int64  `json:"wickets_taken"`
	TournamentID int64  `json:"tournament_id"`
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
	PlayerID     int64  `json:"player_id"`
}

func (server *Server) updateCricketMatchScoreBowling(ctx *gin.Context) {
	var req updateCricketMatchScoreBowlingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	arg := db.UpdateCricketTeamPlayerScoreBowlingParams{
		OversBowled:  req.OversBowled,
		RunsConceded: req.RunsConceded,
		WicketsTaken: req.WicketsTaken,
		TournamentID: req.TournamentID,
		MatchID:      req.MatchID,
		TeamID:       req.TeamID,
		PlayerID:     req.PlayerID,
	}

	response, err := server.store.UpdateCricketTeamPlayerScoreBowling(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
