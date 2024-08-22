package cricket

import (
	"context"
	"fmt"
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

func (s *CricketServer) GetCricketScore(matches []db.Match, tournament db.Tournament) []map[string]interface{} {
	ctx := context.Background()
	var matchDetail []map[string]interface{}
	for _, match := range matches {

		homeTeam, err1 := s.store.GetTeam(ctx, match.HomeTeamID)
		if err1 != nil {
			s.logger.Error("Failed to get club details for team1: %v", err1)
			continue
		}
		awayTeam, err2 := s.store.GetTeam(ctx, match.AwayTeamID)
		if err2 != nil {
			s.logger.Error("Failed to get club details for team2: %v", err2)
			continue
		}

		homeTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.HomeTeamID}
		awayTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.AwayTeamID}
		homeScore, err := s.store.GetCricketScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 1:", err)
		}
		awayScore, err := s.store.GetCricketScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 2:", err)
		}

		var awayScoreMap map[string]interface{}
		var homeScoreMap map[string]interface{}
		var emptyScore db.CricketScore
		if awayScore != emptyScore {
			awayScoreMap = map[string]interface{}{"id": awayScore.ID, "score": awayScore.Score, "wickets": homeScore.Wickets, "overs": awayScore.Overs, "inning": awayScore.Inning, "runRate": awayScore.RunRate, "targetRunRate": awayScore.TargetRunRate}
		}

		if homeScore != emptyScore {
			homeScoreMap = map[string]interface{}{"id": homeScore.ID, "score": homeScore.Score, "wickets": homeScore.Wickets, "overs": homeScore.Overs, "inning": homeScore.Inning, "runRate": homeScore.RunRate, "targetRunRate": homeScore.TargetRunRate}
		}
		fmt.Println("Home: ", homeScoreMap)
		fmt.Println("Away: ", awayScore)

		matchMap := map[string]interface{}{
			"matchId":        match.ID,
			"tournament":     map[string]interface{}{"id": tournament.ID, "name": tournament.TournamentName, "slug": tournament.Slug, "country": tournament.Country, "sports": tournament.Sports},
			"homeTeam":       map[string]interface{}{"id": homeTeam.ID, "name": homeTeam.Name, "slug": homeTeam.Slug, "shortName": homeTeam.Shortname, "gender": homeTeam.Gender, "national": homeTeam.National, "country": homeTeam.Country, "type": homeTeam.Type},
			"homeScore":      homeScoreMap,
			"awayTeam":       map[string]interface{}{"id": awayTeam.ID, "name": awayTeam.Name, "slug": awayTeam.Slug, "shortName": awayTeam.Shortname, "gender": awayTeam.Gender, "national": awayTeam.National, "country": awayTeam.Country, "type": awayTeam.Type},
			"awayScore":      awayScoreMap,
			"startTimeStamp": match.StartTimestamp,
			"end_timestamp":  match.EndTimestamp,
			"status":         match.StatusCode,
		}
		matchDetail = append(matchDetail, matchMap)
	}
	return matchDetail
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
