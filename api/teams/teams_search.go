package teams

import (
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchTeamRequest struct {
	Name string `json:"name"`
}

func (s *TeamsServer) SearchTeamFunc(ctx *gin.Context) {
	var req searchTeamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	searchQuery := "%" + req.Name + "%"

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := s.store.SearchTeam(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to search team",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
	return
}
