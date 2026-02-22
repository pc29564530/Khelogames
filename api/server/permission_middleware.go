package server

import (
	"context"
	"net/http"

	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/pkg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// Role names
const (
	RoleOrganizer       = "organizer"
	RoleTournamentAdmin = "tournament_admin"
	RoleScorer          = "scorer"
	RoleTeamManager     = "team_manager"
	RoleCoach           = "coach"
	RoleAdmin           = "admin"
	RoleNormal          = "normal"
)

// Resource types
const (
	ResourceTournament = "tournament"
	ResourceMatch      = "match"
	ResourceTeam       = "team"
)

// Permission → Roles that can perform it
var permissionRoles = map[string][]string{
	PermUpdateMatch: {
		RoleScorer,
		RoleTournamentAdmin,
		RoleOrganizer,
		RoleAdmin,
	},
	PermUpdateTournament: {
		RoleOrganizer,
		RoleTournamentAdmin,
		RoleAdmin,
	},
	PermUpdateTournamentAdmin: {
		RoleOrganizer,
		RoleAdmin,
	},
	PermUpdateTeam: {
		RoleTeamManager,
		RoleCoach,
		RoleAdmin,
	},
	PermUpdateCommunity: {
		RoleAdmin,
	},
}

// RequiredPermission Middleware

func (s *Server) RequiredPermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

		// First try URL path params
		matchPublicIDStr := ctx.Param("match_public_id")
		tournamentPublicIDStr := ctx.Param("tournament_public_id")
		teamPublicIDStr := ctx.Param("team_public_id")

		// Fallback: try query string
		if matchPublicIDStr == "" {
			matchPublicIDStr = ctx.Query("match_public_id")
		}
		if tournamentPublicIDStr == "" {
			tournamentPublicIDStr = ctx.Query("tournament_public_id")
		}
		if teamPublicIDStr == "" {
			teamPublicIDStr = ctx.Query("team_public_id")
		}

		// Fallback: try JSON request body (uses ShouldBindBodyWith so body is preserved for the handler)
		if matchPublicIDStr == "" && tournamentPublicIDStr == "" && teamPublicIDStr == "" {
			var body struct {
				MatchPublicID      string `json:"match_public_id"`
				TournamentPublicID string `json:"tournament_public_id"`
				TeamPublicID       string `json:"team_public_id"`
			}
			if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err == nil {
				matchPublicIDStr = body.MatchPublicID
				tournamentPublicIDStr = body.TournamentPublicID
				teamPublicIDStr = body.TeamPublicID
			}
		}

		var allowed bool
		var err error

		switch {
		case matchPublicIDStr != "":
			allowed, err = s.canPerformMatchAction(ctx, authPayload, matchPublicIDStr, permission)

		case tournamentPublicIDStr != "":
			allowed, err = s.canPerformTournamentAction(ctx, authPayload, tournamentPublicIDStr, permission)

		case teamPublicIDStr != "":
			allowed, err = s.canPerformTeamAction(ctx, authPayload, teamPublicIDStr, permission)

		default:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "CALL_ERROR",
					"message": "No resource context found",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		if err != nil || !allowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "You don't have permission to perform this action",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		ctx.Next()
	}
}

// Tournament Permission Check

func (s *Server) canPerformTournamentAction(
	ctx context.Context,
	authPayload *token.Payload,
	tournamentPublicIDStr string,
	permission string,
) (bool, error) {
	tournamentPublicID, err := uuid.Parse(tournamentPublicIDStr)
	if err != nil {
		return false, err
	}

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("canPerformTournamentAction: failed to get tournament: ", err)
		return false, nil
	}

	// Owner always has full access
	if authPayload.UserID == tournament.UserID {
		return true, nil
	}

	// Check via user_role_assignments for each allowed role
	return s.hasAnyRoleForResource(
		ctx,
		authPayload.UserID,
		permission,
		ResourceTournament,
		tournament.ID,
	)
}

// Match Permission Check

func (s *Server) canPerformMatchAction(
	ctx context.Context,
	authPayload *token.Payload,
	matchPublicIDStr string,
	permission string,
) (bool, error) {
	matchPublicID, err := uuid.Parse(matchPublicIDStr)
	if err != nil {
		return false, err
	}

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("canPerformMatchAction: failed to get match: ", err)
		return false, nil
	}

	// Check via user_role_assignments for match-level roles
	// Also check tournament-level roles (organizer of tournament can manage matches)
	matchAllowed, err := s.hasAnyRoleForResource(
		ctx,
		authPayload.UserID,
		permission,
		ResourceMatch,
		match.ID,
	)
	if err != nil {
		return false, err
	}
	if matchAllowed {
		return true, nil
	}

	// Fallback: check tournament-level role for this match's tournament
	tournamentAllowed, err := s.hasAnyRoleForResource(
		ctx,
		authPayload.UserID,
		permission,
		ResourceTournament,
		int64(match.TournamentID),
	)
	return tournamentAllowed, err
}

// Team Permission Check

func (s *Server) canPerformTeamAction(
	ctx context.Context,
	authPayload *token.Payload,
	teamPublicIDStr string,
	permission string,
) (bool, error) {
	teamPublicID, err := uuid.Parse(teamPublicIDStr)
	if err != nil {
		return false, err
	}

	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("canPerformTeamAction: failed to get team: ", err)
		return false, nil
	}

	// Owner always has access
	if authPayload.UserID == team.UserID {
		return true, nil
	}

	return s.hasAnyRoleForResource(
		ctx,
		authPayload.UserID,
		permission,
		ResourceTeam,
		team.ID,
	)
}

// Core DB Permission Check

// hasAnyRoleForResource checks if the user has any of the roles
// allowed for the given permission, scoped to the given resource
func (s *Server) hasAnyRoleForResource(
	ctx context.Context,
	userID int32,
	permission string,
	resourceType string,
	resourceID int64,
) (bool, error) {
	allowedRoles, exists := permissionRoles[permission]
	if !exists {
		return false, nil
	}

	rt := resourceType

	for _, roleName := range allowedRoles {
		rn := roleName
		allowed, err := s.store.HasRolePermission(ctx, db.HasPermissionParams{
			UserID:       int64(userID),
			RoleName:     rn,
			ResourceType: &rt,
			ResourceID:   &resourceID,
		})
		if err != nil {
			s.logger.Error("hasAnyRoleForResource: DB error: ", err)
			continue
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}
