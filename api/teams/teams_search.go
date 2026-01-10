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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Failed to bind the request",
		})
		return
	}
	searchQuery := "%" + req.Name + "%"

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := s.store.SearchTeam(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to search team",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
