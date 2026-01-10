package sports

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getGamesRequest struct {
	ID int64 `uri:"ID"`
}

func (s *SportsServer) GetGameFunc(ctx *gin.Context) {
	var req getGamesRequest
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

	response, err := s.store.GetGame(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get the games: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get game",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}

func (s *SportsServer) GetGamesFunc(ctx *gin.Context) {

	response, err := s.store.GetGames(ctx)
	if err != nil {
		s.logger.Error("Failed to get the games: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get games",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
