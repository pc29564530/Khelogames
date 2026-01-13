package football

import (
	"khelogames/core/token"
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"
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
	s.logger.Info("Received request to get football lineup")
	var req getLineUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid team UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetFootballMatchSquad(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get the player in lineup: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get the player in lineup",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully retrieved football lineup")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
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
	s.logger.Info("Received request to add football squad")
	var req MatchSquadRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid team UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match data: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user tournament role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	if !isExists {
		s.logger.Error("User does not own this match")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You do not own this match",
			},
			"request_id": ctx.GetString("request_id"),
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
			s.logger.Error("Invalid player UUID format: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "Invalid player UUID format",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		squad, err = s.store.AddFootballSquad(ctx, matchPublicID, teamPublicID, playerPublicID, substitute)
		if err != nil {
			s.logger.Error("Failed to add football squad: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to add football squad",
				},
				"request_id": ctx.GetString("request_id"),
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

	s.logger.Info("Successfully added football squad")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Football squad added successfully",
			"squad":   footballSquad,
		},
	})
}

func (s *FootballServer) GetFootballMatchSquadFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get football match squad")

	matchIDString := ctx.Query("match_public_id")
	teamIDString := ctx.Query("team_public_id")

	fieldErrors := make(map[string]string)

	if matchIDString == "" {
		fieldErrors["match_public_id"] = "Match public ID is required"
	}

	if teamIDString == "" {
		fieldErrors["team_public_id"] = "Team public ID is required"
	}

	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(matchIDString)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(teamIDString)
	if err != nil {
		s.logger.Error("Invalid team UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetFootballMatchSquad(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get football match squad: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get football match squad",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully retrieved football match squad")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
