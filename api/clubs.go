package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createClubRequest struct {
	ClubName  string `json:"club_name"`
	AvatarURL string `json:"avatar_url,omit_emtpy"`
	Sport     string `json:"sport,omit_empty"`
}

func (server *Server) createClub(ctx *gin.Context) {
	var req createClubRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//add the path for avatar url
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateClubParams{
		ClubName:  req.ClubName,
		AvatarUrl: req.AvatarURL,
		Owner:     authPayload.Username,
		Sport:     req.Sport,
	}
	fmt.Println("arg: ", arg)

	response, err := server.store.CreateClub(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}
	fmt.Println("Club: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getClubs(ctx *gin.Context) {

	response, err := server.store.GetClubs(ctx)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response, err := server.store.GetClub(ctx, req.ID)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	searchQuery := "%" + req.ClubName + "%"
	fmt.Println("SearchQuieru: ", req.ClubName)

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := server.store.SearchTeam(ctx, searchQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("Response: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getClubsBySport(ctx *gin.Context) {

	sports := ctx.Param("sport")

	response, err := server.store.GetClubsBySport(ctx, sports)
	if err != nil {
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
		fmt.Errorf("Unable to parse the id: ", err)
		return
	}
	response, err := server.store.GetMatchByClubName(ctx, clubID)
	if err != nil {
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}

	fmt.Println("Match: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}
