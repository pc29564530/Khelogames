package teams

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchTeamRequest struct {
	Name string `json:"name"`
}

func (s *TeamsServer) SearchTeamFunc(ctx *gin.Context) {
	var req searchTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	searchQuery := "%" + req.Name + "%"

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := s.store.SearchTeam(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
