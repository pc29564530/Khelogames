package teams

import (
	"encoding/base64"
	"fmt"
	db "khelogames/db/sqlc"

	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type addTeamsRequest struct {
	Name     string `json:"name"`
	MediaURL string `json:"media_url"`
	Gender   string `jsong:"gender"`
	National bool   `json:"national"`
	Country  string `json:"country"`
	Type     string `json:"type"`
	Sports   string `json:"sports"`
}

func (s *TeamsServer) AddTeam(ctx *gin.Context) {
	var req addTeamsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind create club request: %v", err)
		return
	}
	saveImageStruct := util.NewSaveImageStruct(s.logger)
	var path string
	if req.MediaURL != "" {
		b64Data := req.MediaURL[strings.IndexByte(req.MediaURL, ',')+1:]

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
	slug := util.GenerateSlug(req.Name)
	shortName := util.GenerateShortName(req.Name)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.NewTeamsParams{
		Name:      req.Name,
		Slug:      slug,
		Shortname: shortName,
		Admin:     authPayload.Username,
		MediaUrl:  path,
		Gender:    req.Gender,
		National:  req.National,
		Country:   req.Country,
		Type:      req.Type,
		Sports:    req.Sports,
	}

	response, err := s.store.NewTeams(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TeamsServer) GetTeamsFunc(ctx *gin.Context) {

	response, err := s.store.GetTeams(ctx)
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

func (s *TeamsServer) GetTeamFunc(ctx *gin.Context) {
	var req getClubRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	response, err := s.store.GetTeam(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TeamsServer) GetTeamsBySportFunc(ctx *gin.Context) {

	sports := ctx.Param("sport")
	response, err := s.store.GetTeamsBySport(ctx, sports)
	if err != nil {
		s.logger.Error("Failed to get club by sport: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	fmt.Println("Response: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}
