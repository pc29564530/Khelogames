package tournaments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *TournamentServer) GetCricketTournamentMostRunsFunc(ctx *gin.Context) {
	var req struct {
		TournamentID int64 `uri:"id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	stats, err := s.store.GetCricketTournamentMostRuns(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("failed to get cricket stats: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

func (s *TournamentServer) GetCricketTournamentHighestRunsFunc(ctx *gin.Context) {
	var req struct {
		TournamentID int64 `uri:"id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	stats, err := s.store.GetCricketTournamentHighestRuns(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("failed to get cricket highest runs: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentMostSixes

func (s *TournamentServer) GetCricketTournamentMostSixesFunc(ctx *gin.Context) {
	var req struct {
		TournamentID int64 `uri:"id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	stats, err := s.store.GetCricketTournamentMostSixes(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("failed to get cricket most sixes: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}
