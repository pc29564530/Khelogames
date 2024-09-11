package tournaments

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (s *TournamentServer) CreateTournamentStandingFunc(ctx *gin.Context) {
	var req createTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
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
		s.logger.Error("Failed to create tournament standing: ", err)
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

func (s *TournamentServer) GetTournamentStandingFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")
	sport := ctx.Query("sport_type")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse group id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTournamentStandingParams{
		TournamentID: tournamentID,
		GroupID:      groupID,
		Sports:       sport,
	}

	response, err := s.store.GetTournamentStanding(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournament standing: ", err)
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

func (s *TournamentServer) UpdateTournamentStandingFunc(ctx *gin.Context) {
	var req updateTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("bind the request: ", req)
	arg := db.UpdateTournamentStandingParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	response, err := s.store.UpdateTournamentStanding(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update tournament standing: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("successfully tournament standing: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
