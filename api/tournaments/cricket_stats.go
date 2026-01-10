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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentMostRuns(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most runs: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get cricket most runs",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentHighestRuns(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket highest runs: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get highest runs",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentMostSixes(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most sixes: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Failed to get cricket most sixes",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentMostFours(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most fours: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get most fours",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentMostFifties(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most fifties: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get most fifties",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentMostHundreds(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most hundreds: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get most hundreds",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentMostWickets(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket most wickets: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABSE_ERROR",
			"message": "Failed to get most wickets",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingEconomyRate(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket bowling economy rate: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get tournament bowling economy rate",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingAverage(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket bowling average: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get bowling average",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingStrikeRate(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket bowling strike rate: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get bowling strike rate",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentBowlingFiveWicketHaul(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket  five wicket haul: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get five wicket haul",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentBattingAverage(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket batting average: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get batting average",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	stats, err := s.store.GetCricketTournamentBattingStrikeRate(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("failed to get cricket batting strike: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get batting strike",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, stats)
	return
}
