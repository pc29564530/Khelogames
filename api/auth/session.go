package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *AuthServer) DeleteSessionFunc(ctx *gin.Context) {
	var req struct {
		PublicID uuid.UUID `uri:"public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	err = s.store.DeleteSessions(ctx, req.PublicID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	ctx.JSON(http.StatusAccepted, "Deleted Session ")
	return
}
