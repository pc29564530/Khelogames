package players

import (
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *PlayerServer) AddCricketPlayerBattingStatsFunc(ctx *gin.Context) {
	var req struct {
		PlayerID     int32  `json:"player_id"`
		MatchType    string `json:"match_type"`
		TotalMatches int    `json:"total_matches"`
		TotalInnings int    `json:"total_innings"`
		Runs         int    `json:"runs"`
		Balls        int    `json:"balls"`
		Sixes        int    `json:"sixes"`
		Fours        int    `json:"fours"`
		Fifties      int    `json:"fifties"`
		Hundreds     int    `json:"hundreds"`
		BestScore    int    `json:"best_score"`
		Average      string `json:"average"`
		StrikeRate   string `json:"strike_rate"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	playerStats, err := s.store.AddPlayerBattingStats(ctx, req.PlayerID, req.MatchType, req.TotalMatches, req.TotalInnings, req.Runs, req.Balls, req.Fours, req.Sixes, req.Fifties, req.Hundreds, req.BestScore, req.Average, req.StrikeRate)
	if err != nil {
		s.logger.Error("Failed to get the player batting stats: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, playerStats)
}

func (s *PlayerServer) GetCricketPlayerBattingStatsFunc(ctx *gin.Context) {
	var req struct {
		PlayerID int32 `uri:"player_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	playerStats, err := s.store.GetCricketPlayerBattingStats(ctx, int64(req.PlayerID))
	if err != nil {
		s.logger.Error("Failed to get the player batting stats: ", err)
		return
	}
	var testStats *models.PlayerBattingStats
	var odiStats *models.PlayerBattingStats
	var t20Stats *models.PlayerBattingStats
	for _, stats := range *playerStats {
		if stats.MatchType == "test" {
			testStats = &stats
		} else if stats.MatchType == "odi" {
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

func (s *PlayerServer) AddCricketPlayerBowlingStatsFunc(ctx *gin.Context) {
	var req struct {
		PlayerID    int32  `json:"player_id"`
		MatchType   string `json:"match_type"`
		Matches     int    `json:"matches"`
		Innings     int    `json:"innings"`
		Wickets     int    `json:"wickets"`
		Runs        int    `json:"runs"`
		Ball        int    `json:"balls"`
		Average     string `json:"average"`
		StrikeRate  string `json:"strike_rate"`
		EconomyRate string `json:"economy_rate"`
		FourWickets int    `json:"four_wickets"`
		FiveWickets int    `json:"five_wickets"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	playerStats, err := s.store.AddPlayerBowlingStats(ctx, req.PlayerID, req.MatchType, req.Matches, req.Innings, req.Wickets, req.Runs, req.Ball, req.Average, req.StrikeRate, req.EconomyRate, req.FourWickets, req.FiveWickets)
	if err != nil {
		s.logger.Error("Failed to get the player bowling stats: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, playerStats)
}

func (s *PlayerServer) GetCricketPlayerBowlingStatsFunc(ctx *gin.Context) {
	var req struct {
		PlayerID int32 `uri:"player_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	playerStats, err := s.store.GetCricketPlayerBowlingStats(ctx, int64(req.PlayerID))
	if err != nil {
		s.logger.Error("Failed to get the player bowling stats: ", err)
		return
	}
	var testStats *models.PlayerBowlingStats
	var odiStats *models.PlayerBowlingStats
	var t20Stats *models.PlayerBowlingStats
	for _, stats := range *playerStats {
		if stats.MatchType == "test" {
			testStats = &stats
		} else if stats.MatchType == "odi" {
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

func (s *PlayerServer) AddOrUpdateFootballPlayerStatsFunc(ctx *gin.Context) {
	var req struct {
		MatchID int64 `uri:"match_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	playersStats, err := s.store.AddORUpdateFootballPlayerStats(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to add or update football player stats: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, playersStats)
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
