package cricket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *CricketServer) GetCurrentBatsmanFunc(ctx *gin.Context) {
	matchIDString := ctx.Query("match_id")
	teamIDString := ctx.Query("team_id")
	inningStr := ctx.Query("inning_number")
	matchID, err := strconv.ParseInt(matchIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match id ", err)
		return
	}

	teamID, err := strconv.ParseInt(teamIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse team id ", err)
		return
	}

	inning, err := strconv.Atoi(inningStr)
	if err != nil {
		s.logger.Error("Failed to parse inning ", err)
		return
	}

	currentBatsmansResponse, err := s.store.GetCurrentBatsman(ctx, matchID, teamID, inning)
	if err != nil {
		s.logger.Error("Failed to get current batsman score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"team": currentBatsmansResponse.(map[string]interface{})["team"], "batting": currentBatsmansResponse.(map[string]interface{})["batsman"]})
}
