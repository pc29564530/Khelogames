package cricket

import (
	"context"
	db "khelogames/database"
	"khelogames/database/models"
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

func (s *CricketServer) GetCricketScore(matches []db.GetMatchByIDRow, tournamentID int64) []map[string]interface{} {
	ctx := context.Background()

	tournament, err := s.store.GetTournament(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
	}

	var matchDetail []map[string]interface{}
	for _, match := range matches {

		homeTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.HomeTeamID}
		awayTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.AwayTeamID}
		homeScore, err := s.store.GetCricketScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for home team:", err)
		}
		awayScore, err := s.store.GetCricketScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for away team:", err)
		}

		var awayScoreMap map[string]interface{}
		var homeScoreMap map[string]interface{}
		var emptyScore models.CricketScore
		if awayScore != emptyScore {
			awayScoreMap = map[string]interface{}{"id": awayScore.ID, "score": awayScore.Score, "wickets": homeScore.Wickets, "overs": awayScore.Overs, "inning": awayScore.Inning, "runRate": awayScore.RunRate, "targetRunRate": awayScore.TargetRunRate}
		}

		if homeScore != emptyScore {
			homeScoreMap = map[string]interface{}{"id": homeScore.ID, "score": homeScore.Score, "wickets": homeScore.Wickets, "overs": homeScore.Overs, "inning": homeScore.Inning, "runRate": homeScore.RunRate, "targetRunRate": homeScore.TargetRunRate}
		}

		game, err := s.store.GetGame(ctx, match.HomeGameID)
		if err != nil {
			s.logger.Error("Failed to get the game: ", err)
		}

		matchMap := map[string]interface{}{
			"matchId":        match.ID,
			"tournament":     map[string]interface{}{"id": tournament.ID, "name": tournament.TournamentName, "slug": tournament.Slug, "country": tournament.Country, "sports": tournament.Sports},
			"homeTeam":       map[string]interface{}{"id": match.HomeTeamID, "name": match.HomeTeamName, "slug": match.HomeTeamSlug, "shortName": match.HomeTeamShortname, "gender": match.HomeTeamGender, "national": match.HomeTeamNational, "country": match.HomeTeamCountry, "type": match.HomeTeamType, "player_count": match.HomeTeamPlayerCount},
			"homeScore":      homeScoreMap,
			"awayTeam":       map[string]interface{}{"id": match.AwayTeamID, "name": match.AwayTeamName, "slug": match.AwayTeamSlug, "shortName": match.AwayTeamShortname, "gender": match.AwayTeamGender, "national": match.AwayTeamNational, "country": match.AwayTeamCountry, "type": match.AwayTeamType, "player_count": match.AwayTeamPlayerCount},
			"awayScore":      awayScoreMap,
			"startTimeStamp": match.StartTimestamp,
			"end_timestamp":  match.EndTimestamp,
			"status":         match.StatusCode,
			"game":           game,
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
