package tournaments

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TournamentServer) GetCricketTournamentMostRunsFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentMostRuns(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most runs: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

func (s *TournamentServer) GetCricketTournamentHighestRunsFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentHighestRuns(ctx, tournamentPublicID)
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
		TournamentPublicID string `uri:"tournament_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentMostSixes(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most sixes: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentMostFours
func (s *TournamentServer) GetCricketTournamentMostFoursFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentMostFours(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most fours: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentMostFifties
func (s *TournamentServer) GetCricketTournamentMostFiftiesFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentMostFifties(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most fifties: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentMostHundreds
func (s *TournamentServer) GetCricketTournamentMostHundredsFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentMostHundreds(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most hundreds: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// Bowling Function
// GetCricketTournamentMostWickets
func (s *TournamentServer) GetCricketTournamentMostWicketsFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentMostWickets(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most wickets: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentBowlwingEconomyRate
func (s *TournamentServer) GetCricketTournamentBowlingEconomyRateFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingEconomyRate(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket bowling economy rate: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentBowlingAverage
func (s *TournamentServer) GetCricketTournamentBowlingAverageFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingAverage(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket bowling average: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentBowlwingEconomyRate
func (s *TournamentServer) GetCricketTournamentBowlingStrikeRateFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingStrikeRate(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket bowling strike rate: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentBowlwingEconomyRate
func (s *TournamentServer) GetCricketTournamentBowlingFiveWicketHaulFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingFiveWicketHaul(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket  five wicket haul: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentBattingAverage
func (s *TournamentServer) GetCricketTournamentBattingAverageFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentBattingAverage(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket batting average: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}

// GetCricketTournamentBattingAverage
func (s *TournamentServer) GetCricketTournamentBattingStrikeFunc(ctx *gin.Context) {
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
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	stats, err := s.store.GetCricketTournamentBattingStrikeRate(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket batting strike: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}
