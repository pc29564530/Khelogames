package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addFootballMatchScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
	GoalScore    int64 `json:"goal_score"`
}

func (server *Server) addFootballMatchScore(ctx *gin.Context) {

	var req addFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
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
		GoalScore:    req.GoalScore,
	}

	response, err := server.store.AddFootballMatchScore(ctx, arg)
	if err != nil {
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
		ctx.Error(err)
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	arg := db.GetFootballMatchScoreParams{
		MatchID: matchId,
		TeamID:  teamID,
	}

	response, err := server.store.GetFootballMatchScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (server *Server) updateFootballMatchScore(ctx *gin.Context) {
	tournamentIdStr := ctx.Query("tournament_id")
	matchIdStr := ctx.Query("match_id")
	teamIdStr := ctx.Query("team_id")
	goalScoreStr := ctx.Query("goal_score")
	matchID, err := strconv.ParseInt(matchIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("unable to parse the match Id")
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("unable to parse the match Id")
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	goalScore, err := strconv.ParseInt(goalScoreStr, 10, 64)
	if err != nil {
		fmt.Errorf("unable to parse the match Id")
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	tournamentId, err := strconv.ParseInt(tournamentIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("unable to parse the match Id")
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.UpdateFootballMatchScoreParams{
		TournamentID: tournamentId,
		MatchID:      matchID,
		TeamID:       teamID,
		GoalScore:    goalScore,
	}

	response, err := server.store.UpdateFootballMatchScore(ctx, arg)
	if err != nil {
		fmt.Errorf("unable to get the response ")
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
