package football

import (
	db "khelogames/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addLineUpRequest struct {
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

func (s *FootballServer) AddFootballLineUpFunc(ctx *gin.Context) {
	var req addLineUpRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.AddFootballLineUpParams{
		TeamID:   req.TeamID,
		PlayerID: req.PlayerID,
		MatchID:  req.MatchID,
		Position: req.Position,
	}

	response, err := s.store.AddFootballLineUp(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the player in lineup: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type getLineUpRequest struct {
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

func (s *FootballServer) GetFootballLineUpFunc(ctx *gin.Context) {
	var req getLineUpRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.GetFootballLineUpParams{
		TeamID:  req.TeamID,
		MatchID: req.MatchID,
	}

	response, err := s.store.GetFootballLineUp(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get the player in lineup: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type updateSubsAndLineUpRequest struct {
	LineUpID int64 `json:"lineup_id"`
	SubsID   int64 `json:"subs_id"`
}

func (s *FootballServer) UpdateFootballSubsAndLineUpFunc(ctx *gin.Context) {
	var req updateSubsAndLineUpRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.UpdateFootballSubsAndLineUpParams{
		ID:   req.LineUpID,
		ID_2: req.SubsID,
	}

	response, err := s.store.UpdateFootballSubsAndLineUp(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the football and Subs")
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

// substitution code

type addSubsRequest struct {
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

func (s *FootballServer) AddFootballSubstitionFunc(ctx *gin.Context) {
	var req addSubsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.AddFootballSubstitutionParams{
		TeamID:   req.TeamID,
		PlayerID: req.PlayerID,
		MatchID:  req.MatchID,
		Position: req.Position,
	}

	response, err := s.store.AddFootballSubstitution(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the player in lineup: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type getSubstitutionpRequest struct {
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

func (s *FootballServer) GetFootballSubstitutionFunc(ctx *gin.Context) {
	var req getSubstitutionpRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.GetFootballSubstitutionParams{
		TeamID:  req.TeamID,
		MatchID: req.MatchID,
	}

	response, err := s.store.GetFootballSubstitution(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get the player in lineup: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

func (s *FootballServer) AddFootballSquadFunc(ctx *gin.Context) {
	var req struct {
		MatchID     int64  `json:"match_id"`
		TeamID      int64  `json:"team_id"`
		PlayerID    int64  `json:"player_id"`
		Position    string `json:"position"`
		IsSubstitue bool   `json:"is_substitue"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}

	footballSquad, err := s.store.AddFootballSquad(ctx, req.MatchID, req.TeamID, req.PlayerID, req.Position, req.IsSubstitue)
	if err != nil {
		s.logger.Error("Failed to add football squad: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, footballSquad)
}
