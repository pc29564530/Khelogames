package cricket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *CricketServer) GetCurrentBowlerFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID uuid.UUID `json: "match_public_id"`
		TeamPublicID  uuid.UUID `json:"team_public_id"`
		InningNumber  int       `json:"inning_number"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		return
	}

	currentBowlerResponse, err := s.store.GetCurrentBowler(ctx, req.MatchPublicID, req.TeamPublicID, req.InningNumber)
	if err != nil {
		s.logger.Error("Failed to get current bowler score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if currentBowlerResponse == nil {
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"team": currentBowlerResponse.(map[string]interface{})["team"], "bowling": currentBowlerResponse.(map[string]interface{})["bowler"]})
}
