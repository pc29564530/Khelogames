package football

import (
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getLineUpRequest struct {
	MatchPublicID  uuid.UUID `json:"match_public_id"`
	TeamPublicID   uuid.UUID `json:"team_public_id"`
	PlayerPublicID uuid.UUID `json:"player_public_id"`
	Position       string    `json:"position"`
}

func (s *FootballServer) GetFootballLineUpFunc(ctx *gin.Context) {
	var req getLineUpRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	response, err := s.store.GetFootballMatchSquad(ctx, req.MatchPublicID, req.TeamPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player in lineup: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type Player struct {
	PublicID uuid.UUID `json:"public_id"`
}

type MatchSquadRequest struct {
	MatchPublicID uuid.UUID   `json:"match_public_id"`
	TeamPublicID  uuid.UUID   `json:"team_public_id"`
	Player        []Player    `json:"player"`
	IsSubstituted []uuid.UUID `json:"is_substituted"`
}

func (s *FootballServer) AddFootballSquadFunc(ctx *gin.Context) {

	var req MatchSquadRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}

	substitutedMap := make(map[uuid.UUID]bool)

	for _, substitutedID := range req.IsSubstituted {
		substitutedMap[substitutedID] = true
	}

	var footballSquad []map[string]interface{}
	for _, player := range req.Player {
		var squad models.FootballSquad
		var err error

		substitute := substitutedMap[player.PublicID]

		squad, err = s.store.AddFootballSquad(ctx, *&req.MatchPublicID, req.TeamPublicID, player.PublicID, substitute)
		if err != nil {
			s.logger.Error("Failed to add football squad: ", err)
			return
		}

		footballSquad = append(footballSquad, map[string]interface{}{
			"id":            squad.ID,
			"public_id":     squad.PublicID,
			"match_id":      squad.MatchID,
			"team_id":       squad.TeamID,
			"player":        player,
			"positions":     squad.Position,
			"is_substitute": squad.IsSubstitute,
			"role":          squad.Role,
			"created_at":    squad.CreatedAT,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Football squad added successfully",
		"squad":   footballSquad,
	})
}

func (s *FootballServer) GetFootballMatchSquadFunc(ctx *gin.Context) {

	matchIDString := ctx.Query("match_public_id")
	matchPublicID, err := uuid.Parse(matchIDString)
	if err != nil {
		s.logger.Error("`Failed to parse int: ", err)
		return
	}

	teamIDString := ctx.Query("team_public_id")
	teamPublicID, err := uuid.Parse(teamIDString)
	if err != nil {
		s.logger.Error("Failed to parse int: ", err)
		return
	}

	response, err := s.store.GetFootballMatchSquad(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get football match squad: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
