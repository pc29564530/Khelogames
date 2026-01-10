package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *AuthServer) DeleteSessionFunc(ctx *gin.Context) {
	var req struct {
		PublicID string `uri:"public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind the delete session request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}
	err = s.store.DeleteSessions(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to delete session: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to delete session",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, "Deleted Session ")
	return
}
