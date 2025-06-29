package tournaments

import (
	"database/sql"
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	"khelogames/util"
	"net/http"
	"strconv"
	"strings"

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
	sports := strings.TrimSpace(ctx.Param("sport"))
	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentID, sports))
	s.logger.Debug("Tournament match params: ", tournamentID)

	matches, err := s.store.GetMatchByTournamentID(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament match: ", err)
		return
	}

	checkSportServer := util.NewCheckSport(s.store, s.logger)
	matchDetailsWithScore := checkSportServer.CheckSport(sports, matches, tournamentID)

	s.logger.Info("successfully  get the tournament match: ", matchDetailsWithScore)
	ctx.JSON(http.StatusAccepted, matchDetailsWithScore)
}

type createTournamentMatchRequest struct {
	ID              int64   `json:"id"`
	TournamentID    int64   `json:"tournament_id"`
	AwayTeamID      int64   `json:"away_team_id"`
	HomeTeamID      int64   `json:"home_team_id"`
	StartTimestamp  string  `json:"start_timestamp"`
	EndTimestamp    string  `json:"end_timestamp"`
	Type            string  `json:"type"`
	StatusCode      string  `json:"status_code"`
	Result          *int64  `json:"result"`
	Stage           string  `json:"stage"`
	KnockoutLevelID *int32  `json:"knockout_level_id"`
	MatchFormat     *string `json:"match_format"`
}

