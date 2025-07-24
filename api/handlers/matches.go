package handlers

import (
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) GetAllMatchesFunc(ctx *gin.Context) {

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

func (s *HandlersServer) GetMatchByMatchIDFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID uuid.UUID `uri:"match_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}

	match, err := s.store.GetMatchByMatchID(ctx, req.MatchPublicID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by match id: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, match)
	return
}
