package cricket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *CricketServer) GetCurrentBatsmanFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `json: "match_public_id"`
		TeamPublicID  string `json:"team_public_id"`
		InningNumber  int    `json:"inning_number"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	currentBatsmansResponse, err := s.store.GetCurrentBatsman(ctx, matchPublicID, teamPublicID, req.InningNumber)
	if err != nil {
		s.logger.Error("Failed to get current batsman score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"team": currentBatsmansResponse.(map[string]interface{})["team"], "batting": currentBatsmansResponse.(map[string]interface{})["batsman"]})
}
