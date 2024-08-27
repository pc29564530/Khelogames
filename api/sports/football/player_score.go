package football

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addFootballPlayerScoreRequest struct {
	MatchID  int64 `json:"match_id"`
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
	GoalTime int64 `json:"goal_time"`
}

func (s *FootballServer) AddFootballGoalByPlayerFunc(ctx *gin.Context) {

	var req addFootballPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind add football goal: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddFootballPlayerScoreParams{
		MatchID:  req.MatchID,
		TeamID:   req.TeamID,
		PlayerID: req.PlayerID,
		GoalTime: req.GoalTime,
	}

	response, err := s.store.AddFootballPlayerScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add football player score: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type getFootballPlayerScoreRequest struct {
	MatchID  int64 `json:"match_id"`
	PlayerID int64 `json:"player_id"`
}

func (s *FootballServer) GetFootballPlayerScoreFunc(ctx *gin.Context) {
	var req addFootballPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind add football goal: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.GetFootballPlayerScoreParams{
		MatchID:  req.MatchID,
		PlayerID: req.PlayerID,
	}

	response, err := s.store.GetFootballPlayerScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add football player score: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type countGoalByPlayerTeamRequest struct {
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
}

func (s *FootballServer) CountGoalByPlayerTeamFunc(ctx *gin.Context) {
	var req countGoalByPlayerTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind add football goal: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CountGoalByPlayerTeamParams{
		TeamID:   req.TeamID,
		PlayerID: req.PlayerID,
	}

	response, err := s.store.CountGoalByPlayerTeam(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add football player score: ", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
