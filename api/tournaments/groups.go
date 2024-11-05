package tournaments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *TournamentServer) GetGroupsFunc(ctx *gin.Context) {

	response, err := s.store.GetGroups(ctx)
	if err != nil {
		s.logger.Error("Failed to get groups: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully get group: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
