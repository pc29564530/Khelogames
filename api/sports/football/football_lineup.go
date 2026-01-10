package football

import (
	"khelogames/core/token"
	"khelogames/database/models"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getLineUpRequest struct {
	MatchPublicID  string `json:"match_public_id"`
	TeamPublicID   string `json:"team_public_id"`
	PlayerPublicID string `json:"player_public_id"`
	Position       string `json:"position"`
}

func (s *FootballServer) GetFootballLineUpFunc(ctx *gin.Context) {
	var req getLineUpRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetFootballMatchSquad(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player in lineup: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get the player in lineup",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type Player struct {
	PublicID string `json:"public_id"`
}

type MatchSquadRequest struct {
	MatchPublicID string   `json:"match_public_id"`
	TeamPublicID  string   `json:"team_public_id"`
	Player        []Player `json:"player"`
	IsSubstituted []string `json:"is_substituted"`
}

func (s *FootballServer) AddFootballSquadFunc(ctx *gin.Context) {

	var req MatchSquadRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match data: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get match details",
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get user tournament role",
		})
		return
	}
	if !isExists {
		s.logger.Error("User does not own this match")
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    "NOT_FOUND",
			"message": "You do not own this match",
		})
		return
	}

	substitutedMap := make(map[string]bool)

	for _, substitutedID := range req.IsSubstituted {
		substitutedMap[substitutedID] = true
	}

	var footballSquad []map[string]interface{}
	for _, player := range req.Player {
		var squad models.FootballSquad
		var err error

		substitute := substitutedMap[player.PublicID]

		playerPublicID, err := uuid.Parse(player.PublicID)
		if err != nil {
			s.logger.Error("Invalid UUID format", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			})
			return
		}

		squad, err = s.store.AddFootballSquad(ctx, matchPublicID, teamPublicID, playerPublicID, substitute)
		if err != nil {
			s.logger.Error("Failed to add football squad: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "DATABASE_ERROR",
				"message": "Failed to add football squad",
			})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	teamIDString := ctx.Query("team_public_id")
	teamPublicID, err := uuid.Parse(teamIDString)
	if err != nil {
		s.logger.Error("Failed to parse int: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetFootballMatchSquad(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get football match squad: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get football match squad",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
