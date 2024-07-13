package clubs

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *ClubServer) GetClubPlayedTournamentFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get club played tournament")

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse the tournament ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tournament ID"})
		return
	}
	s.logger.Debug("Parsed tournament ID: %d", tournamentID)

	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse the club ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid club ID"})
		return
	}
	s.logger.Debug("Parsed club ID: %d", clubID)

	arg := db.GetClubPlayedTournamentParams{
		TournamentID: tournamentID,
		ClubID:       clubID,
	}

	response, err := s.store.GetClubPlayedTournament(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get club played tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get club played tournament"})
		return
	}

	s.logger.Info("Successfully retrieved club played tournament")
	ctx.JSON(http.StatusOK, response)
}

func (s *ClubServer) GetClubPlayedTournamentsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get club played tournaments")

	clubIDStr := ctx.Query("club_id")
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse the club ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid club ID"})
		return
	}
	s.logger.Debug("Parsed club ID: %d", clubID)

	response, err := s.store.GetClubPlayedTournaments(ctx, clubID)
	if err != nil {
		s.logger.Error("Failed to get club played tournaments: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get club played tournaments"})
		return
	}

	s.logger.Info("Successfully retrieved club played tournaments")
	ctx.JSON(http.StatusOK, response)
	return
}
