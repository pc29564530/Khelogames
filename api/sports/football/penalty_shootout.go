package football

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addPenaltyRequest struct {
	MatchID  int64 `json:"match_id"`
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
	Scored   bool  `json:"scored"`
}

func (s *FootballServer) AddFootballPenaltyFunc(ctx *gin.Context) {
	var req addPenaltyRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.AddFootballPenaltyShootoutParams{
		MatchID:  req.MatchID,
		TeamID:   req.TeamID,
		PlayerID: req.PlayerID,
		Scored:   req.Scored,
	}

	response, err := s.store.AddFootballPenaltyShootout(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the penalty: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}

type getPenaltyRequest struct {
	MatchID  int64 `json:"match_id"`
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
	Scored   bool  `json:"scored"`
}

func (s *FootballServer) GetFootballPenaltyFunc(ctx *gin.Context) {
	var req getPenaltyRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.GetFootballPenaltyShootoutParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.GetFootballPenaltyShootout(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the penalty: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}

type updatePenaltyRequest struct {
	ID     int64 `json:"id"`
	Scored bool  `json:"scored"`
}

func (s *FootballServer) UpdateFootballPenaltyFunc(ctx *gin.Context) {
	var req updatePenaltyRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.UpdateFootballPenaltyShootoutParams{
		ID:     req.ID,
		Scored: req.Scored,
	}

	response, err := s.store.UpdateFootballPenaltyShootout(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the penalty: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}
