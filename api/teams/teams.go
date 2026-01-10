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
	Name        string `json:"name" binding:"required,min=3,max=100"`
	MediaURL    string `json:"media_url" binding:"omitempty,url"`
	Gender      string `json:"gender" binding:"required,oneof=male female mixed"`
	National    bool   `json:"national"`
	Country     string `json:"country" binding:"required,min=2,max=100"`
	Type        string `json:"type" binding:"required,min=2,max=50"`
	PlayerCount int    `json:"player_count" binding:"required,min=1,max=100"`
	GameID      int32  `json:"game_id" binding:"required,min=1"`
	Latitude    string `json:"latitude" binding:"omitempty"`
	Longitude   string `json:"longitude" binding:"omitempty"`
	City        string `json:"city" binding:"required,min=2,max=100"`
	State       string `json:"state" binding:"required,min=2,max=100"`
}

func (s *TeamsServer) AddTeam(ctx *gin.Context) {
	var req addTeamsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind create club request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    "VALIDATION_ERROR",
				"message": "Invalid latitude format",
			})
			return
		}

		longitude, err = strconv.ParseFloat(req.Longitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse to float: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    "VALIDATION_ERROR",
				"message": "Invalid longitude format",
			})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to create team",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, team)
	return
}

func (s *TeamsServer) GetTeamsFunc(ctx *gin.Context) {

	response, err := s.store.GetTeams(ctx)
	if err != nil {
		s.logger.Error("Failed to get club: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get teams",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetTeamByPublicID(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get club: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get team",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	rows, err := s.store.GetTeamsBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get teams by sport",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, rows)
}

func (s *TeamsServer) GetTeamsByPlayerFunc(ctx *gin.Context) {

	playerPublicIDString := ctx.Param("player_public_id")

	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetTeamByPlayer(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get teams by player",
		})
		return
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	players, err := s.store.GetPlayerByTeam(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get club by sport: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get players by team",
		})
		return
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	err = ctx.ShouldBindJSON(&reqJSON)
	if err != nil {
		s.logger.Error("Failed to bind json: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	teamPublicID, err := uuid.Parse(reqURI.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	latitude, err := strconv.ParseFloat(reqJSON.Latitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse float: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid latitude format",
		})
		return
	}

	longitude, err := strconv.ParseFloat(reqJSON.Longitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse float: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid longitude format",
		})
		return
	}

	team, err := s.txStore.UpdateTeamTx(ctx, teamPublicID, reqJSON.City, reqJSON.State, reqJSON.Country, latitude, longitude)

	ctx.JSON(http.StatusAccepted, team)
}
