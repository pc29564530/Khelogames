package players

import (
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *PlayerServer) GetPlayerCricketStatsByMatchType(ctx *gin.Context) {
	playerPublicIDString := ctx.Query("player_public_id")
	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parst to uuid: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}
	playerStats, err := s.store.GetCricketPlayerBowlingStats(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player bowling stats: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get player bowling stats",
		})
		return
	}
	var testStats *models.PlayerBowlingStats
	var odiStats *models.PlayerBowlingStats
	var t20Stats *models.PlayerBowlingStats
	for _, stats := range *playerStats {
		if stats.MatchType == "Test" {
			testStats = &stats
		} else if stats.MatchType == "ODI" {
			odiStats = &stats
		} else {
			t20Stats = &stats
		}
	}

	playerStatsByType := map[string]interface{}{
		"test": testStats,
		"odi":  odiStats,
		"t20":  t20Stats,
	}

	ctx.JSON(http.StatusAccepted, playerStatsByType)
}

func (s *PlayerServer) GetFootballPlayerStatsFunc(ctx *gin.Context) {
	var req struct {
		PlayerPublicID string `uri:"player_public_id"`
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

	playerPublicID, err := uuid.Parse(req.PlayerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	playersStats, err := s.store.GetFootballPlayerStats(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get football player stats: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get football player stats",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, playersStats)
}

func (s *PlayerServer) GetPlayerCricketStatsByMatchTypeFunc(ctx *gin.Context) {
	playerPublicIDString := ctx.Query("player_id")
	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parst to uuid: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	playerStats, err := s.store.GetPlayerCricketStatsByMatchType(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player stats for cricket: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get player cricket stats",
		})
		return
	}

	var testStats models.CricketPlayerStats
	var odiStats models.CricketPlayerStats
	var t20Stats models.CricketPlayerStats

	for _, item := range *playerStats {
		if item.MatchType == "Test" {
			testStats = item
		} else if item.MatchType == "ODI" {
			odiStats = item
		} else {
			t20Stats = item
		}
	}

	stats := map[string]interface{}{
		"Test": testStats,
		"ODI":  odiStats,
		"T20":  t20Stats,
	}

	ctx.JSON(http.StatusAccepted, stats)
}
