package handlers

import (
	"khelogames/core/token"
	"khelogames/pkg"
	"net/http"

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
		return
	}

	eventPublicID, err := uuid.Parse(req.EventPublicID)
	if err != nil {
		s.logger.Error("Failed to parse event public id to uuid: ", err)
		return
	}

	//implement the tx for add location need to add the location id into the tournament, player, team
	location, err := s.txStore.AddLocationTx(ctx, eventPublicID, req.LocationOF, req.City, req.State, req.Country, req.Latitude, req.Longitude)
	if err != nil {
		s.logger.Error("Failed to add location: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, location)
}

func (s *HandlersServer) UpdateUserLocationFunc(ctx *gin.Context) {
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// need to update location id to the profile
	profile, err := s.store.GetProfileByUserID(ctx, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get profile by user id: ", err)
		return
	}

	location, err := s.store.UpdateUserLocation(ctx, profile.LocationID, req.Latitude, req.Longitude)
	if err != nil {
		s.logger.Error("Failed to update user location: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, location)
}
