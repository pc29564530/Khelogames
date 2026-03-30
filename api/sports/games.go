package sports

import (
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getGamesRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *SportsServer) GetGameFunc(ctx *gin.Context) {
	var req getGamesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		s.logger.Error("Failed to bind URI: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Info("Received request to get game with ID: ", req.ID)

	response, err := s.store.GetGame(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get the games: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get game",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *SportsServer) GetGamesFunc(ctx *gin.Context) {
	response, err := s.store.GetGames(ctx)
	if err != nil {
		s.logger.Error("Failed to get the games: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get games",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}
