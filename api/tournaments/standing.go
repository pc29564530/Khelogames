package tournaments

import (
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createTournamentStandingRequest struct {
	TournamentPublicID string `json:"tournament_public_id"`
	GroupID            int32  `json:"group_id"`
	TeamPublicID       string `json:"team_public_id"`
}

func (s *TournamentServer) CreateTournamentStandingFunc(ctx *gin.Context) {
	var req createTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	game := ctx.Param("sport")

	if game == "football" {
		var footballStanding models.FootballStanding
		footballStanding, err := s.store.CreateFootballStanding(ctx, tournamentPublicID, req.GroupID, teamPublicID)
		if err != nil {
			s.logger.Error("Failed to create football standing: ", err)
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusAccepted, footballStanding)
	} else if game == "cricket" {
		cricketStanding, err := s.store.CreateCricketStanding(ctx, tournamentPublicID, req.GroupID, teamPublicID)
		if err != nil {
			s.logger.Error("Failed to create cricket standing: ", err)
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusAccepted, cricketStanding)
	}
}
