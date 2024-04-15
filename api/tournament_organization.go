package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type createTournamentOrganizationRequest struct {
	TournamentID    int64     `json:"tournament_id"`
	PlayerCount     int64     `json:"player_count"`
	TeamCount       int64     `json:"team_count"`
	GroupCount      int64     `json:"group_count"`
	AdvancedTeam    int64     `json:"advanced_team"`
	TournamentStart time.Time `json:"tournament_start"`
}

func (server *Server) createTournamentOrganization(ctx *gin.Context) {

	var req createTournamentOrganizationRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Println("Unable to get the struct name: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateTournamentOrganizationParams{
		TournamentID:    req.TournamentID,
		TournamentStart: req.TournamentStart,
		PlayerCount:     req.PlayerCount,
		TeamCount:       req.TeamCount,
		GroupCount:      req.GroupCount,
		AdvancedTeam:    req.AdvancedTeam,
	}

	response, err := server.store.CreateTournamentOrganization(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournamentOrganization(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil || tournamentIDStr == " " {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := server.store.GetTournamentOrganization(ctx, tournamentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
