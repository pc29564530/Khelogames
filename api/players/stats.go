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
	playerStats, err := s.store.GetCricketPlayerBowlingStats(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player bowling stats: ", err)
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
		MatchPublicID string `uri:"match_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	playersStats, err := s.store.GetFootballPlayerStats(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get football player stats: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, playersStats)
}

func (s *PlayerServer) GetPlayerCricketStatsByMatchTypeFunc(ctx *gin.Context) {
	playerPublicIDString := ctx.Query("player_id")
	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parst to uuid: ", err)
		return
	}

	playerStats, err := s.store.GetPlayerCricketStatsByMatchType(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player stats for cricket: ", err)
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
