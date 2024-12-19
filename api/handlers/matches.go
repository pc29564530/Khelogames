package handlers

import (
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getAllMatchesReq struct {
	StartTimestamp string `json:"start_timestamp"`
}

func (s *HandlersServer) GetAllMatchesFunc(ctx *gin.Context) {

	// var req getAllMatchesReq
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	s.logger.Error("Failed to bind: ", err)
	// 	return
	// }

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}
	startDateString := ctx.Query("start_timestamp")
	startDate, err := util.ConvertTimeStamp(startDateString)
	if err != nil {
		s.logger.Error("Failed to convert to second: ", err)
	}
	response, err := s.store.GetAllMatches(ctx, int32(startDate), game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
