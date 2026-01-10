package football

import (
	db "khelogames/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addFootballStatisticsRequest struct {
	MatchID         int32 `json:"match_id"`
	TeamID          int32 `json:"team_id"`
	ShotsOnTarget   int32 `json:"shots_on_target"`
	TotalShots      int32 `json:"total_shots"`
	CornerKicks     int32 `json:"corner_kicks"`
	Fouls           int32 `json:"fouls"`
	GoalkeeperSaves int32 `json:"goalkeeper_saves"`
	FreeKicks       int32 `json:"free_kicks"`
	YellowCards     int32 `json:"yellow_cards"`
	RedCards        int32 `json:"red_cards"`
}

func (s *FootballServer) AddFootballStatisticsFunc(ctx *gin.Context) {
	var req addFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	arg := db.CreateFootballStatisticsParams{
		MatchID:         int32(req.MatchID),
		TeamID:          int32(req.TeamID),
		ShotsOnTarget:   req.ShotsOnTarget,
		TotalShots:      req.TotalShots,
		CornerKicks:     req.CornerKicks,
		Fouls:           req.Fouls,
		GoalkeeperSaves: req.GoalkeeperSaves,
		FreeKicks:       req.FreeKicks,
		YellowCards:     req.YellowCards,
		RedCards:        req.RedCards,
	}

	response, err := s.store.CreateFootballStatistics(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the football statistics: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to add football statistics",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response)

}

type getFootballStatisticsRequest struct {
	MatchPublicID string `json:"match_public_id"`
	TeamPublicID  string `json:"team_public_id"`
}

func (s *FootballServer) GetFootballStatisticsFunc(ctx *gin.Context) {
	var req getFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetFootballStatistics(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get the football statistics: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get football statistics",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
