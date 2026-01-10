package cricket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *CricketServer) GetCurrentBowlerFunc(ctx *gin.Context) {

	matchPublicIDStr := ctx.Query("match_public_id")
	teamPublicIDStr := ctx.Query("team_public_id")
	inningNumberStr := ctx.Query("inning_number")

	matchPublicID, err := uuid.Parse(matchPublicIDStr)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	teamPublicID, err := uuid.Parse(teamPublicIDStr)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	inningNumber, err := strconv.Atoi(inningNumberStr)
	if err != nil {
		s.logger.Error("Failed to parse to int: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid inning number format",
		})
		return
	}

	currentBowlerResponse, err := s.store.GetCurrentBowler(ctx, matchPublicID, teamPublicID, inningNumber)
	if err != nil {
		s.logger.Error("Failed to get current bowler score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get current bowler score",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, currentBowlerResponse)
}
