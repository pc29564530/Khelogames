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
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	err = s.store.DeleteSessions(ctx, publicID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	ctx.JSON(http.StatusAccepted, "Deleted Session ")
	return
}
