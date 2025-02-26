package teams

import (
	"encoding/base64"
	"encoding/json"
	db "khelogames/database"

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
	PlayerCount int    `json:"player_count"`
	GameID      int32  `json:"game_id"`
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

	defer tx.Rollback()

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
			s.logger.Error("Failed to create file: ", err)
			return
		}
	}
	slug := util.GenerateSlug(req.Name)
	shortName := util.GenerateShortName(req.Name)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.NewTeamsParams{
		Name:        req.Name,
		Slug:        slug,
		Shortname:   shortName,
		Admin:       authPayload.Username,
		MediaUrl:    path,
		Gender:      req.Gender,
		National:    false,
		Country:     req.Country,
		Type:        req.Type,
		PlayerCount: int32(req.PlayerCount),
		GameID:      req.GameID,
	}

	response, err := s.store.NewTeams(ctx, arg)
	if err != nil {
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

type getTeamsBySportRequest struct {
	GameID int64 `uri:"game_id"`
}

func (s *TeamsServer) GetTeamsBySportFunc(ctx *gin.Context) {

	var req getTeamsBySportRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	rows, err := s.store.GetTeamsBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	var result map[string]interface{}
	var teamData []map[string]interface{}
	var gameDetail map[string]interface{}

	for _, row := range rows {
		var teamDetails map[string]interface{}
		err := json.Unmarshal(row.TeamData, &teamDetails)
		if err != nil {
			s.logger.Error("Failed to unmarshal team data: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process team data"})
			return
		}

		teamData = append(teamData, teamDetails)
		gameDetail = map[string]interface{}{
			"id":          row.ID,
			"name":        row.Name,
			"min_players": row.MinPlayers,
		}
	}

	result = map[string]interface{}{
		"game":  gameDetail,
		"teams": teamData,
	}

	ctx.JSON(http.StatusAccepted, result)
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
