package football

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *FootballServer) GetFootballTopPerformerFunc(ctx *gin.Context) {
	res, err := s.store.GetFootballTopPerformer(ctx)
	if err != nil {
		s.logger.Error("Failed to get top performer", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Unable to get top performer",
			},
			"request_id": ctx.GetString("request_id"),
		})
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    res,
	})
}
