package api

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createTournamentGroupRequest struct {
	GroupName     string `json:"group_name"`
	TournamentID  int64  `json:"tournament_id"`
	GroupStrength int64  `json:"group_strength"`
}

func (server *Server) createTournamentGroup(ctx *gin.Context) {
	var req createTournamentGroupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateTournamentGroupParams{
		GroupName:     req.GroupName,
		TournamentID:  req.TournamentID,
		GroupStrength: req.GroupStrength,
	}

	response, err := server.store.CreateTournamentGroup(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournamentGroup(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTournamentGroupParams{
		GroupID:      groupID,
		TournamentID: tournamentID,
	}

	response, err := server.store.GetTournamentGroup(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournamentGroups(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	response, err := server.store.GetTournamentGroups(ctx, tournamentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
