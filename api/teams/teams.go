package teams

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

type addTeamsRequest struct {
	Name        string `json:"name"`
	MediaURL    string `json:"media_url"`
	Gender      string `jsong:"gender"`
	National    bool   `json:"national"`
	Country     string `json:"country"`
	Type        string `json:"type"`
	Sports      string `json:"sports"`
	PlayerCount int    `json:"player_count"`
}

func (s *TeamsServer) AddTeam(ctx *gin.Context) {
	var req addTeamsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind create club request: ", err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	tx.Rollback()

	saveImageStruct := util.NewSaveImageStruct(s.logger)
	var path string
	if req.MediaURL != "" {
		b64Data := req.MediaURL[strings.IndexByte(req.MediaURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			s.logger.Error("Failed to decode string: ", err)
			return
		}

		path, err = saveImageStruct.SaveImageToFile(data, "image")
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to create file: ", err)
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
		tx.Rollback()
		s.logger.Error("Failed to create club: ", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TeamsServer) GetTeamsFunc(ctx *gin.Context) {

	response, err := s.store.GetTeams(ctx)
	if err != nil {
		s.logger.Error("Failed to get club: ", err)
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
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	response, err := s.store.GetTeam(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get club: ", err)
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
		s.logger.Error("Failed to get club by sport: ", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTeamByPlayerRequest struct {
	PlayerID int64 `uri:"player_id`
}

func (s *TeamsServer) GetTeamsByPlayer(ctx *gin.Context) {
	var req getTeamByPlayerRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	response, err := s.store.GetTeamByPlayer(ctx, req.PlayerID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
	}
	var teamDetails []map[string]interface{}
	for _, team := range response {

		teamDetail := map[string]interface{}{
			"id":         team.ID,
			"name":       team.Name,
			"gender":     team.Gender,
			"media_url":  team.MediaUrl,
			"short_name": team.Shortname,
			"slug":       team.Slug,
			"country":    team.Country,
			"national":   team.National,
		}
		teamDetails = append(teamDetails, teamDetail)

	}

	ctx.JSON(http.StatusAccepted, teamDetails)
}

type getPlayersByTeamRequest struct {
	TeamID int64 `uri:"teamID`
}

func (s *TeamsServer) GetPlayersByTeamFunc(ctx *gin.Context) {
	var req getPlayersByTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	response, err := s.store.GetPlayerByTeam(ctx, req.TeamID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
	}
	var teamDetails []map[string]interface{}
	for _, player := range response {

		teamDetail := map[string]interface{}{
			"id":          player.ID,
			"player_name": player.PlayerName,
			"slug":        player.Slug,
			"short_name":  player.ShortName,
			"position":    player.Positions,
			"country":     player.Country,
			"sports":      player.Sports,
			"media_url":   player.MediaUrl,
		}
		teamDetails = append(teamDetails, teamDetail)

	}

	ctx.JSON(http.StatusAccepted, teamDetails)
}

type UpdateCurrentTeamByPlayerRequest struct {
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
}

func (s *TeamsServer) UpdateCurrentTeamByPlayerFunc(ctx *gin.Context) {
	var req UpdateCurrentTeamByPlayerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.UpdateCurrentTeamParams{
		TeamID:   req.TeamID,
		PlayerID: req.PlayerID,
	}

	response, err := s.store.UpdateCurrentTeam(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the current team: ", err)
	}

	ctx.JSON(http.StatusAccepted, response)
}
