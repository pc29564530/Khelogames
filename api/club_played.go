package api

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) getClubPlayedTournament(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse the tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse the club id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetClubPlayedTournamentParams{
		TournamentID: tournamentID,
		ClubID:       clubID,
	}

	response, err := server.store.GetClubPlayedTournament(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to get club played tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getClubPlayedTournaments(ctx *gin.Context) {
	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse the club id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	response, err := server.store.GetClubPlayedTournaments(ctx, clubID)
	if err != nil {
		server.logger.Error("Failed to get club played tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
