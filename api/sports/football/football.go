package football

import (
	"context"
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addFootballMatchScoreRequest struct {
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
	FirstHalf  int32 `json:"first_half"`
	SecondHalf int32 `json:"second_half"`
	Goals      int64 `json:"goal_for"`
}

func (s *FootballServer) AddFootballMatchScoreFunc(ctx *gin.Context) {

	var req addFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind football match score: %v", err)
		return
	}

	arg := db.NewFootballScoreParams{
		MatchID:    req.MatchID,
		TeamID:     req.TeamID,
		FirstHalf:  req.FirstHalf,
		SecondHalf: req.SecondHalf,
		Goals:      req.Goals,
	}

	response, err := s.store.NewFootballScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add football match score: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type getFootballScoreRequest struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

// func (s *FootballServer) GetFootballScore(ctx *gin.Context) {
// 	var req getFootballScoreRequest
// 	err := ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		s.logger.Error("Failed to bind football score: %v", err)
// 		return
// 	}

// 	arg := db.GetFootballScoreParams{
// 		MatchID: req.MatchID,
// 		TeamID:  req.TeamID,
// 	}

// 	response, err := s.store.GetFootballScore(ctx, arg)
// 	if err != nil {
// 		s.logger.Error("Failed to fetch football match score: %v", err)
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return
// }

func (s *FootballServer) GetFootballScore(matches []db.Match, tournament db.Tournament) []map[string]interface{} {
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

		homeTeamArg := db.GetFootballScoreParams{MatchID: match.ID, TeamID: match.HomeTeamID}
		awayTeamArg := db.GetFootballScoreParams{MatchID: match.ID, TeamID: match.AwayTeamID}
		homeScore, err := s.store.GetFootballScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 1:", err)
			return nil
		}
		awayScore, err := s.store.GetFootballScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 2:", err)
			return nil
		}

		matchMap := map[string]interface{}{
			"tournament":     map[string]interface{}{"id": tournament.ID, "name": tournament.TournamentName, "slug": tournament.Slug, "country": tournament.Country, "sports": tournament.Sports},
			"homeTeam":       map[string]interface{}{"id": homeTeam.ID, "name": homeTeam.Name, "slug": homeTeam.Slug, "shortName": homeTeam.Shortname, "gender": homeTeam.Gender, "national": homeTeam.National, "country": homeTeam.Country, "type": homeTeam.Type},
			"homeScore":      map[string]interface{}{"id": homeScore.ID, "score": homeScore.Goals, "firstHalf": homeScore.FirstHalf, "secondHalf": homeScore.SecondHalf},
			"awayTeam":       map[string]interface{}{"id": awayTeam.ID, "name": awayTeam.Name, "slug": awayTeam.Slug, "shortName": awayTeam.Shortname, "gender": awayTeam.Gender, "national": awayTeam.National, "country": awayTeam.Country, "type": awayTeam.Type},
			"awayScore":      map[string]interface{}{"id": awayScore.ID, "score": awayScore.Goals, "firstHalf": homeScore.FirstHalf, "secondHalf": homeScore.SecondHalf},
			"startTimeStamp": match.StartTimestamp,
			"end_timestamp":  match.EndTimestamp,
			"status":         match.StatusCode,
		}
		matchDetail = append(matchDetail, matchMap)

	}

	return matchDetail

}

type updateFootballMatchScoreRequest struct {
	Goals   int64 `json:"goals"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *FootballServer) UpdateFootballMatchScoreFunc(ctx *gin.Context) {

	var req updateFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind update football match score: %v", err)
		return
	}

	arg := db.UpdateFootballScoreParams{
		Goals:   req.Goals,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.UpdateFootballScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update football match score: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateFootballMatchScoreFirstHalfRequest struct {
	FirstHalf int32 `json:"first_half"`
	MatchID   int64 `json:"match_id"`
	TeamID    int64 `json:"team_id"`
}

func (s *FootballServer) UpdateFootballMatchScoreFirstHalfFunc(ctx *gin.Context) {

	var req updateFootballMatchScoreFirstHalfRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind update football match score: %v", err)
		return
	}

	arg := db.UpdateFirstHalfScoreParams{
		FirstHalf: req.FirstHalf,
		MatchID:   req.MatchID,
		TeamID:    req.TeamID,
	}

	response, err := s.store.UpdateFirstHalfScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update football match score: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateFootballMatchScoreSecondHalfRequest struct {
	SecondHalf int32 `json:"second_half"`
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
}

func (s *FootballServer) UpdateFootballMatchScoreSecondHalfFunc(ctx *gin.Context) {

	var req updateFootballMatchScoreSecondHalfRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind update football match score: %v", err)
		return
	}

	arg := db.UpdateSecondHalfScoreParams{
		SecondHalf: req.SecondHalf,
		MatchID:    req.MatchID,
		TeamID:     req.TeamID,
	}

	response, err := s.store.UpdateSecondHalfScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update football match score: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
