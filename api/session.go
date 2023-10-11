package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteSessionRequest struct {
	Username string `uri:"username"`
}

func (server *Server) deleteSession(ctx *gin.Context) {
	var req deleteSessionRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteSessions(ctx, req.Username)
	fmt.Println("line 22", err)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, "Deleted Session ")
	return
}
