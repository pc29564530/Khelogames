package cricket

import (
	errorhandler "khelogames/error_handler"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *CricketServer) GetCurrentBowlerFunc(ctx *gin.Context) {
	matchPublicIDStr := ctx.Query("match_public_id")
	teamPublicIDStr := ctx.Query("team_public_id")
	inningNumberStr := ctx.Query("inning_number")

	// Validate required query parameters
	if matchPublicIDStr == "" || teamPublicIDStr == "" || inningNumberStr == "" {
		fieldErrors := make(map[string]string)
		if matchPublicIDStr == "" {
			fieldErrors["match_public_id"] = "Match public ID is required"
		}
		if teamPublicIDStr == "" {
			fieldErrors["team_public_id"] = "Team public ID is required"
		}
		if inningNumberStr == "" {
			fieldErrors["inning_number"] = "Inning number is required"
		}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(matchPublicIDStr)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(teamPublicIDStr)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	inningNumber, err := strconv.Atoi(inningNumberStr)
	if err != nil {
		s.logger.Error("Failed to parse to int: ", err)
		fieldErrors := map[string]string{"inning_number": "Invalid inning number format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	currentBowlerResponse, err := s.store.GetCurrentBowler(ctx, matchPublicID, teamPublicID, inningNumber)
	if err != nil {
		s.logger.Error("Failed to get current bowler score : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get current bowler score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    currentBowlerResponse,
	})
}
