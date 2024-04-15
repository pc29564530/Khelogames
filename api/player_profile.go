package api

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) addPlayerProfile(ctx *gin.Context) {
	playerName := ctx.Query("player_name")
	playerAvatarURL := ctx.Query("player_avatar_url")
	playerBio := ctx.Query("player_bio")
	playerPlayingCategory := ctx.Query("player_playing_category")
	playerSport := ctx.Query("player_sport")
	nation := ctx.Query("nation")

	arg := db.AddPlayerProfileParams{
		PlayerName:            playerName,
		PlayerAvatarUrl:       playerAvatarURL,
		PlayerBio:             playerBio,
		PlayerPlayingCategory: playerPlayingCategory,
		PlayerSport:           playerSport,
		Nation:                nation,
	}

	response, err := server.store.AddPlayerProfile(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getPlayerProfile(ctx *gin.Context) {
	playerIdStr := ctx.Query("player_id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	response, err := server.store.GetPlayerProfile(ctx, playerID)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) updatePlayerProfileAvatar(ctx *gin.Context) {
	playerIdStr := ctx.Query("id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	playerAvatarUrl := ctx.Query("player_avatar_url")

	arg := db.UpdatePlayerProfileAvatarParams{
		PlayerAvatarUrl: playerAvatarUrl,
		ID:              playerID,
	}

	response, err := server.store.UpdatePlayerProfileAvatar(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
