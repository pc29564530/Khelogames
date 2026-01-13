package players

import (
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchPlayerRequest struct {
	Name string `json:"name"`
}

func (s *PlayerServer) SearchPlayerFunc(ctx *gin.Context) {
	var req searchPlayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	searchQuery := "%" + req.Name + "%"

	response, err := s.store.SearchPlayer(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to search player profile",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
