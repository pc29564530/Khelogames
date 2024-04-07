package api

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

func (server *Server) createTournamentStanding(ctx *gin.Context) {
	var req createTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
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
	response, err := server.store.CreateTournamentStanding(ctx, arg)
	if err != nil {
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

func (server *Server) getTournamentStanding(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")
	sport := ctx.Query("sport_type")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTournamentStandingParams{
		TournamentID: tournamentID,
		GroupID:      groupID,
		SportType:    sport,
	}

	response, err := server.store.GetTournamentStanding(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
