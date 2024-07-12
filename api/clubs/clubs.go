package clubs

import (
	"encoding/base64"
	db "khelogames/db/sqlc"

	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type createClubRequest struct {
	ClubName  string `json:"club_name"`
	AvatarURL string `json:"avatar_url"`
	Sport     string `json:"sport"`
}

func (s *ClubServer) CreateClubFunc(ctx *gin.Context) {
	var req createClubRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind create club request: %v", err)
		return
	}
	saveImageStruct := util.NewSaveImageStruct(s.logger)
	var path string
	if req.AvatarURL != "" {
		b64Data := req.AvatarURL[strings.IndexByte(req.AvatarURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			s.logger.Error("Failed to decode string: %v", err)
			return
		}

		path, err = saveImageStruct.SaveImageToFile(data, "image")
		if err != nil {
			s.logger.Error("Failed to create file: %v", err)
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
		s.logger.Error("Failed to create club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetClubsFunc(ctx *gin.Context) {

	response, err := s.store.GetClubs(ctx)
	if err != nil {
		s.logger.Error("Failed to get club: %v", err)
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
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	response, err := s.store.GetClub(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetClubsBySportFunc(ctx *gin.Context) {

	sports := ctx.Param("sport")

	response, err := s.store.GetClubsBySport(ctx, sports)
	if err != nil {
		s.logger.Error("Failed to get club by sport: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
