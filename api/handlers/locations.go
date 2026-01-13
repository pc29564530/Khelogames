package handlers

import (
	"fmt"
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) AddLocationFunc(ctx *gin.Context) {
	var req struct {
		LocationOF    string  `json:"location_of" binding:"required"`
		EventPublicID string  `json:"event_public_id" binding:"required"`
		City          string  `json:"city" binding:"required,min=2,max=100"`
		State         string  `json:"state" binding:"required,min=2,max=100"`
		Country       string  `json:"country" binding:"required,min=2,max=100"`
		Latitude      float64 `json:"latitude" binding:"required"`
		Longitude     float64 `json:"longitude" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	eventPublicID, err := uuid.Parse(req.EventPublicID)
	if err != nil {
		s.logger.Error("Failed to parse event public id to uuid: ", err)
		fieldErrors := map[string]string{"event_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	//implement the tx for add location need to add the location id into the tournament, player, team
	location, err := s.txStore.AddLocationTx(ctx, req.LocationOF, eventPublicID, req.City, req.State, req.Country, req.Latitude, req.Longitude)
	if err != nil {
		s.logger.Error("Failed to add location: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add location",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    location,
	})
}

func (s *HandlersServer) UpdateUserLocationFunc(ctx *gin.Context) {
	var req struct {
		Latitude  string `json:"latitude" binding:"required"`
		Longitude string `json:"longitude" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	latitude, err := strconv.ParseFloat(req.Latitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse to float: ", err)
		fieldErrors := map[string]string{"latitude": "Invalid latitude format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	longitude, err := strconv.ParseFloat(req.Longitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse to float: ", err)
		fieldErrors := map[string]string{"longitude": "Invalid longitude format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// need to update location id to the profile
	profile, err := s.store.GetProfileByUserID(ctx, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get profile by user id: ", err)
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

	fmt.Println("Latitude: ", latitude)
	fmt.Println("Longitude: ", longitude)

	var h3Index string
	location, err := s.store.UpdateUserLocation(ctx, *profile.LocationID, latitude, longitude, h3Index)
	if err != nil {
		s.logger.Error("Failed to update user location: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update user location",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    location,
	})
}
