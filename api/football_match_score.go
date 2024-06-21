package api

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type addFootballMatchScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
	GoalFor      int64 `json:"goal_for"`
	GoalAgainst  int64 `json:"goal_against"`
}

func (server *Server) addFootballMatchScore(ctx *gin.Context) {

	var req addFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind football match score: %v", err)
		return
	}

	// matchIdStr := ctx.Query("match_id")
	// matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 15: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	// tournamentIdStr := ctx.Query("tournament_id")
	// tournamentId, err := strconv.ParseInt(tournamentIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 24: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	// teamIdStr := ctx.Query("team_id")
	// teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 32: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	// goalScoreStr := ctx.Query("goal_score")
	// goalScore, err := strconv.ParseInt(goalScoreStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 40: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	arg := db.AddFootballMatchScoreParams{
		MatchID:      req.MatchID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
		GoalFor:      req.GoalFor,
		GoalAgainst:  req.GoalAgainst,
	}

	response, err := server.store.AddFootballMatchScore(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to add football match score: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (server *Server) getFootballMatchScore(ctx *gin.Context) {
	matchIdStr := ctx.Query("match_id")
	matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse match id: %v", err)
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse team id: %v", err)
		return
	}

	arg := db.GetFootballMatchScoreParams{
		MatchID: matchId,
		TeamID:  teamID,
	}

	response, err := server.store.GetFootballMatchScore(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to get football match score: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type updateFootballMatchScoreRequest struct {
	// TournamentID int64 `json:"tournament_id"`
	// MatchID      int64 `json:"match_id"`
	// TeamID       int64 `json:"team_id"`

	GoalFor     int64 `json:"goal_for"`
	GoalAgainst int64 `json:"goal_against"`
	MatchID     int64 `json:"match_id"`
	TeamID      int64 `json:"team_id"`
}

func (server *Server) updateFootballMatchScore(ctx *gin.Context) {

	var req updateFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind update football match score: %v", err)
		return
	}

	// tournamentIdStr := ctx.Query("tournament_id")
	// matchIdStr := ctx.Query("match_id")
	// teamIdStr := ctx.Query("team_id")
	// goalScoreStr := ctx.Query("goal_score")
	// matchID, err := strconv.ParseInt(matchIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	// teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	// goalScore, err := strconv.ParseInt(goalScoreStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	// tournamentId, err := strconv.ParseInt(tournamentIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	arg := db.UpdateFootballMatchScoreParams{
		GoalFor:     req.GoalFor,
		GoalAgainst: req.GoalAgainst,
		MatchID:     req.MatchID,
		TeamID:      req.TeamID,
	}

	response, err := server.store.UpdateFootballMatchScore(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update football match score: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type addFootballteamPlayerScoreRequest struct {
	MatchID       int64     `json:"match_id"`
	TeamID        int64     `json:"team_id"`
	PlayerID      int64     `json:"player_id"`
	TournamentID  int64     `json:"tournament_id"`
	GoalScoreTime time.Time `json:"goal_score_time"`
}

func (server *Server) addFootballGoalByPlayer(ctx *gin.Context) {

	var req addFootballteamPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind add football goal: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddFootballGoalByPlayerParams{
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		PlayerID:      req.PlayerID,
		TournamentID:  req.TournamentID,
		GoalScoreTime: req.GoalScoreTime,
	}

	response, err := server.store.AddFootballGoalByPlayer(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to add football goal by player : %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}
