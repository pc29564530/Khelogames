package handlers

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerProfileServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewPlayerProfileServer(store *db.Store, logger *logger.Logger) *PlayerProfileServer {
	return &PlayerProfileServer{store: store, logger: logger}
}

type addPlayerProfileRequest struct {
	PlayerName            string `json:"player_name"`
	PlayerAvatarUrl       string `json:"player_avatar_url"`
	PlayerBio             string `json:"player_bio"`
	PlayerSport           string `json:"player_sport"`
	PlayerPlayingCategory string `json:"player_playing_category"`
	Nation                string `json:"nation"`
}

func (s *PlayerProfileServer) AddPlayerProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add player profile")
	var req addPlayerProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		return
	}
	s.logger.Debug("Requested data: %v", req)
	arg := db.AddPlayerProfileParams{
		PlayerName:            req.PlayerName,
		PlayerAvatarUrl:       req.PlayerAvatarUrl,
		PlayerBio:             req.PlayerBio,
		PlayerPlayingCategory: req.PlayerPlayingCategory,
		PlayerSport:           req.PlayerSport,
		Nation:                req.Nation,
	}

	response, err := s.store.AddPlayerProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add player profile: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Added player profile: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerProfileServer) GetPlayerProfileFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get player profile")
	playerIdStr := ctx.Query("player_id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse player id: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Parse the player id: %v", playerID)

	response, err := s.store.GetPlayerProfile(ctx, playerID)
	if err != nil {
		s.logger.Error("Failed to get player profile: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	s.logger.Debug("Successfully get the player profile: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerProfileServer) GetAllPlayerProfileFunc(ctx *gin.Context) {

	response, err := s.store.GetAllPlayerProfile(ctx)
	if err != nil {
		s.logger.Error("Failed to get player profile: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Successfully get all player profile: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerProfileServer) UpdatePlayerProfileAvatarFunc(ctx *gin.Context) {
	playerIdStr := ctx.Query("id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse player id: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Parse the player id: %v", playerID)

	playerAvatarUrl := ctx.Query("player_avatar_url")
	s.logger.Debug("Parse the player avatar ur: %v", playerAvatarUrl)
	arg := db.UpdatePlayerProfileAvatarParams{
		PlayerAvatarUrl: playerAvatarUrl,
		ID:              playerID,
	}

	response, err := s.store.UpdatePlayerProfileAvatar(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update player profile avatar: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	s.logger.Debug("Update the player profile Avatar: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}