func (s *TournamentServer) CreateTournamentMatch(ctx *gin.Context) {

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	var req createTournamentMatchRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("bind the request: ", req)

	gameName := ctx.Param("sport")

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
	var matchFormat string
	if gameName == "cricket" {
		if req.MatchFormat != nil {
			matchFormat = *req.MatchFormat
		} else {
			s.logger.Error("Match format is required for cricket matches")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Match format is required for cricket matches. Please provide a valid match format.",
			})
			return
		}
	}

	arg := db.NewMatchParams{
		TournamentID:    req.TournamentID,
		AwayTeamID:      req.AwayTeamID,
		HomeTeamID:      req.HomeTeamID,
		StartTimestamp:  startTimeStamp,
		EndTimestamp:    endTimeStamp,
		Type:            req.Type,
		StatusCode:      req.StatusCode,
		Result:          req.Result,
		Stage:           req.Stage,
		KnockoutLevelID: req.KnockoutLevelID,
		MatchFormat:     &matchFormat,
	}

	s.logger.Debug("Create match params: ", arg)

	response, err := s.store.NewMatch(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create match: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	if gameName == "football" {

		argAway := db.NewFootballScoreParams{
			MatchID:    response.AwayTeamID,
			TeamID:     response.AwayTeamID,
			FirstHalf:  0,
			SecondHalf: 0,
			Goals:      0,
		}

		argHome := db.NewFootballScoreParams{
			MatchID:    response.HomeTeamID,
			TeamID:     response.HomeTeamID,
			FirstHalf:  0,
			SecondHalf: 0,
			Goals:      0,
		}

		_, err := s.store.NewFootballScore(ctx, argAway)
		if err != nil {
			s.logger.Error("Failed to add away score: ", err)
			return
		}

		_, err = s.store.NewFootballScore(ctx, argHome)
		if err != nil {
			s.logger.Error("Failed to add home score: ", err)
			return
		}

	} else if gameName == "cricket" {

	}
	s.logger.Debug("Successfully create match: ", response)
	s.logger.Info("Successfully create match")

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

func updateFootballStatusCode(ctx *gin.Context, updatedMatchData models.Match, gameID int64, s *TournamentServer, tx *sql.Tx) {
	if updatedMatchData.StatusCode == "not_started" {
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

		argStatisticsHome := db.CreateFootballStatisticsParams{
			MatchID:         updatedMatchData.ID,
			TeamID:          updatedMatchData.HomeTeamID,
			ShotsOnTarget:   0,
			TotalShots:      0,
			CornerKicks:     0,
			Fouls:           0,
			GoalkeeperSaves: 0,
			FreeKicks:       0,
			YellowCards:     0,
			RedCards:        0,
		}

		argStatisticsAway := db.CreateFootballStatisticsParams{
			MatchID:         updatedMatchData.ID,
			TeamID:          updatedMatchData.AwayTeamID,
			ShotsOnTarget:   0,
			TotalShots:      0,
			CornerKicks:     0,
			Fouls:           0,
			GoalkeeperSaves: 0,
			FreeKicks:       0,
			YellowCards:     0,
			RedCards:        0,
		}

		_, err = s.store.CreateFootballStatistics(ctx, argStatisticsHome)
		if err != nil {
			s.logger.Error("Failed to add the football statistics: ", err)
		}

		_, err = s.store.CreateFootballStatistics(ctx, argStatisticsAway)
		if err != nil {
			s.logger.Error("Failed to add the football statistics: ", err)
		}
	} else if updatedMatchData.StatusCode == "finished" {
		argAway := db.GetFootballScoreParams{
			MatchID: updatedMatchData.ID,
			TeamID:  updatedMatchData.AwayTeamID,
		}

		awayScore, err := s.store.GetFootballScore(ctx, argAway)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		argHome := db.GetFootballScoreParams{
			MatchID: updatedMatchData.ID,
			TeamID:  updatedMatchData.HomeTeamID,
		}

		homeScore, err := s.store.GetFootballScore(ctx, argHome)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		if awayScore.Goals > homeScore.Goals {
			_, err := s.store.UpdateMatchResult(ctx, updatedMatchData.ID, updatedMatchData.AwayTeamID)
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		} else if homeScore.Goals > awayScore.Goals {
			_, err := s.store.UpdateMatchResult(ctx, updatedMatchData.ID, updatedMatchData.HomeTeamID)
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		}
	}
}

func updateCricketStatusCode(ctx *gin.Context, updatedMatchData models.Match, gameID int64, s *TournamentServer, tx *sql.Tx) {
	if updatedMatchData.StatusCode == "finished" {
		argAway := db.GetCricketScoreParams{
			MatchID: updatedMatchData.ID,
			TeamID:  updatedMatchData.AwayTeamID,
		}

		awayScore, err := s.store.GetCricketScore(ctx, argAway)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		argHome := db.GetCricketScoreParams{
			MatchID: updatedMatchData.ID,
			TeamID:  updatedMatchData.HomeTeamID,
		}

		homeScore, err := s.store.GetCricketScore(ctx, argHome)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		if awayScore.Score > homeScore.Score {
			_, err := s.store.UpdateMatchResult(ctx, updatedMatchData.ID, updatedMatchData.AwayTeamID)
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		} else if homeScore.Score > awayScore.Score {
			_, err := s.store.UpdateMatchResult(ctx, updatedMatchData.ID, updatedMatchData.HomeTeamID)
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		}
	}
}

type updateStatusRequest struct {
	ID         int64  `json:"id"`
	StatusCode string `json:"status_code"`
}

func (s *TournamentServer) UpdateMatchStatusFunc(ctx *gin.Context) {

	var req updateStatusRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	game := ctx.Param("sport")

	gameID, err := s.store.GetGamebyName(ctx, game)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
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

	if updatedMatchData.StatusCode == "finished" {

		_, err := s.store.AddORUpdateFootballPlayerStats(ctx, req.ID)
		if err != nil {
			s.logger.Error("Failed to add or update player stats: ", err)
			return
		}

		if gameID.Name == "football" {
			updateFootballStatusCode(ctx, updatedMatchData, gameID.ID, s, tx)
		} else if gameID.Name == "cricket" {
			updateCricketStatusCode(ctx, updatedMatchData, gameID.ID, s, tx)
		}
	}

	s.logger.Info("successfully updated the match status")

	match, err := s.store.GetMatchByMatchID(ctx, updatedMatchData.ID, gameID.ID)
	if err != nil {
		s.logger.Error("Failed to get match: ", err)
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, match)
}

type updateMatchResultRequest struct {
	ID     int64 `json:"id"`
	Result int64 `json:"result"`
}

func (s *TournamentServer) UpdateMatchResultFunc(ctx *gin.Context) {
	var req updateMatchResultRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	response, err := s.store.UpdateMatchResult(ctx, req.ID, req.Result)
	if err != nil {
		s.logger.Error("Failed to update result: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	s.logger.Info("Successfully update match result")
	ctx.JSON(http.StatusAccepted, response)
}
