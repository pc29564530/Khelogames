package handlers

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClubTournamentServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewClubTournamentServer(store *db.Store, logger *logger.Logger) *ClubTournamentServer {
	return &ClubTournamentServer{store, logger}
}

func (s *ClubTournamentServer) GetClubPlayedTournamentFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse the tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse the club id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetClubPlayedTournamentParams{
		TournamentID: tournamentID,
		ClubID:       clubID,
	}

	response, err := s.store.GetClubPlayedTournament(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to get club played tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubTournamentServer) GetClubPlayedTournamentsFunc(ctx *gin.Context) {
	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse the club id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	response, err := s.store.GetClubPlayedTournaments(ctx, clubID)
	if err != nil {
		fmt.Errorf("Failed to get club played tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
