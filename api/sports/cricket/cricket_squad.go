package cricket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *CricketServer) AddCricketSquadFunc(ctx *gin.Context) {
	var req struct {
		MatchID  *int64 `json:"match_id"`
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

	cricketSquad, err := s.store.AddCricketSquad(ctx, *req.MatchID, req.TeamID, req.PlayerID, req.Role, req.OnBench)
	if err != nil {
		s.logger.Error("Failed to add football squad: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, cricketSquad)
}

func (s *CricketServer) GetCricketMatchSquadFunc(ctx *gin.Context) {
	var req struct {
		MatchID *int64 `json:"match_id"`
		TeamID  int64  `json:"team_id"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}

	cricketSquad, err := s.store.GetCricketMatchSquad(ctx, *req.MatchID, req.TeamID)
	if err != nil {
		s.logger.Error("Failed to get cricket squad: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, cricketSquad)
}
