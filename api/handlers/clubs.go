package handlers

import (
	"encoding/base64"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClubServer struct {
	store  *db.Store
	logger *logger.Logger
}
type createClubRequest struct {
	ClubName  string `json:"club_name"`
	AvatarURL string `json:"avatar_url"`
	Sport     string `json:"sport"`
}

func NewClubServer(store *db.Store, logger *logger.Logger) *ClubServer {
	return &ClubServer{store: store, logger: logger}
}

func (s *ClubServer) CreateClubFunc(ctx *gin.Context) {
	var req createClubRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind create club request: %v", err)
		return
	}

	var path string
	if req.AvatarURL != "" {
		b64Data := req.AvatarURL[strings.IndexByte(req.AvatarURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			fmt.Errorf("Failed to decode string: %v", err)
			return
		}

		path, err = util.SaveImageToFile(data, "image")
		if err != nil {
			fmt.Errorf("Failed to create file: %v", err)
			return
		}
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateClubParams{
		ClubName:  req.ClubName,
		AvatarUrl: path,
		Sport:     req.Sport,
		Owner:     authPayload.Username,
	}

	response, err := s.store.CreateClub(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetClubsFunc(ctx *gin.Context) {

	response, err := s.store.GetClubs(ctx)
	if err != nil {
		fmt.Errorf("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getClubRequest struct {
	ID int64 `uri:"id"`
}

func (s *ClubServer) GetClubFunc(ctx *gin.Context) {
	var req getClubRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	response, err := s.store.GetClub(ctx, req.ID)
	if err != nil {
		fmt.Errorf("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url"`
	ClubName  string `json:"club_name"`
}

func (s *ClubServer) UpdateClubAvatarUrlFunc(ctx *gin.Context) {
	var req updateAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateAvatarUrlParams{
		AvatarUrl: req.AvatarUrl,
		ClubName:  req.ClubName,
	}

	response, err := s.store.UpdateAvatarUrl(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update avatar url: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
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
// 		ctx.JSON(http.StatusInternalServerError, (err))
// 		return
// 	}
// 	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
// 	arg := db.UpdateClubNameParams{
// 		ClubName: req.ClubName,
// 	}

// 	response, err :=s.store.UpdateClubName(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, (err))
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return
// }

type updateClubSport struct {
	ClubName string `json:"club_name"`
	Sport    string `json:"sport"`
}

func (s *ClubServer) UpdateClubSportFunc(ctx *gin.Context) {
	var req updateClubSport
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateClubSportParams{
		Sport:    req.Sport,
		ClubName: req.ClubName,
	}

	response, err := s.store.UpdateClubSport(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update club sport: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type searchTeamRequest struct {
	ClubName string `json:"club_name"`
}

func (s *ClubServer) SearchTeamFunc(ctx *gin.Context) {
	var req searchTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	searchQuery := "%" + req.ClubName + "%"

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := s.store.SearchTeam(ctx, searchQuery)
	if err != nil {
		fmt.Errorf("Failed to search team : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetClubsBySportFunc(ctx *gin.Context) {

	sports := ctx.Param("sport")

	response, err := s.store.GetClubsBySport(ctx, sports)
	if err != nil {
		fmt.Errorf("Failed to get club by sport: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetTournamentsByClubFunc(ctx *gin.Context) {
	clubName := ctx.Query("club_name")
	response, err := s.store.GetTournamentsByClub(ctx, clubName)
	if err != nil {
		fmt.Errorf("Failed to get tournament by club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetMatchByClubNameFunc(ctx *gin.Context) {
	clubIdStr := ctx.Query("id")
	clubID, err := strconv.ParseInt(clubIdStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse club id: %v", err)
		return
	}
	response, err := s.store.GetMatchByClubName(ctx, clubID)
	if err != nil {
		fmt.Errorf("Failed to get match by clubname: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
