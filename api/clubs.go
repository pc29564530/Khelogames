package api

import (
	"encoding/base64"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type createClubRequest struct {
	ClubName  string `json:"club_name"`
	AvatarURL string `json:"avatar_url"`
	Sport     string `json:"sport"`
}

func (server *Server) createClub(ctx *gin.Context) {
	var req createClubRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind create club request: %v", err)
		return
	}

	var path string
	if req.AvatarURL != "" {
		b64Data := req.AvatarURL[strings.IndexByte(req.AvatarURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			server.logger.Error("Failed to decode string: %v", err)
			return
		}

		path, err = saveImageToFile(data, "image")
		if err != nil {
			server.logger.Error("Failed to create file: %v", err)
			return
		}
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateClubParams{
		ClubName:  req.ClubName,
		AvatarUrl: path,
		Sport:     req.Sport,
		Owner:     authPayload.Username,
	}

	response, err := server.store.CreateClub(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to create club: %v", err)
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getClubs(ctx *gin.Context) {

	response, err := server.store.GetClubs(ctx)
	if err != nil {
		server.logger.Error("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getClubRequest struct {
	ID int64 `uri:"id"`
}

func (server *Server) getClub(ctx *gin.Context) {
	var req getClubRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response, err := server.store.GetClub(ctx, req.ID)
	if err != nil {
		server.logger.Error("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url"`
	ClubName  string `json:"club_name"`
}

func (server *Server) updateClubAvatarUrl(ctx *gin.Context) {
	var req updateAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateAvatarUrlParams{
		AvatarUrl: req.AvatarUrl,
		ClubName:  req.ClubName,
	}

	response, err := server.store.UpdateAvatarUrl(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update avatar url: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

// type updateClubNameRequest struct {
// 	AvatarUrl string `json:"avatar_url"`
// 	ClubName  string `json:"club_name"`
// }

// func (server *Server) updateClubName(ctx *gin.Context) {
// 	var req updateClubNameRequest
// 	err := ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
// 	arg := db.UpdateClubNameParams{
// 		ClubName: req.ClubName,
// 	}

// 	response, err := server.store.UpdateClubName(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return
// }

type updateClubSport struct {
	ClubName string `json:"club_name"`
	Sport    string `json:"sport"`
}

func (server *Server) updateClubSport(ctx *gin.Context) {
	var req updateClubSport
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateClubSportParams{
		Sport:    req.Sport,
		ClubName: req.ClubName,
	}

	response, err := server.store.UpdateClubSport(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update club sport: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type searchTeamRequest struct {
	ClubName string `json:"club_name"`
}

func (server *Server) searchTeam(ctx *gin.Context) {
	var req searchTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	searchQuery := "%" + req.ClubName + "%"

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := server.store.SearchTeam(ctx, searchQuery)
	if err != nil {
		server.logger.Error("Failed to search team : %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getClubsBySport(ctx *gin.Context) {

	sports := ctx.Param("sport")

	response, err := server.store.GetClubsBySport(ctx, sports)
	if err != nil {
		server.logger.Error("Failed to get club by sport: %v", err)
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournamentsByClub(ctx *gin.Context) {
	clubName := ctx.Query("club_name")
	response, err := server.store.GetTournamentsByClub(ctx, clubName)
	if err != nil {
		server.logger.Error("Failed to get tournament by club: %v", err)
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getMatchByClubName(ctx *gin.Context) {
	clubIdStr := ctx.Query("id")
	clubID, err := strconv.ParseInt(clubIdStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse club id: %v", err)
		return
	}
	response, err := server.store.GetMatchByClubName(ctx, clubID)
	if err != nil {
		server.logger.Error("Failed to get match by clubname: %v", err)
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
