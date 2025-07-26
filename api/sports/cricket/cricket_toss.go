package cricket

import (
	db "khelogames/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addCricketTossRequest struct {
	MatchPublicID string `json:"match_public_id"`
	TossDecision  string `json:"toss_decision"`
	TossWin       string `json:"toss_win"`
}

func (s *CricketServer) AddCricketToss(ctx *gin.Context) {
	var req addCricketTossRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	tossWin, err := uuid.Parse(req.TossWin)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response, err := s.store.AddCricketToss(ctx, matchPublicID, req.TossDecision, tossWin)
	if err != nil {
		s.logger.Error("Failed to add cricket match toss : ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get the match by id: ", err)
		return
	}
	var teamID int32
	if req.TossDecision == "batting" {
		teamID = response.TossWin
	} else {
		if match.HomeTeamID != response.TossWin {
			teamID = match.AwayTeamID
		} else {
			teamID = match.HomeTeamID
		}
	}
	inningR := db.NewCricketScoreParams{
		MatchID:           int32(response.MatchID),
		TeamID:            int32(teamID),
		InningNumber:      1,
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
	MatchPublicID string `uri:"match_id" form:"match_id"`
}

func (s *CricketServer) GetCricketTossFunc(ctx *gin.Context) {

	var req getTossRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		s.logger.Error("Failed to bind cricket toss : ", err)
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	cricketToss, err := s.store.GetCricketToss(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get the cricket match toss: ", err)
		return
	}

	tossWonTeamMap, ok := cricketToss["toss_won_team"].(map[string]interface{})
	if !ok {
		s.logger.Error("Invalid toss_won_team format")
		return
	}

	publicIDStr, ok := tossWonTeamMap["public_id"].(string)
	if !ok {
		s.logger.Error("Invalid public_id format")
		return
	}

	publicID, err := uuid.Parse(publicIDStr)
	if err != nil {
		s.logger.Error("Failed to parse public_id as UUID: ", err)
		return
	}

	tossWonTeam, err := s.store.GetTeamByPublicID(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get team: ", err)
		return
	}

	tossDetails := map[string]interface{}{

		"tossWonTeam": map[string]interface{}{
			"id":        tossWonTeam.ID,
			"public_id": tossWonTeam.PublicID,
			"name":      tossWonTeam.Name,
			"slug":      tossWonTeam.Slug,
			"shortName": tossWonTeam.Shortname,
			"gender":    tossWonTeam.Gender,
			"national":  tossWonTeam.National,
			"country":   tossWonTeam.Country,
			"type":      tossWonTeam.Type,
		},
		"tossDecision": cricketToss["toss_decision"].(map[string]string),
	}

	s.logger.Debug("toss won team details: ", tossDetails)

	ctx.JSON(http.StatusAccepted, tossDetails)
}
