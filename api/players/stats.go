package players

import (
	"khelogames/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *PlayerServer) GetPlayerCricketStatsByMatchType(ctx *gin.Context) {
	playerIDString := ctx.Query("player_id")
	playerID, err := strconv.ParseInt(playerIDString, 10, 64)
	playerStats, err := s.store.GetCricketPlayerBowlingStats(ctx, playerID)
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
		MatchID int64 `uri:"match_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	playersStats, err := s.store.GetFootballPlayerStats(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to get football player stats: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, playersStats)
}

func (s *PlayerServer) GetPlayerCricketStatsByMatchTypeFunc(ctx *gin.Context) {
	playerIDString := ctx.Query("player_id")
	playerID, err := strconv.ParseInt(playerIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parst to int: ", err)
		return
	}

	playerStats, err := s.store.GetPlayerCricketStatsByMatchType(ctx, playerID)
	if err != nil {
		s.logger.Error("Failed to get the player stats for cricket: ", err)
		return
	}

	var testStats models.PlayerCricketStats
	var odiStats models.PlayerCricketStats
	var t20Stats models.PlayerCricketStats

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
