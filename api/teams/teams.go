package teams

import (
	"khelogames/core/token"
	"khelogames/pkg"
	"strconv"

	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	City        string `json:"city"`
	State       string `json:"state"`
}

func (s *TeamsServer) AddTeam(ctx *gin.Context) {
	var req addTeamsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind create club request: ", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	slug := util.GenerateSlug(req.Name)
	shortName := util.GenerateShortName(req.Name)

	var emptyString string
	var latitude float64
	var longitude float64

	if req.Latitude != emptyString && req.Longitude != emptyString {
		latitude, err = strconv.ParseFloat(req.Latitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse to float: ", err)
			return
		}

		longitude, err = strconv.ParseFloat(req.Longitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse to float: ", err)
			return
		}
	}

	team, err := s.txStore.CreateTeamsTx(ctx,
		authPayload.PublicID,
		req.Name,
		slug,
		shortName,
		req.MediaURL,
		req.Gender,
		false,
		req.Type,
		int32(req.PlayerCount),
		req.GameID,
		req.City,
		req.State,
		req.Country,
		latitude,
		longitude,
	)
	if err != nil {
		s.logger.Error("Failed to create new team: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, team)
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
	PublicID string `uri:"public_id"`
}

func (s *TeamsServer) GetTeamFunc(ctx *gin.Context) {
	var req getClubRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response, err := s.store.GetTeamByPublicID(ctx, publicID)
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

	ctx.JSON(http.StatusAccepted, rows)
}

func (s *TeamsServer) GetTeamsByPlayerFunc(ctx *gin.Context) {

	playerPublicIDString := ctx.Param("player_public_id")

	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response, err := s.store.GetTeamByPlayer(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
	}

	ctx.JSON(http.StatusAccepted, response)
}

type getPlayersByTeamRequest struct {
	TeamPublicID string `uri:"team_public_id`
}

func (s *TeamsServer) GetPlayersByTeamFunc(ctx *gin.Context) {
	var req getPlayersByTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	players, err := s.store.GetPlayerByTeam(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
	}

	ctx.JSON(http.StatusAccepted, players)
}

func (s *TeamsServer) UpdateTeamLocationFunc(ctx *gin.Context) {
	var reqURI struct {
		TeamPublicID string `uri:"team_public_id"`
	}
	var reqJSON struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
		City      string `json:"city"`
		State     string `json:"state"`
		Country   string `json:"country"`
	}
	err := ctx.ShouldBindUri(&reqURI)
	if err != nil {
		s.logger.Error("Failed to bind uri: ", err)
		return
	}
	err = ctx.ShouldBindJSON(&reqJSON)
	if err != nil {
		s.logger.Error("Failed to bind json: ", err)
		return
	}

	teamPublicID, err := uuid.Parse(reqURI.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	latitude, err := strconv.ParseFloat(reqJSON.Latitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse float: ", err)
		return
	}

	longitude, err := strconv.ParseFloat(reqJSON.Longitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse float: ", err)
		return
	}

	team, err := s.txStore.UpdateTeamTx(ctx, teamPublicID, reqJSON.City, reqJSON.State, reqJSON.Country, latitude, longitude)

	ctx.JSON(http.StatusAccepted, team)
}
