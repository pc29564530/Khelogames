package tournaments

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *TournamentServer) GetTournamentMatch(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	sports := ctx.Query("sports")
	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentID, sports))
	s.logger.Debug("Tournament match params: %v", tournamentID)

	matches, err := s.store.GetMatchesByTournamentID(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	tournament, err := s.store.GetTournament(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var matchDetails []map[string]interface{}
	for _, matchData := range matches {
		// argScore1 := db.GetFootballMatchScoreParams{
		// 	MatchID:      matchData.MatchID,
		// 	TeamID:       matchData.Team1ID,
		// 	TournamentID: matchData.TournamentID,
		// }
		// argScore2 := db.GetFootballMatchScoreParams{
		// 	MatchID:      matchData.MatchID,
		// 	TeamID:       matchData.Team2ID,
		// 	TournamentID: matchData.TournamentID,
		// }
		// _, err := s.store.GetFootballMatchScore(ctx, argScore1)
		// if err != nil {
		// 	s.logger.Error("Failed to get football match score for team1: %v", err)
		// 	continue
		// }

		// _, err = s.store.GetFootballMatchScore(ctx, argScore2)
		// if err != nil {
		// 	s.logger.Error("Failed to get football match score for team1: %v", err)
		// 	continue
		// }

		homeTeamID, err1 := s.store.GetTeam(ctx, matchData.HomeTeamID)
		if err1 != nil {
			s.logger.Error("Failed to get club details for team1: %v", err1)
			continue
		}
		awayTeamID, err2 := s.store.GetTeam(ctx, matchData.AwayTeamID)
		if err2 != nil {
			s.logger.Error("Failed to get club details for team2: %v", err2)
			continue
		}

		matchDetail := map[string]interface{}{
			"tournament_id":   matchData.TournamentID,
			"tournament_name": tournament.TournamentName,
			"match_id":        matchData.ID,
			"home_team_id":    matchData.HomeTeamID,
			"away_team_id":    matchData.AwayTeamID,
			"away_team_name":  awayTeamID.Name,
			"home_team_name":  homeTeamID.Name,
			"start_time":      matchData.StartTimestamp,
			"sports":          matchData.EndTimestamp,
		}
		//s.logger.Debug("football match details: %v ", matchDetails)
		matchDetails = append(matchDetails, matchDetail)

	}
	checkSportServer := util.NewCheckSport(s.store, s.logger)
	tournamentMatches := checkSportServer.CheckSport(sports, matches, matchDetails)
	s.logger.Info("successfully  get the tournament match: %v", tournamentMatches)
	ctx.JSON(http.StatusAccepted, tournamentMatches)
	return
}

type createTournamentMatchRequest struct {
	ID             int64  `json:"id"`
	TournamentID   int64  `json:"tournament_id"`
	AwayTeamID     int64  `json:"away_team_id"`
	HomeTeamID     int64  `json:"home_team_id"`
	StartTimestamp string `json:"start_timestamp"`
	EndTimestamp   string `json:"end_timestamp"`
	StatusCode     int64  `json:"status_code"`
	Type           string `json:"type"`
}

func (s *TournamentServer) CreateTournamentMatch(ctx *gin.Context) {
	var req createTournamentMatchRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)

	startTimeStamp, err := util.ConvertTimeStamp(req.StartTimestamp)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		return
	}
	var endTimeStamp int64
	if req.EndTimestamp != "" {
		endTimeStamp, err = util.ConvertTimeStamp(req.EndTimestamp)
		if err != nil {
			s.logger.Error("unable to convert time to second: ", err)
			return
		}
	}

	arg := db.NewMatchParams{
		TournamentID:   req.TournamentID,
		AwayTeamID:     req.AwayTeamID,
		HomeTeamID:     req.HomeTeamID,
		StartTimestamp: startTimeStamp,
		EndTimestamp:   endTimeStamp,
		StatusCode:     req.StatusCode,
		Type:           req.Type,
	}

	s.logger.Debug("Create match params: %v", arg)

	response, err := s.store.NewMatch(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("Successfully create match: %v", response)
	s.logger.Info("Successfully create match")

	ctx.JSON(http.StatusAccepted, response)
	return
}
