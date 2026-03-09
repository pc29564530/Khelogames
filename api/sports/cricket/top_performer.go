package cricket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *CricketServer) GetCricketTopPerformerFunc(ctx *gin.Context) {
	batting, err := s.store.GetCricketTopBattingPerformer(ctx)
	if err != nil {
		s.logger.Error("Failed to get top performer", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Unable to get batting top performer",
			},
			"request_id": ctx.GetString("request_id"),
		})
	}
	bowling, err := s.store.GetCricketTopBowlingPerformer(ctx)
	if err != nil {
		s.logger.Error("Failed to get top performer", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Unable to get bowling top performer",
			},
			"request_id": ctx.GetString("request_id"),
		})
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": gin.H{
			"batting": batting,
			"bowling": bowling,
		},
	})
}
