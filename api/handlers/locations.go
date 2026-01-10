package handlers

import (
	"fmt"
	"khelogames/core/token"
	"khelogames/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) AddLocationFunc(ctx *gin.Context) {
	var req struct {
		LocationOF    string  `json:"city"`
		EventPublicID string  `json:"event_public_id"`
		City          string  `json:"city"`
		State         string  `json:"state"`
		Country       string  `json:"country"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	eventPublicID, err := uuid.Parse(req.EventPublicID)
	if err != nil {
		s.logger.Error("Failed to parse event public id to uuid: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid event UUID format",
		})
		return
	}

	//implement the tx for add location need to add the location id into the tournament, player, team
	location, err := s.txStore.AddLocationTx(ctx, req.LocationOF, eventPublicID, req.City, req.State, req.Country, req.Latitude, req.Longitude)
	if err != nil {
		s.logger.Error("Failed to add location: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to add location",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, location)
}

func (s *HandlersServer) UpdateUserLocationFunc(ctx *gin.Context) {
	var req struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	latitude, err := strconv.ParseFloat(req.Latitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse to float: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid latitude format",
		})
		return
	}

	longitude, err := strconv.ParseFloat(req.Longitude, 64)
	if err != nil {
		s.logger.Error("Failed to parse to float: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid longitude format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// need to update location id to the profile
	profile, err := s.store.GetProfileByUserID(ctx, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get profile by user id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get user profile",
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
			"code":    "INTERNAL_ERROR",
			"message": "Failed to update user location",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, location)
}
