package auth

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteSessionRequest struct {
	Username string `uri:"username"`
}

type DeleteSessionServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewSessionServer(store *db.Store, logger *logger.Logger) *DeleteSessionServer {
	return &DeleteSessionServer{store: store, logger: logger}
}

func (s *DeleteSessionServer) DeleteSessionFunc(ctx *gin.Context) {
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
