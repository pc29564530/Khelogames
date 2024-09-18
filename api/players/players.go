package players

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type newPlayerRequest struct {
	Positions string `json:"positions"`
	Sports    string `json:"sports"`
	Country   string `json:"country"`
	GameID    int64  `json:"game_id"`
}

func (s *PlayerServer) NewPlayerFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add player profile")
	var req newPlayerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	s.logger.Debug("Requested data: ", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	userProfile, err := s.store.GetProfile(ctx, authPayload.Username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to get the profile: %s", err))
		return
	}
	fullNameSlug := util.GenerateSlug(userProfile.FullName)
	shortName := util.GenerateShortName(userProfile.FullName)

	arg := db.NewPlayerParams{
		Username:   authPayload.Username,
		Slug:       fullNameSlug,
		ShortName:  shortName,
		MediaUrl:   userProfile.AvatarUrl,
		Positions:  req.Positions,
		Sports:     req.Sports,
		Country:    req.Country,
		PlayerName: userProfile.FullName,
		GameID:     req.GameID,
	}

	response, err := s.store.NewPlayer(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add player profile: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Added player profile: ", response)
	ctx.JSON(http.StatusAccepted, response)
}

func (s *PlayerServer) GetAllPlayerFunc(ctx *gin.Context) {

	response, err := s.store.GetAllPlayer(ctx)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	s.logger.Debug("Successfully get the player profile: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerServer) GetPlayerFunc(ctx *gin.Context) {

	playerIDStr := ctx.Query("id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse player id: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	response, err := s.store.GetPlayer(ctx, playerID)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	s.logger.Debug("Successfully get the player profile: ", response)

	ctx.JSON(http.StatusAccepted, response)
}

func (s *PlayerServer) GetPlayerSearchFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get player profile")
	playerName := ctx.Query("name")
	s.logger.Debug("Parse the player id: ", playerName)

	response, err := s.store.SearchPlayer(ctx, playerName)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	s.logger.Debug("Successfully get the player profile: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerServer) GetPlayerByCountry(ctx *gin.Context) {
	country := ctx.Query("country")
	response, err := s.store.GetPlayersCountry(ctx, country)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Successfully get all player profile: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerServer) UpdatePlayerMediaFunc(ctx *gin.Context) {
	playerIdStr := ctx.Query("id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse player id: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Parse the player id: ", playerID)

	playerMediaURL := ctx.Query("media_url")
	s.logger.Debug("Parse the player avatar ur: ", playerMediaURL)
	arg := db.UpdatePlayerMediaParams{
		MediaUrl: playerMediaURL,
		ID:       playerID,
	}

	response, err := s.store.UpdatePlayerMedia(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update player profile avatar: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Update the player profile Avatar: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerServer) UpdatePlayerPositionFunc(ctx *gin.Context) {
	playerIdStr := ctx.Query("id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse player id: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Parse the player id: ", playerID)

	playerPosition := ctx.Query("position")
	s.logger.Debug("Parse the player avatar ur: ", playerPosition)
	arg := db.UpdatePlayerPositionParams{
		Positions: playerPosition,
		ID:        playerID,
	}

	response, err := s.store.UpdatePlayerPosition(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update player profile avatar: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Update the player profile Avatar: ", response)

	ctx.JSON(http.StatusAccepted, response)
}
