package badminton

import (
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *BadmintonServer) UpdateBadmintonScoreFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `json:"match_public_id"`
		TeamPublicID  string `json:"team_public_id"`
		SetNumber     int    `json:"set_number"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	score, err := s.store.UpdateBadmintonScore(ctx, matchPublicID, teamPublicID, req.SetNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update badminton score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    score,
	})
}

func (s *BadmintonServer) GetBadmintonScoreFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `json:"match_public_id"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	score, err := s.store.GetBadmintonScore(ctx, matchPublicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update badminton score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    score,
	})
}
