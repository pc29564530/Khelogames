package cricket

import (
	"context"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addCricketScoreRequest struct {
	MatchID       int64   `json:"match_id"`
	TeamID        int64   `json:"team_id"`
	Inning        string  `json:"inning"`
	Score         int32   `json:"score"`
	Wickets       int32   `json:"wickets"`
	Overs         int32   `json:"overs"`
	RunRate       float64 `json:"run_rate"`
	TargetRunRate float64 `json:"target_run_rate"`
}

func (s *CricketServer) AddCricketScoreFunc(ctx *gin.Context) {

	var req addCricketScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	runRate := strconv.FormatFloat(req.RunRate, 'f', 2, 64)
	targetRunRate := strconv.FormatFloat(req.TargetRunRate, 'f', 2, 64)

	arg := db.NewCricketScoreParams{
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		Inning:        req.Inning,
		Score:         req.Score,
		Wickets:       req.Wickets,
		Overs:         req.Overs,
		RunRate:       runRate,
		TargetRunRate: targetRunRate,
	}

	response, err := s.store.NewCricketScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (s *CricketServer) GetCricketScore(matches []db.Match, matchDetails []map[string]interface{}) []map[string]interface{} {
	ctx := context.Background()
	for i, match := range matches {

		homeTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.HomeTeamID}
		awayTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.AwayTeamID}
		matchScoreData1, err := s.store.GetCricketScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 1:", err)
			return nil
		}
		matchScoreData2, err := s.store.GetCricketScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 2:", err)
			return nil
		}

		matchDetails[i]["home_team_score"] = matchScoreData1.Score
		matchDetails[i]["home_team_wickets"] = matchScoreData1.Wickets
		matchDetails[i]["home_team_overs"] = matchScoreData1.Overs
		matchDetails[i]["home_team_innings"] = matchScoreData1.Inning
		matchDetails[i]["home_team_run_rate"] = matchScoreData1.RunRate
		matchDetails[i]["home_team_target_run_rate"] = matchScoreData1.TargetRunRate
		matchDetails[i]["away_team_score"] = matchScoreData2.Score
		matchDetails[i]["away_team_wickets"] = matchScoreData2.Wickets
		matchDetails[i]["away_team_overs"] = matchScoreData2.Overs
		matchDetails[i]["away_team_innings"] = matchScoreData2.Inning
		matchDetails[i]["away_team_run_rate"] = matchScoreData2.RunRate
		matchDetails[i]["away_team_target_run_rate"] = matchScoreData2.TargetRunRate
	}
	return matchDetails
}

type updateInningRequest struct {
	Inning  string `json:"inning"`
	MatchID int64  `json:"match_id"`
	TeamID  int64  `json:"team_id"`
}

func (s *CricketServer) UpdateCricketInningsFunc(ctx *gin.Context) {
	var req updateInningRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("unable to bind the json: ", err)
		return
	}

	arg := db.UpdateCricketInningsParams{
		Inning:  req.Inning,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateCricketInnings(ctx, arg)
	if err != nil {
		s.logger.Error("unable to update the innings: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
