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
		s.logger.Error("Failed to parse tournament id: ", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	sports := ctx.Query("sports")
	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentID, sports))
	s.logger.Debug("Tournament match params: ", tournamentID)

	matches, err := s.store.GetMatchByID(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament match: ", err)
		return
	}

	checkSportServer := util.NewCheckSport(s.store, s.logger)
	matchDetailsWithScore := checkSportServer.CheckSport(sports, matches, tournamentID)

	s.logger.Info("successfully  get the tournament match: ", matchDetailsWithScore)
	ctx.JSON(http.StatusAccepted, matchDetailsWithScore)
	return
}

type createTournamentMatchRequest struct {
	ID             int64  `json:"id"`
	TournamentID   int64  `json:"tournament_id"`
	AwayTeamID     int64  `json:"away_team_id"`
	HomeTeamID     int64  `json:"home_team_id"`
	StartTimestamp string `json:"start_timestamp"`
	EndTimestamp   string `json:"end_timestamp"`
	Type           string `json:"type"`
	StatusCode     string `json:"status_code"`
}

func (s *TournamentServer) CreateTournamentMatch(ctx *gin.Context) {
	var req createTournamentMatchRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

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
		Type:           req.Type,
		StatusCode:     req.StatusCode,
	}

	s.logger.Debug("Create match params: ", arg)

	response, err := s.store.NewMatch(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create match: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("Successfully create match: ", response)
	s.logger.Info("Successfully create match")

	ctx.JSON(http.StatusAccepted, response)
}

type updateStatusRequest struct {
	ID         int64  `json:"id"`
	StatusCode string `json:"status_code"`
	Sports     string `json:"sports"`
}

func (s *TournamentServer) UpdateMatchStatusFunc(ctx *gin.Context) {

	var req updateStatusRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("unable to begin tx: ", err)
		return
	}

	arg := db.UpdateMatchStatusParams{
		ID:         req.ID,
		StatusCode: req.StatusCode,
	}

	updatedMatchData, err := s.store.UpdateMatchStatus(ctx, arg)
	if err != nil {
		tx.Rollback()
		s.logger.Error("unable to update the match status: ", err)
		return
	}

	s.logger.Info("successfully updated the match status")

	if updatedMatchData.StatusCode == "in_progress" && req.Sports == "Football" {
		argAway := db.NewFootballScoreParams{
			MatchID:    updatedMatchData.ID,
			TeamID:     updatedMatchData.AwayTeamID,
			FirstHalf:  0,
			SecondHalf: 0,
			Goals:      0,
		}

		_, err := s.store.NewFootballScore(ctx, argAway)
		if err != nil {
			tx.Rollback()
			s.logger.Error("unable to add the football match score: ", err)
			return
		}

		argHome := db.NewFootballScoreParams{
			MatchID:    updatedMatchData.ID,
			TeamID:     updatedMatchData.HomeTeamID,
			FirstHalf:  0,
			SecondHalf: 0,
			Goals:      0,
		}

		_, err = s.store.NewFootballScore(ctx, argHome)
		if err != nil {
			tx.Rollback()
			s.logger.Error("unable to add the football match score: ", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, updatedMatchData)
}
