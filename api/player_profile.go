package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addPlayerProfileRequest struct {
	PlayerName            string `json:"player_name"`
	PlayerAvatarUrl       string `json:"player_avatar_url"`
	PlayerBio             string `json:"player_bio"`
	PlayerSport           string `json:"player_sport"`
	PlayerPlayingCategory string `json:"player_playing_category"`
	Nation                string `json:"nation"`
}

func (server *Server) addPlayerProfile(ctx *gin.Context) {
	var req addPlayerProfileRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Println("unable to get the data from frontend: ", err)
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

func (server *Server) getAllPlayerProfile(ctx *gin.Context) {

	response, err := server.store.GetAllPlayerProfile(ctx)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	fmt.Println("Line no 73: ", response)
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
