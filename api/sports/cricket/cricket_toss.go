package cricket

import (
	db "khelogames/database"
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
		s.logger.Error("Failed to bind : ", err)
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
		s.logger.Error("Failed to add cricket match toss : ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	match, err := s.store.GetTournamentMatchByMatchID(ctx, arg.MatchID)
	if err != nil {
		s.logger.Error("Failed to get the match by id: ", err)
		return
	}
	var teamID int64
	if arg.TossDecision == "batting" {
		teamID = arg.TossWin
	} else {
		if match.HomeTeamID != arg.TossWin {
			teamID = match.AwayTeamID
		}
	}
	inningR := db.NewCricketScoreParams{
		MatchID:           response.ID,
		TeamID:            teamID,
		Inning:            "inning1",
		Score:             0,
		Wickets:           0,
		Overs:             0,
		RunRate:           "0.00",
		TargetRunRate:     "0.00",
		FollowOn:          false,
		IsInningCompleted: false,
		Declared:          false,
	}

	_, err = s.store.NewCricketScore(ctx, inningR)
	if err != nil {
		s.logger.Error("Failed to add the team score: ", err)
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
		s.logger.Error("Failed to get the cricket match toss: ", err)
		return
	}

	tossWonTeam, err := s.store.GetTeam(ctx, response.TossWin)
	if err != nil {
		s.logger.Error("Failed to get team: ", err)
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
