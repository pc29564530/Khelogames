package cricket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *CricketServer) AddFootballSquadFunc(ctx *gin.Context) {
	var req struct {
		MatchID  int64  `json:"match_id"`
		TeamID   int64  `json:"team_id"`
		PlayerID int64  `json:"player_id"`
		Role     string `json:"role"`
		OnBench  bool   `json:"on_bench"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}

	footballSquad, err := s.store.AddCricketSquad(ctx, req.MatchID, req.TeamID, req.PlayerID, req.Role, req.OnBench)
	if err != nil {
		s.logger.Error("Failed to add football squad: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, footballSquad)
}
