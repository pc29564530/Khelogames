package tournaments

import (
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTournamentStandingRequest struct {
	TournamentID int64 `json:"tournament_id"`
	GroupID      int64 `json:"group_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) CreateTournamentStandingFunc(ctx *gin.Context) {
	var req createTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	game := ctx.Param("sport")

	if game == "football" {
		var footballStanding models.FootballStanding
		footballStanding, err := s.store.CreateFootballStanding(ctx, req.TournamentID, req.GroupID, req.TeamID)
		if err != nil {
			s.logger.Error("Failed to create football standing: ", err)
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusAccepted, footballStanding)
	} else if game == "cricket" {
		cricketStanding, err := s.store.CreateCricketStanding(ctx, req.TournamentID, req.GroupID, req.TeamID)
		if err != nil {
			s.logger.Error("Failed to create cricket standing: ", err)
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusAccepted, cricketStanding)
	}
}
