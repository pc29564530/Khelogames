package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"khelogames/core/token"
	"khelogames/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var RolePermissions = map[string][]string{
	"SCORER": {
		PermUpdateMatch,
	},
	"TEAM_MANAGER": {
		PermUpdateTeam,
	},
	"TOURNAMENT_ADMIN": {
		PermUpdateMatch,
		PermUpdateTournament,
	},
	"TOURNAMENT_ORGANIZER": {
		PermUpdateMatch,
		PermUpdateTournament,
		PermUpdateTournamentAdmin,
		PermUpdateTeam,
	},
}

func (s *Server) RequiredPermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
		matchPublicIDStr := ctx.Param("match_public_id")
		tournamentPublicIDStr := ctx.Param("tournament_public_id")
		fmt.Println("Tournament: ", tournamentPublicIDStr)
		fmt.Println("Match: ", matchPublicIDStr)
		var allowed bool
		var err error

		switch {
		case matchPublicIDStr != "":
			//Match public id
			allowed, err = s.canPerformMatchAction(
				ctx,
				authPayload,
				matchPublicIDStr,
				permission,
			)
		case tournamentPublicIDStr != "":
			allowed, err = s.canPerformTournamentAction(
				ctx,
				authPayload,
				tournamentPublicIDStr,
				permission,
			)
		default:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "CALL_ERROR",
					"message": "No context found",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		if err != nil || !allowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "AUTHENTICATION_ERROR",
					"message": "Not allowed to make change",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		ctx.Next()
	}
}

func (s *Server) canPerformTeamAction(ctx context.Context, authPayload *token.Payload, teamPublicIDStr string, permission string) (bool, error) {
	teamPublicID, err := uuid.Parse(teamPublicIDStr)
	if err != nil {
		return false, err
	}
	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team public_id: ", err)
		return false, nil
	}

	if authPayload.UserID == team.UserID {
		return true, nil
	} else {
		return false, nil
	}
	return hasPermission("team_manager", permission), nil
}

func (s *Server) canPerformTournamentAction(ctx context.Context, authPayload *token.Payload, tournamentPublicIDStr string, permission string) (bool, error) {
	tournamentPublicID, err := uuid.Parse(tournamentPublicIDStr)
	if err != nil {
		return false, err
	}

	fmt.Println("Line no 102: ")

	// GetTournament expects UUID, let's use it
	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		return false, nil
	}
	fmt.Println("Tou User: ", tournament.UserID)
	fmt.Println("Auth: ", authPayload.UserID)
	if authPayload.UserID == tournament.UserID {
		return true, nil
	} else {
		return false, nil
	}
	return hasPermission("tournament_organizer", permission), nil
}

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

	userRole, err := s.store.GetMatchUserRole(ctx, matchPublicID, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user role: ", err)
		return false, err
	}

	if userRole.UserID != authPayload.UserID {
		fmt.Println("No match")
		return false, nil
	}

	return hasPermission(userRole.Role, permission), nil
}

func hasPermission(role string, permission string) bool {
	permissions, exists := RolePermissions[strings.ToUpper(role)]
	if !exists {
		return false
	}
	for _, perm := range permissions {
		fmt.Println("Perm: ", perm)
		fmt.Println("Permission: ", permission)
		if perm == permission {
			return true
		}
	}
	return false
}
