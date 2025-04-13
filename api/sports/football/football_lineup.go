package football

import (
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"
	"strconv"

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

	response, err := s.store.GetFootballMatchSquad(ctx, req.MatchID, req.TeamID)
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

type Player struct {
	ID         int64  `json:"id"`
	PlayerName string `json:"player_name"`
	ShortName  string `json:"short_name"`
	Slug       string `json:"slug"`
	Country    string `json:"country"`
	Position   string `json:"position"`
	MediaURL   string `json:"media_url"`
	Sports     string `json:"sports"`
}

type MatchSquadRequest struct {
	MatchID       *int64   `json:"match_id"`
	TeamID        int64    `json:"team_id"`
	Player        []Player `json:"player"`
	IsSubstituted []int64  `json:"is_substituted"`
}

func (s *FootballServer) AddFootballSquadFunc(ctx *gin.Context) {

	var req MatchSquadRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}

	substitutedMap := make(map[int64]bool)

	for _, substitutedID := range req.IsSubstituted {
		substitutedMap[substitutedID] = true
	}

	var footballSquad []map[string]interface{}
	for _, player := range req.Player {
		var squad models.FootballSquad
		var err error

		substitute := substitutedMap[player.ID]

		var role string
		squad, err = s.store.AddFootballSquad(ctx, *req.MatchID, req.TeamID, player.ID, player.Position, substitute, role)
		if err != nil {
			s.logger.Error("Failed to add football squad: ", err)
			return
		}

		footballSquad = append(footballSquad, map[string]interface{}{
			"id":            squad.ID,
			"match_id":      squad.MatchID,
			"team_id":       squad.TeamID,
			"player":        player,
			"positions":     squad.Position,
			"is_substitute": squad.IsSubstitute,
			"role":          squad.Role,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Football squad added successfully",
		"squad":   footballSquad,
	})
}

type GetMatchSquadRequest struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *FootballServer) GetFootballMatchSquadFunc(ctx *gin.Context) {

	matchIDString := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse int: ", err)
		return
	}

	teamIDString := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse int: ", err)
		return
	}

	response, err := s.store.GetFootballMatchSquad(ctx, matchID, teamID)
	if err != nil {
		s.logger.Error("Failed to get football match squad: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
