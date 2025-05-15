package tournaments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *TournamentServer) GetFootballTournamentPlayersGoalsFunc(ctx *gin.Context) {
	var req struct {
		TournamentID int64 `uri:"id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersGoals(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("Failed to get football tournament goals: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

func (s *TournamentServer) GetFootballTournamentPlayersYellowCardFunc(ctx *gin.Context) {
	var req struct {
		TournamentID int64 `uri:"id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersYellowCard(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("Failed to get football tournament yellow cards: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

func (s *TournamentServer) GetFootballTournamentPlayersRedCardFunc(ctx *gin.Context) {
	var req struct {
		TournamentID int64 `uri:"id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	stats, err := s.store.GetFootballTournamentPlayersRedCard(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("Failed to get football tournament red cards: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, stats)
	return
}
