package players

import (
	"fmt"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type newPlayerRequest struct {
	Positions string `json:"positions"`
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
	userPlayer, err := s.store.GetProfile(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to get the profile: %s", err))
		return
	}
	fullNameSlug := util.GenerateSlug(userPlayer.FullName)
	shortName := util.GenerateShortName(userPlayer.FullName)

	arg := db.NewPlayerParams{
		UserPublicID: authPayload.PublicID,
		Name:         userPlayer.FullName,
		Slug:         fullNameSlug,
		ShortName:    shortName,
		MediaUrl:     userPlayer.AvatarUrl,
		Positions:    req.Positions,
		Country:      req.Country,
		GameID:       req.GameID,
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
	var req struct {
		PlayerPublicID string `uri:"public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	playerPublicID, err := uuid.Parse(req.PlayerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response, err := s.store.GetPlayerByPublicID(ctx, playerPublicID)
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

func (s *PlayerServer) GetPlayersBySportFunc(ctx *gin.Context) {
	gameIDString := ctx.Query("game_id")
	gameID, err := strconv.ParseInt(gameIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse game id: ", err)
		return
	}

	response, err := s.store.GetPlayersBySport(ctx, int32(gameID))
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Successfully get all player profile: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

// func (s *PlayerServer) UpdatePlayerMediaFunc(ctx *gin.Context) {
// 	playerIdStr := ctx.Query("id")
// 	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
// 	if err != nil {
// 		s.logger.Error("Failed to parse player id: ", err)
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}
// 	s.logger.Debug("Parse the player id: ", playerID)

// 	playerMediaURL := ctx.Query("media_url")
// 	s.logger.Debug("Parse the player avatar ur: ", playerMediaURL)
// 	arg := db.UpdatePlayerMediaParams{
// 		MediaUrl: playerMediaURL,
// 		ID:       playerID,
// 	}

// 	response, err := s.store.UpdatePlayerMedia(ctx, arg)
// 	if err != nil {
// 		s.logger.Error("Failed to update player profile avatar: ", err)
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}
// 	s.logger.Debug("Update the player profile Avatar: ", response)

// 	ctx.JSON(http.StatusAccepted, response)
// 	return
// }

// func (s *PlayerServer) UpdatePlayerPositionFunc(ctx *gin.Context) {
// 	playerPublicIdStr := ctx.Query("player_public_id")
// 	playerPublicID, err := uuid.Parse(playerPublicIdStr)
// 	if err != nil {
// 		s.logger.Error("Failed to parse player id: ", err)
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}
// 	s.logger.Debug("Parse the player id: ", playerPublicID)

// 	playerPosition := ctx.Query("position")
// 	s.logger.Debug("Parse the player avatar ur: ", playerPosition)
// 	arg := db.UpdatePlayerPositionParams{
// 		PublicID:  playerPublicID,
// 		Positions: playerPosition,
// 	}

// 	response, err := s.store.UpdatePlayerPosition(ctx, arg)
// 	if err != nil {
// 		s.logger.Error("Failed to update player profile avatar: ", err)
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}
// 	s.logger.Debug("Update the player profile Avatar: ", response)

// 	ctx.JSON(http.StatusAccepted, response)
// }
