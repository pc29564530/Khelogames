package cricket

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addCricketTossRequest struct {
	MatchID      int64  `json:"match_id"`
	TossDecision string `json:"toss_decision"`
	TossWin      int64  `json:"toss_win"`
}

func (s *CricketServer) AddCricketToss(ctx *gin.Context) {
	var req addCricketTossRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	arg := db.AddCricketTossParams{
		MatchID:      req.MatchID,
		TossDecision: req.TossDecision,
		TossWin:      req.TossWin,
	}

	response, err := s.store.AddCricketToss(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add cricket match toss : %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *CricketServer) GetCricketTossFunc(ctx *gin.Context) {

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match id: %v", err)
		return
	}

	response, err := s.store.GetCricketToss(ctx, matchID)
	if err != nil {
		s.logger.Error("Failed to get the cricket match toss: %v", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
