package clubs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchTeamRequest struct {
	ClubName string `json:"club_name"`
}

func (s *ClubServer) SearchTeamFunc(ctx *gin.Context) {
	var req searchTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	searchQuery := "%" + req.ClubName + "%"

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
