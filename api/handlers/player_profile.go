package handlers

import (
	"fmt"
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
	var req addPlayerProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		return
	}

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
		fmt.Errorf("Failed to add player profile: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerProfileServer) GetPlayerProfileFunc(ctx *gin.Context) {
	playerIdStr := ctx.Query("player_id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse player id: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	response, err := s.store.GetPlayerProfile(ctx, playerID)
	if err != nil {
		fmt.Errorf("Failed to get player profile: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerProfileServer) GetAllPlayerProfileFunc(ctx *gin.Context) {

	response, err := s.store.GetAllPlayerProfile(ctx)
	if err != nil {
		fmt.Errorf("Failed to get player profile: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *PlayerProfileServer) UpdatePlayerProfileAvatarFunc(ctx *gin.Context) {
	playerIdStr := ctx.Query("id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse player id: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	playerAvatarUrl := ctx.Query("player_avatar_url")

	arg := db.UpdatePlayerProfileAvatarParams{
		PlayerAvatarUrl: playerAvatarUrl,
		ID:              playerID,
	}

	response, err := s.store.UpdatePlayerProfileAvatar(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update player profile avatar: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
