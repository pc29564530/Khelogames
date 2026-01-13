package teams

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
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
	fieldErrors := make(map[string]string)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors = errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// Latitude / Longitude validation
	if req.Latitude != "" {
		if _, err := strconv.ParseFloat(req.Latitude, 64); err != nil {
			fieldErrors["latitude"] = "Invalid format"
		}
	}

	if req.Longitude != "" {
		if _, err := strconv.ParseFloat(req.Longitude, 64); err != nil {
			fieldErrors["longitude"] = "Invalid format"
		}
	}

	// If ANY validation error exists → return once
	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	latitude, _ := strconv.ParseFloat(req.Latitude, 64)
	longitude, _ := strconv.ParseFloat(req.Longitude, 64)

	slug := util.GenerateSlug(req.Name)
	shortName := util.GenerateShortName(req.Name)

	team, err := s.txStore.CreateTeamsTx(
		ctx,
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
		s.logger.Error("Failed to create team: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create team",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    team,
	})
}

func (s *TeamsServer) GetTeamsFunc(ctx *gin.Context) {

	response, err := s.store.GetTeams(ctx)
	if err != nil {
		s.logger.Error("Failed to get teams: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get teams",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

type getClubRequest struct {
	PublicID string `uri:"public_id"`
}

func (s *TeamsServer) GetTeamFunc(ctx *gin.Context) {
	var req getClubRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetTeamByPublicID(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get team: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

type getTeamsBySportRequest struct {
	GameID int64 `uri:"game_id"`
}

func (s *TeamsServer) GetTeamsBySportFunc(ctx *gin.Context) {

	var req getTeamsBySportRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	rows, err := s.store.GetTeamsBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get teams by sport: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get teams by sport",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    rows,
	})
}

func (s *TeamsServer) GetTeamsByPlayerFunc(ctx *gin.Context) {

	playerPublicIDString := ctx.Param("player_public_id")

	playerPublicID, err := uuid.Parse(playerPublicIDString)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"player_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetTeamByPlayer(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get teams by player: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get teams by player",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

type getPlayersByTeamRequest struct {
	TeamPublicID string `uri:"team_public_id`
}

func (s *TeamsServer) GetPlayersByTeamFunc(ctx *gin.Context) {
	var req getPlayersByTeamRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	players, err := s.store.GetPlayerByTeam(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get players by team: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get players by team",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    players,
	})
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
	if err := ctx.ShouldBindUri(&reqURI); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(reqURI.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	latitude, err := strconv.ParseFloat(reqJSON.Latitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse latitude: ", err)
		fieldErrors := map[string]string{"latitude": "Invalid format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	longitude, err := strconv.ParseFloat(reqJSON.Longitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse longitude: ", err)
		fieldErrors := map[string]string{"longitude": "Invalid format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	team, err := s.txStore.UpdateTeamTx(ctx, teamPublicID, reqJSON.City, reqJSON.State, reqJSON.Country, latitude, longitude)
	if err != nil {
		s.logger.Error("Failed to update team location: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update team location",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    team,
	})
}
