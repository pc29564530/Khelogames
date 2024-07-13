package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteSessionRequest struct {
	Username string `uri:"username"`
}

func (s *AuthServer) DeleteSessionFunc(ctx *gin.Context) {
	var req deleteSessionRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	err = s.store.DeleteSessions(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	ctx.JSON(http.StatusAccepted, "Deleted Session ")
	return
}
