package badminton

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *BadmintonServer) GetBadmintonTopPerformerFunc(ctx *gin.Context) {

	topPerformer, err := s.store.GetBadmintonTopPerformer(ctx)
	if err != nil {
		s.logger.Error("Failed to get badminton top performer", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Unable to get badminton top performer",
			},
			"request_id": ctx.GetString("request_id"),
		})
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    topPerformer,
	})
}
