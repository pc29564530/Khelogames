package players

import (
	"fmt"
	"khelogames/core/token"
	db "khelogames/database"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type newPlayerRequest struct {
	Positions string `json:"positions" binding:"required,min=2,max=100"`
	Country   string `json:"country" binding:"required,min=2,max=100"`
	GameID    int64  `json:"game_id" binding:"required,min=1"`
}

func (s *PlayerServer) NewPlayerFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add player profile")
	var req newPlayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Requested data: ", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	userPlayer, err := s.store.GetProfile(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("unable to get the profile: %s", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	fullNameSlug := util.GenerateSlug(userPlayer.FullName)
	shortName := util.GenerateShortName(userPlayer.FullName)

	arg := db.NewPlayerParams{
		UserPublicID: authPayload.PublicID,
		GameID:       req.GameID,
		Name:         userPlayer.FullName,
		Slug:         fullNameSlug,
		ShortName:    shortName,
		MediaUrl:     userPlayer.AvatarUrl,
		Positions:    req.Positions,
		Country:      req.Country,
	}

	response, err := s.store.NewPlayer(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add player profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add player profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Added player profile: ", response)
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *PlayerServer) GetPlayerByProfilePublicIDFunc(ctx *gin.Context) {
	var req struct {
		ProfilePublicID string `uri:"profile_public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	profilePublicID, err := uuid.Parse(req.ProfilePublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: ", err)
		fieldErrors := map[string]string{"profile_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	player, err := s.store.GetPlayerByProfile(ctx, profilePublicID)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    player,
	})
}

func (s *PlayerServer) GetAllPlayerFunc(ctx *gin.Context) {
	response, err := s.store.GetAllPlayer(ctx)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player profiles",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("Successfully get the player profile: ", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *PlayerServer) GetPlayerFunc(ctx *gin.Context) {
	var req struct {
		PlayerPublicID string `uri:"public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	playerPublicID, err := uuid.Parse(req.PlayerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetPlayer(ctx, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("Successfully get the player profile: ", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *PlayerServer) GetPlayerSearchFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get player profile")
	playerName := ctx.Query("name")

	fieldErrors := make(map[string]string)
	if playerName == "" {
		fieldErrors["name"] = "Name query parameter is required"
	}
	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Debug("Parse the player id: ", playerName)

	response, err := s.store.SearchPlayer(ctx, playerName)
	if err != nil {
		s.logger.Error("Failed to search player profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to search player profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("Successfully get the player profile: ", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *PlayerServer) GetPlayerByCountry(ctx *gin.Context) {
	country := ctx.Query("country")

	fieldErrors := make(map[string]string)
	if country == "" {
		fieldErrors["country"] = "Country query parameter is required"
	}
	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetPlayersCountry(ctx, country)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player profiles by country",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully get all player profile: ", response)
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *PlayerServer) GetPlayersBySportFunc(ctx *gin.Context) {
	var req struct {
		GameID int32 `uri:"game_id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetPlayersBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get player profile: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player profiles by sport",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully get all player profile: ", response)
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}
