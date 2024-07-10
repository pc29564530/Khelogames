package tournaments

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TournamentStandingServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewTournamentStanding(store *db.Store, logger *logger.Logger) *TournamentStandingServer {
	return &TournamentStandingServer{store: store, logger: logger}
}

type createTournamentStandingRequest struct {
	TournamentID   int64 `json:"tournament_id"`
	GroupID        int64 `json:"group_id"`
	TeamID         int64 `json:"team_id"`
	Wins           int64 `json:"wins"`
	Loss           int64 `json:"loss"`
	Draw           int64 `json:"draw"`
	GoalFor        int64 `json:"goal_for"`
	GoalAgainst    int64 `json:"goal_against"`
	GoalDifference int64 `json:"goal_difference"`
	Points         int64 `json:"points"`
}

func (s *TournamentStandingServer) CreateTournamentStandingFunc(ctx *gin.Context) {
	var req createTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateTournamentStandingParams{
		TournamentID:   req.TournamentID,
		GroupID:        req.GroupID,
		TeamID:         req.TeamID,
		Wins:           req.Wins,
		Loss:           req.Loss,
		Draw:           req.Draw,
		GoalFor:        req.GoalFor,
		GoalAgainst:    req.GoalAgainst,
		GoalDifference: req.GoalDifference,
		Points:         req.Points,
	}
	response, err := s.store.CreateTournamentStanding(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create tournament standing: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTournamentStandingRequest struct {
	TournamentID int64  `json:"tournament_id"`
	GroupID      int64  `json:"group_id"`
	SportType    string `json:"sport_type"`
}

func (s *TournamentStandingServer) GetTournamentStandingFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")
	sport := ctx.Query("sport_type")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse group id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTournamentStandingParams{
		TournamentID: tournamentID,
		GroupID:      groupID,
		SportType:    sport,
	}

	response, err := s.store.GetTournamentStanding(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournament standing: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateTournamentStandingRequest struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentStandingServer) UpdateTournamentStandingFunc(ctx *gin.Context) {
	var req updateTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("bind the request: %v", req)
	arg := db.UpdateTournamentStandingParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	response, err := s.store.UpdateTournamentStanding(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update tournament standing: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("successfully tournament standing: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
