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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public ID: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersGoals(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get football tournament goals: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Invalid request format",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersYellowCard(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get football tournament yellow cards: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersRedCard(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get football tournament red cards: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get player red card",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, stats)
	return
}
