package tournaments

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addTournamentTeamRequest struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) AddTournamentTeamFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add a team")
	var req addTournamentTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.NewTournamentTeamParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	response, err := s.store.NewTournamentTeam(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add team: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully added team: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTournamentTeamsRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (s *TournamentServer) GetTournamentTeamsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get teams for a tournament")
	var req getTournamentTeamsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := s.store.GetTournamentTeams(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("Failed to get teams: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not found"})
		return
	}
	s.logger.Info("Successfully retrieved teams: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}
