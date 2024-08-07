package football

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addFootballMatchScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
	GoalFor      int64 `json:"goal_for"`
	GoalAgainst  int64 `json:"goal_against"`
}

func (s *FootballServer) AddFootballMatchScoreFunc(ctx *gin.Context) {

	var req addFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind football match score: %v", err)
		return
	}

	arg := db.AddFootballMatchScoreParams{
		MatchID:      req.MatchID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
		GoalFor:      req.GoalFor,
		GoalAgainst:  req.GoalAgainst,
	}

	response, err := s.store.AddFootballMatchScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add football match score: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type updateFootballMatchScoreRequest struct {
	GoalFor     int64 `json:"goal_for"`
	GoalAgainst int64 `json:"goal_against"`
	MatchID     int64 `json:"match_id"`
	TeamID      int64 `json:"team_id"`
}

func (s *FootballServer) UpdateFootballMatchScoreFunc(ctx *gin.Context) {

	var req updateFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind update football match score: %v", err)
		return
	}

	arg := db.UpdateFootballMatchScoreParams{
		GoalFor:     req.GoalFor,
		GoalAgainst: req.GoalAgainst,
		MatchID:     req.MatchID,
		TeamID:      req.TeamID,
	}

	response, err := s.store.UpdateFootballMatchScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update football match score: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
