package tournaments

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TournamentServer) GetFootballTournamentPlayersGoalsFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public ID: ", err)
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersGoals(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get football tournament goals: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

func (s *TournamentServer) GetFootballTournamentPlayersYellowCardFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public id: ", err)
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersYellowCard(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get football tournament yellow cards: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

func (s *TournamentServer) GetFootballTournamentPlayersRedCardFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public id: ", err)
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersRedCard(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get football tournament red cards: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, stats)
	return
}
