package api

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) addFootballMatchScore(ctx *gin.Context) {
	matchIdStr := ctx.Query("match_id")
	matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	tournamentIdStr := ctx.Query("tournament_id")
	tournamentId, err := strconv.ParseInt(tournamentIdStr, 10, 64)
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

	goalScoreStr := ctx.Query("goal_score")
	goalScore, err := strconv.ParseInt(goalScoreStr, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	arg := db.AddFootballMatchScoreParams{
		MatchID:      matchId,
		TournamentID: tournamentId,
		TeamID:       teamID,
		GoalScore:    goalScore,
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
