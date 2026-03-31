package handlers

import (
	db "khelogames/database"
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	h3 "github.com/uber/h3-go/v4"
)

func (s *HandlersServer) GetMatchesByLocationFunc(ctx *gin.Context) {
	sport := ctx.Param("sport")

	// Get game by sport name
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get game by name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Parse query parameters
	startDateString := ctx.Query("start_timestamp")
	latitudeString := ctx.Query("latitude")
	longitudeString := ctx.Query("longitude")
	radiusString := ctx.DefaultQuery("radius", "10") // Default 10km

	// Validate required parameters
	if latitudeString == "" || longitudeString == "" {
		fieldErrors := make(map[string]string)
		if latitudeString == "" {
			fieldErrors["latitude"] = "Latitude is required"
		}
		if longitudeString == "" {
			fieldErrors["longitude"] = "Longitude is required"
		}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Convert start date
	startDate, err := util.ConvertTimeStamp(startDateString)
	if err != nil {
		s.logger.Error("Failed to convert timestamp: ", err)
		fieldErrors := map[string]string{"start_timestamp": "Invalid timestamp format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Parse coordinates
	latitude, err := strconv.ParseFloat(latitudeString, 64)
	if err != nil {
		s.logger.Error("Failed to parse latitude: ", err)
		fieldErrors := map[string]string{"latitude": "Invalid latitude format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	longitude, err := strconv.ParseFloat(longitudeString, 64)
	if err != nil {
		s.logger.Error("Failed to parse longitude: ", err)
		fieldErrors := map[string]string{"longitude": "Invalid longitude format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	radius, err := strconv.ParseFloat(radiusString, 64)
	if err != nil {
		s.logger.Error("Failed to parse radius: ", err)
		radius = 10.0
	}

	//Convert user location to H3
	latLng := h3.NewLatLng(latitude, longitude)
	userCell, err := h3.LatLngToCell(latLng, 9) // Resolution 9
	if err != nil {
		s.logger.Error("Failed to get cell from lat and long: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to process location",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("User H3 Index: ", userCell.String())

	//Calculate k-ring based on radius
	kRing := calculateKRing(radius, 9)
	s.logger.Info("K-Ring: ", kRing, " for radius: ", radius, " km")

	//Get neighboring H3 cells
	neighbors, err := h3.GridDisk(userCell, kRing)
	if err != nil {
		s.logger.Error("Failed to get the neighbors: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to process neighboring locations",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Convert to string array
	h3Cells := make([]string, len(neighbors))
	for i, cell := range neighbors {
		h3Cells[i] = cell.String()
	}

	s.logger.Info("Searching in ", len(h3Cells), " H3 cells")

	// Query matches with those H3 indexes
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get matches by location",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	var matches []map[string]interface{}
	if game.Name == "badminton" {
		for _, match := range listMatches {
			matchPublicIDStr, ok := match["public_id"].(string)
			if !ok {
				s.logger.Error("Invalid match public_id format: ", match["public_id"])
				continue
			}

			matchPublicID, err := uuid.Parse(matchPublicIDStr)
			if err != nil {
				s.logger.Error("Failed to parse match public_id: ", err)
				continue
			}

			badmintonScore, err := s.store.GetBadmintonMatchScore(ctx, matchPublicID)
			if err != nil {
				s.logger.Error("Failed to get badminton match score: ", err)
				continue
			}
			match["homeScore"] = badmintonScore.HomeSetsWon
			match["awayScore"] = badmintonScore.AwaySetsWon
			matches = append(matches, match)
		}
	} else {
		matches = append(matches, listMatches...)
	}
	s.logger.Info("Found ", len(listMatches), " listMatches")
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    matches,
	})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get game by name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	startDateString := ctx.Query("start_timestamp")
	startDate, err := util.ConvertTimeStamp(startDateString)
	if err != nil {
		s.logger.Error("Failed to convert to second: ", err)
		fieldErrors := map[string]string{"start_timestamp": "Invalid timestamp format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	res, err := s.store.ListMatches(ctx, int32(startDate), game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get matches by game",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	var matches []map[string]interface{}
	if game.Name == "badminton" {
		for _, match := range res {
			matchPublicIDStr, ok := match["public_id"].(string)
			if !ok {
				s.logger.Error("Invalid match public_id format: ", match["public_id"])
				continue
			}

			matchPublicID, err := uuid.Parse(matchPublicIDStr)
			if err != nil {
				s.logger.Error("Failed to parse match public_id: ", err)
				continue
			}

			badmintonScore, err := s.store.GetBadmintonMatchScore(ctx, matchPublicID)
			if err != nil {
				s.logger.Error("Failed to get badminton match score: ", err)
				continue
			}
			match["homeScore"] = badmintonScore.HomeSetsWon
			match["awayScore"] = badmintonScore.AwaySetsWon
			matches = append(matches, match)
		}
	} else {
		matches = append(matches, res...)
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    matches,
	})
}

func (s *HandlersServer) GetLiveMatchesFunc(ctx *gin.Context) {

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid game name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	res, err := s.store.GetLiveMatches(ctx, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get live matches by game",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	var matches []map[string]interface{}
	if game.Name == "badminton" {
		for _, match := range res {
			matchPublicIDStr, ok := match["public_id"].(string)
			if !ok {
				s.logger.Error("Invalid match public_id format: ", match["public_id"])
				continue
			}

			matchPublicID, err := uuid.Parse(matchPublicIDStr)
			if err != nil {
				s.logger.Error("Failed to parse match public_id: ", err)
				continue
			}

			badmintonScore, err := s.store.GetBadmintonMatchScore(ctx, matchPublicID)
			if err != nil {
				s.logger.Error("Failed to get badminton match score: ", err)
				continue
			}
			match["homeScore"] = badmintonScore.HomeSetsWon
			match["awayScore"] = badmintonScore.AwaySetsWon
			matches = append(matches, match)
		}
	} else {
		matches = append(matches, res...)
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    matches,
	})
}

func (s *HandlersServer) GetTrendingMatchesFunc(ctx *gin.Context) {

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid game name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	res, err := s.store.GetTrendingMatches(ctx, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get trending matches by game",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    res,
	})
}

func (s *HandlersServer) GetMatchByMatchIDFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		s.logger.Error("Failed to bind URI parameters: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get game by name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	match, err := s.store.GetMatchByPublicId(ctx, matchPublicID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by match id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match by match id",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	// Fetch typed match model for score lookups
	matchData, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil || matchData == nil {
		s.logger.Error("Failed to get match model by public id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match model",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if sport == "football" {
		homeTeamArg := db.GetFootballScoreParams{MatchID: matchData.ID, TeamID: int64(matchData.HomeTeamID)}
		awayTeamArg := db.GetFootballScoreParams{MatchID: matchData.ID, TeamID: int64(matchData.AwayTeamID)}
		homeScore, err := s.store.GetFootballScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get football match score for home team:", err)
		}
		awayScore, err := s.store.GetFootballScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get football match score for away team: ", err)
		}

		var emptyScore models.FootballScore
		var hScore map[string]interface{}
		if homeScore != emptyScore {
			hScore = map[string]interface{}{
				"public_id":        homeScore.PublicID,
				"first_half":       homeScore.FirstHalf,
				"second_half":      homeScore.SecondHalf,
				"goals":            homeScore.Goals,
				"penalty_shootout": homeScore.PenaltyShootOut,
			}
		}
		match["homeScore"] = hScore
		var aScore map[string]interface{}
		if awayScore != emptyScore {
			aScore = map[string]interface{}{
				"public_id":        awayScore.PublicID,
				"first_half":       awayScore.FirstHalf,
				"second_half":      awayScore.SecondHalf,
				"goals":            awayScore.Goals,
				"penalty_shootout": awayScore.PenaltyShootOut,
			}
		}
		match["awayScore"] = aScore
	} else if sport == "cricket" {
		matchScore, err := s.store.GetCricketScores(ctx, int32(matchData.ID))
		if err != nil {
			s.logger.Error("Failed to get cricket scores: ", err)
		} else {
			var homeScore []models.CricketScore
			var awayScore []models.CricketScore
			for _, score := range matchScore {
				if matchData.HomeTeamID == score.TeamID {
					homeScore = append(homeScore, score)
				} else {
					awayScore = append(awayScore, score)
				}
			}
			match["homeScore"] = homeScore
			match["awayScore"] = awayScore
		}
	} else if sport == "badminton" {
		score, err := s.store.GetBadmintonMatchScore(ctx, matchData.PublicID)
		if err != nil {
			s.logger.Error("Failed to get badminton match score:", err)
		} else if score != nil {
			var hScore int
			var aScore int
			if score.HomeSetsWon != nil {
				hScore = *score.HomeSetsWon
			}
			if score.AwaySetsWon != nil {
				aScore = *score.AwaySetsWon
			}
			match["homeScore"] = hScore
			match["awayScore"] = aScore
		}
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    match,
	})
}
