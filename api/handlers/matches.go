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
	response, err := s.store.ListMatches(ctx, int32(startDate), game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *HandlersServer) GetLiveMatchesFunc(ctx *gin.Context) {

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}
	response, err := s.store.GetLiveMatches(ctx, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *HandlersServer) GetMatchByMatchIDFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}

	match, err := s.store.GetMatchByPublicId(ctx, matchPublicID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by match id: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, match)
	return
}
