package api

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addGroupTeamRequest struct {
	GroupID      int64 `json:"group_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (server *Server) addGroupTeam(ctx *gin.Context) {
	var req addGroupTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddGroupTeamParams{
		GroupID:      req.GroupID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	response, err := server.store.AddGroupTeam(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to add group team: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTeamsByGroup(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to group id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTeamByGroupParams{
		TournamentID: tournamentID,
		GroupID:      groupID,
	}

	response, err := server.store.GetTeamByGroup(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to get team by group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
