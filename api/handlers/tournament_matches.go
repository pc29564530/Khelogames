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
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

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

	//server.GetLogger().Info("Create match params: %v", arg)

	response, err := s.store.CreateMatch(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	//server.GetLogger().Info("Create match response: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentMatchServer) GetAllTournamentMatchFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	sports := ctx.Query("sports")
	arg := db.GetTournamentMatchParams{
		TournamentID: tournamentID,
		Sports:       sports,
	}

	response, err := s.store.GetTournamentMatch(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to get tournament match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

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
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.UpdateMatchScheduleTimeParams{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		MatchID:   req.MatchID,
	}

	response, err := s.store.UpdateMatchScheduleTime(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update match schedule time: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getMatchRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
}

func (s *TournamentMatchServer) GetMatchFunc(ctx *gin.Context) {
	// var req getMatchRequest
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	ctx.JSON(http.StatusInternalServerError, err)
	// 	return
	// }

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse match id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}

	arg := db.GetMatchParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
	}

	response, err := s.store.GetMatch(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to get match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}
