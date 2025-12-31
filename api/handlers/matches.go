package handlers

import (
	"fmt"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uber/h3-go/v4"
)

func (s *HandlersServer) GetMatchesByLocationFunc(ctx *gin.Context) {
	sport := ctx.Param("sport")

	// Get game by sport name
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		ctx.JSON(400, gin.H{"error": "Invalid sport"})
		return
	}

	// Parse query parameters
	startDateString := ctx.Query("start_timestamp")
	latitudeString := ctx.Query("latitude")
	longitudeString := ctx.Query("longitude")
	radiusString := ctx.DefaultQuery("radius", "10") // Default 10km

	// Validate required parameters
	if latitudeString == "" || longitudeString == "" {
		ctx.JSON(400, gin.H{"error": "latitude and longitude are required"})
		return
	}

	// Convert start date
	startDate, err := util.ConvertTimeStamp(startDateString)
	if err != nil {
		s.logger.Error("Failed to convert timestamp: ", err)
		startDate = 0
	}

	// Parse coordinates
	latitude, err := strconv.ParseFloat(latitudeString, 64)
	if err != nil {
		s.logger.Error("Failed to parse latitude: ", err)
		ctx.JSON(400, gin.H{"error": "Invalid latitude"})
		return
	}

	longitude, err := strconv.ParseFloat(longitudeString, 64)
	if err != nil {
		s.logger.Error("Failed to parse longitude: ", err)
		ctx.JSON(400, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(radiusString, 64)
	if err != nil {
		s.logger.Error("Failed to parse radius: ", err)
		radius = 10.0
	}

	// ✅ Step 1: Convert user location to H3
	latLng := h3.NewLatLng(latitude, longitude)
	userCell, err := h3.LatLngToCell(latLng, 9) // Resolution 9
	if err != nil {
		s.logger.Error("Failed to get cell from lat and long: ", err)
	}

	s.logger.Info("User H3 Index: ", userCell.String())

	// ✅ Step 2: Calculate k-ring based on radius
	kRing := calculateKRing(radius, 9)
	s.logger.Info("K-Ring: ", kRing, " for radius: ", radius, " km")

	// ✅ Step 3: Get neighboring H3 cells
	neighbors, err := h3.GridDisk(userCell, kRing)
	if err != nil {
		s.logger.Error("Failed to get the neighbors: ", err)
	}

	// Convert to string array
	h3Cells := make([]string, len(neighbors))
	for i, cell := range neighbors {
		h3Cells[i] = cell.String()
	}

	s.logger.Info("Searching in ", len(h3Cells), " H3 cells")

	// ✅ Step 4: Query matches with those H3 indexes
	// This queries locations with nearby H3, then gets matches at those locations
	listMatches, err := s.store.ListMatchesByLocation(
		ctx,
		int32(startDate),
		latitude,
		longitude,
		h3Cells,
		radius,
		int32(game.ID),
	)
	if err != nil {
		s.logger.Error("Failed to get matches by location: ", err)
		ctx.JSON(500, gin.H{"error": "Failed to fetch matches"})
		return
	}

	s.logger.Info("Found ", len(listMatches), " listMatches")

	fmt.Println("Match: ", listMatches)
	ctx.JSON(200, listMatches)
}

func calculateKRing(radiusKm float64, resolution int) int {
	edgeLengths := map[int]float64{
		7:  5.16,
		8:  0.46,
		9:  0.174,
		10: 0.066,
	}

	edgeLength, ok := edgeLengths[resolution]
	if !ok {
		edgeLength = 0.174
	}

	kRing := int(radiusKm / (2 * edgeLength))
	if kRing < 1 {
		kRing = 1
	}
	if kRing > 10 {
		kRing = 10
	}

	return kRing
}

func (s *HandlersServer) GetAllMatchesFunc(ctx *gin.Context) {

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}
	startDateString := ctx.Query("start_timestamp")
	startDate, err := util.ConvertTimeStamp(startDateString)
	if err != nil {
		s.logger.Error("Failed to convert to second: ", err)
	}
	response, err := s.store.ListMatches(ctx, int32(startDate), game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *HandlersServer) GetLiveMatchesFunc(ctx *gin.Context) {

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}
	response, err := s.store.GetLiveMatches(ctx, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *HandlersServer) GetMatchByMatchIDFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}

	match, err := s.store.GetMatchByPublicId(ctx, matchPublicID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by match id: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, match)
	return
}
