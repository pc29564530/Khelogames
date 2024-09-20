package football

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
