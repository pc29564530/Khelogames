package cricket

import (
	db "khelogames/db/sqlc"
	"net/http"

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

type getTossRequest struct {
	MatchID int64 `json:"match_id" form:"match_id"`
}

func (s *CricketServer) GetCricketTossFunc(ctx *gin.Context) {

	var req getTossRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		s.logger.Error("Failed to bind cricket toss : ", err)
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	response, err := s.store.GetCricketToss(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to get the cricket match toss: %v", err)
		return
	}

	tossWonTeam, err := s.store.GetTeam(ctx, response.TossWin)
	if err != nil {
		s.logger.Error("Failed to get team: %v", err)
		return
	}

	tossDetails := map[string]interface{}{

		"tossWonTeam": map[string]interface{}{
			"id":        tossWonTeam.ID,
			"name":      tossWonTeam.Name,
			"slug":      tossWonTeam.Slug,
			"shortName": tossWonTeam.Shortname,
			"gender":    tossWonTeam.Gender,
			"national":  tossWonTeam.National,
			"country":   tossWonTeam.Country,
			"type":      tossWonTeam.Type,
		},
		"tossDecision": response.TossDecision,
	}

	s.logger.Debug("toss won team details: ", tossDetails)

	ctx.JSON(http.StatusAccepted, tossDetails)
}
