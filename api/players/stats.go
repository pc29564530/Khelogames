package players

import (
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *PlayerServer) GetPlayerCricketStatsByMatchType(ctx *gin.Context) {
	playerPublicIDString := ctx.Query("player_public_id")
	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parst to uuid: ", err)
		fieldErrors := map[string]string{"player_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	playerStats, err := s.store.GetCricketPlayerBowlingStats(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player bowling stats: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player bowling stats",
			},
			"request_id": ctx.GetString("request_id"),
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
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	playerPublicID, err := uuid.Parse(req.PlayerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"player_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	playersStats, err := s.store.GetFootballPlayerStats(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get football player stats: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get football player stats",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    playersStats,
	})
}

func (s *PlayerServer) GetPlayerCricketStatsByMatchTypeFunc(ctx *gin.Context) {
	playerPublicIDString := ctx.Query("player_id")
	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parst to uuid: ", err)
		fieldErrors := map[string]string{"Player_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	playerStats, err := s.store.GetPlayerCricketStatsByMatchType(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player stats for cricket: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player cricket stats",
			},
			"request_id": ctx.GetString("request_id"),
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

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    stats,
	})
}
