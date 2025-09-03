package tournaments

import (
	"database/sql"
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	"khelogames/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TournamentServer) GetTournamentMatch(ctx *gin.Context) {

	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	sports := strings.TrimSpace(ctx.Param("sport"))
	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentPublicID, sports))
	s.logger.Debug("Tournament match params: ", req.TournamentPublicID)

	matches, err := s.store.GetMatchByTournamentPublicID(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament match: ", err)
		return
	}

	checkSportServer := util.NewCheckSport(s.store, s.logger)
	matchDetailsWithScore := checkSportServer.CheckSport(sports, matches, tournamentPublicID)

	s.logger.Info("successfully  get the tournament match: ", matchDetailsWithScore)
	ctx.JSON(http.StatusAccepted, matchDetailsWithScore)
}

type createTournamentMatchRequest struct {
	TournamentPublicID string  `json:"tournament_public_id"`
	AwayTeamPublicID   string  `json:"away_team_public_id"`
	HomeTeamPublicID   string  `json:"home_team_public_id"`
	StartTimestamp     string  `json:"start_timestamp"`
	EndTimestamp       string  `json:"end_timestamp"`
	Type               string  `json:"type"`
	StatusCode         string  `json:"status_code"`
	Result             *int64  `json:"result"`
	Stage              string  `json:"stage"`
	KnockoutLevelID    *int32  `json:"knockout_level_id"`
	MatchFormat        *string `json:"match_format"`
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

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament_public_id: ", err)
		return
	}
	homeTeamPublicID, err := uuid.Parse(req.HomeTeamPublicID)
	if err != nil {
		s.logger.Error("Failed to parse home_team_public_id: ", err)
		return
	}
	awayTeamPublicID, err := uuid.Parse(req.AwayTeamPublicID)
	if err != nil {
		s.logger.Error("Failed to parse away_team_public_id: ", err)
		return
	}

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
		TournamentPublicID: tournamentPublicID,
		AwayTeamPublicID:   awayTeamPublicID,
		HomeTeamPublicID:   homeTeamPublicID,
		StartTimestamp:     startTimeStamp,
		EndTimestamp:       endTimeStamp,
		Type:               req.Type,
		StatusCode:         req.StatusCode,
		Result:             req.Result,
		Stage:              req.Stage,
		KnockoutLevelID:    req.KnockoutLevelID,
		MatchFormat:        &matchFormat,
		DayNumber:          nil,
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
			MatchID:    int32(updatedMatchData.ID),
			TeamID:     int32(updatedMatchData.AwayTeamID),
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
			MatchID:    int32(updatedMatchData.ID),
			TeamID:     int32(updatedMatchData.HomeTeamID),
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
			MatchID:         int32(updatedMatchData.ID),
			TeamID:          int32(updatedMatchData.HomeTeamID),
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
			MatchID:         int32(updatedMatchData.ID),
			TeamID:          int32(updatedMatchData.AwayTeamID),
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
			TeamID:  int64(updatedMatchData.AwayTeamID),
		}

		awayScore, err := s.store.GetFootballScore(ctx, argAway)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		argHome := db.GetFootballScoreParams{
			MatchID: updatedMatchData.ID,
			TeamID:  int64(updatedMatchData.HomeTeamID),
		}

		homeScore, err := s.store.GetFootballScore(ctx, argHome)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		if awayScore.Goals > homeScore.Goals {
			_, err := s.store.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.AwayTeamID))
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		} else if homeScore.Goals > awayScore.Goals {
			_, err := s.store.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.HomeTeamID))
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		}
	}
}

func updateCricketStatusCode(ctx *gin.Context, updatedMatchData models.Match, gameID int64, s *TournamentServer, tx *sql.Tx) {
	if updatedMatchData.StatusCode == "finished" {

		awayScore, err := s.store.GetCricketScore(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.AwayTeamID))
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		homeScore, err := s.store.GetCricketScore(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.HomeTeamID))
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to get away score: ", err)
		}

		if awayScore.Score > homeScore.Score {
			_, err := s.store.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.AwayTeamID))
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		} else if homeScore.Score > awayScore.Score {
			_, err := s.store.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.HomeTeamID))
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update match result: ", err)
			}
		}
	}
}

type updateStatusRequest struct {
	MatchPublicID string `json:"match_public_id"`
	StatusCode    string `json:"status_code"`
}

func (s *TournamentServer) UpdateMatchStatusFunc(ctx *gin.Context) {

	var req updateStatusRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
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

	updatedMatchData, err := s.store.UpdateMatchStatus(ctx, matchPublicID, req.StatusCode)
	if err != nil {
		tx.Rollback()
		s.logger.Error("unable to update the match status: ", err)
		return
	}

	if updatedMatchData.StatusCode == "finished" {

		_, err := s.store.AddORUpdateFootballPlayerStats(ctx, matchPublicID)
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

	match, err := s.store.GetMatchByPublicId(ctx, matchPublicID, gameID.ID)
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
	MatchPublicID string `json:"match_public_id"`
	Result        string `json:"result"`
}

func (s *TournamentServer) UpdateMatchResultFunc(ctx *gin.Context) {
	var req updateMatchResultRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	resultPublicID, err := uuid.Parse(req.Result)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	match, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	team, err := s.store.GetTeamByPublicID(ctx, resultPublicID)
	if err != nil {
		s.logger.Error("Failed to team ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	response, err := s.store.UpdateMatchResult(ctx, int32(match.ID), int32(team.ID))
	if err != nil {
		s.logger.Error("Failed to update result: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	s.logger.Info("Successfully update match result")
	ctx.JSON(http.StatusAccepted, response)
}
