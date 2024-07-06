package handlers

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TournamentMatchServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewTournamentMatch(store *db.Store, logger *logger.Logger) *TournamentMatchServer {
	return &TournamentMatchServer{store: store, logger: logger}
}

type createTournamentMatchRequest struct {
	OrganizerID  int64     `json:"organizer_id"`
	TournamentID int64     `json:"tournament_id"`
	Team1ID      int64     `json:"team1_id"`
	Team2ID      int64     `json:"team2_id"`
	DateON       time.Time `json:"date_on"`
	StartTime    time.Time `json:"start_time"`
	Stage        string    `json:"stage"`
	Sports       string    `json:"sports"`
	EndTime      time.Time `json:"end_time"`
}

func (s *TournamentMatchServer) CreateTournamentMatchFunc(ctx *gin.Context) {
	var req createTournamentMatchRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)
	arg := db.CreateMatchParams{
		OrganizerID:  req.OrganizerID,
		TournamentID: req.TournamentID,
		Team1ID:      req.Team1ID,
		Team2ID:      req.Team2ID,
		DateOn:       req.DateON,
		StartTime:    req.StartTime,
		Stage:        req.Stage,
		Sports:       req.Sports,
		EndTime:      req.EndTime,
	}

	s.logger.Debug("Create match params: %v", arg)

	response, err := s.store.CreateMatch(ctx, arg)
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

//Changes are required to make like we can call getMatchscore() for team based match
//Can divide category wise like football, cricket etc.

func (s *TournamentMatchServer) GetAllTournamentMatchFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	sports := ctx.Query("sports")
	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentID, sports))
	arg := db.GetTournamentMatchParams{
		TournamentID: tournamentID,
		Sports:       sports,
	}

	s.logger.Debug("Tournament match params: %v", arg)

	response, err := s.store.GetTournamentMatch(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournament match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("successfully  get the tournament match: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateMatchScheduleTimeRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	MatchID   int64     `json:"match_id"`
}

func (s *TournamentMatchServer) updateMatchScheduleTime(ctx *gin.Context) {
	var req updateMatchScheduleTimeRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("bind the request: %v", req)
	arg := db.UpdateMatchScheduleTimeParams{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		MatchID:   req.MatchID,
	}
	s.logger.Debug(fmt.Sprintf("params arg: %v", arg))
	response, err := s.store.UpdateMatchScheduleTime(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update match schedule time: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug(fmt.Sprintf("update schedule time: %v", response))
	ctx.JSON(http.StatusAccepted, response)
	return
}

type getMatchRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
}

func (s *TournamentMatchServer) GetMatchFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	s.logger.Debug("get the request tournament_id: %v and match_id: %v", tournamentID, matchID)

	arg := db.GetMatchParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
	}

	response, err := s.store.GetMatch(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug(fmt.Sprintf("successfullly get the match: %v", response))

	ctx.JSON(http.StatusAccepted, response)
	return

}
