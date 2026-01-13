package football

import (
	db "khelogames/database"
	"net/http"

	errorhandler "khelogames/error_handler"

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
	s.logger.Info("Received request to add football statistics")
	var req addFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
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
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add football statistics",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully added football statistics")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

type getFootballStatisticsRequest struct {
	MatchPublicID string `json:"match_public_id"`
	TeamPublicID  string `json:"team_public_id"`
}

func (s *FootballServer) GetFootballStatisticsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get football statistics")
	var req getFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid team UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetFootballStatistics(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get the football statistics: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get football statistics",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully retrieved football statistics")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
