package players

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchPlayerRequest struct {
	Name string `json:"name"`
}

func (s *PlayerServer) SearchProfileFunc(ctx *gin.Context) {
	var req searchPlayerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	searchQuery := "%" + req.Name + "%"

	response, err := s.store.SearchPlayer(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
