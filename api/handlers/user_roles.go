package handlers

import (
	"khelogames/core/token"
	db "khelogames/database"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) AssignRoleFunc(ctx *gin.Context) {
	var req struct {
		UserPublicID     string  `json:"user_public_id" binding:"required"`
		RoleName         string  `json:"role_name" binding:"required"`
		ResourceType     *string `json:"resource_type"`      // "tournament" | "match" | "team" | nil
		ResourcePublicID *string `json:"resource_public_id"` // UUID of tournament/match/team
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// Get target user profile
	userPublicID, err := uuid.Parse(req.UserPublicID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "INVALID_UUID", "message": "Invalid user_public_id"},
		})
		return
	}

	profile, err := s.store.GetProfileByPublicID(ctx, userPublicID)
	if err != nil || profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	// Resolve role by name
	role, err := s.store.GetRoleByName(ctx, req.RoleName)
	if err != nil || role == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "INVALID_ROLE", "message": "Role not found: " + req.RoleName},
		})
		return
	}

	// Resolve resource_id from public_id if provided
	var resourceID *int64
	if req.ResourceType != nil && req.ResourcePublicID != nil {
		resourcePublicID, err := uuid.Parse(*req.ResourcePublicID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   gin.H{"code": "INVALID_UUID", "message": "Invalid resource_public_id"},
			})
			return
		}

		rid, err := s.resolveResourceID(ctx, *req.ResourceType, resourcePublicID)
		if err != nil || rid == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   gin.H{"code": "NOT_FOUND", "message": "Resource not found"},
			})
			return
		}
		resourceID = &rid
	}

	assignedBy := int64(authPayload.UserID)
	assignment, err := s.store.AssignUserRole(ctx, db.AssignUserRoleParams{
		UserID:       int64(profile.UserID),
		RoleID:       role.ID,
		ResourceType: req.ResourceType,
		ResourceID:   resourceID,
		AssignedBy:   &assignedBy,
	})
	if err != nil {
		s.logger.Error("AssignRoleFunc: failed to assign role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "INTERNAL_ERROR", "message": "Failed to assign role"},
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    assignment,
	})
}

// Get User Roles

func (s *HandlersServer) GetUserRolesFunc(ctx *gin.Context) {
	profilePublicIDStr := ctx.Param("profile_public_id")
	profilePublicID, err := uuid.Parse(profilePublicIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "INVALID_UUID", "message": "Invalid profile_public_id"},
		})
		return
	}

	profile, err := s.store.GetProfileByPublicID(ctx, profilePublicID)
	if err != nil || profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	roles, err := s.store.GetUserRoles(ctx, int64(profile.UserID))
	if err != nil {
		s.logger.Error("GetUserRolesFunc: failed to get roles: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "INTERNAL_ERROR", "message": "Failed to get user roles"},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

// Get My Roles

func (s *HandlersServer) GetMyRolesFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	roles, err := s.store.GetUserRoles(ctx, int64(authPayload.UserID))
	if err != nil {
		s.logger.Error("GetMyRolesFunc: failed to get roles: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "INTERNAL_ERROR", "message": "Failed to get roles"},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

// Get Resource Roles

func (s *HandlersServer) GetResourceRolesFunc(ctx *gin.Context) {
	resourceType := ctx.Query("resource_type")
	resourcePublicID := ctx.Query("resource_public_id")

	if resourceType == "" || resourcePublicID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "MISSING_PARAMS", "message": "resource_type and resource_public_id are required"},
		})
		return
	}

	parsedUUID, err := uuid.Parse(resourcePublicID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "INVALID_UUID", "message": "Invalid resource_public_id"},
		})
		return
	}

	resourceID, err := s.resolveResourceID(ctx, resourceType, parsedUUID)
	if err != nil || resourceID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": "Resource not found"},
		})
		return
	}

	roles, err := s.store.GetResourceUserRoles(ctx, resourceType, resourceID)
	if err != nil {
		s.logger.Error("GetResourceRolesFunc: failed: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "INTERNAL_ERROR", "message": "Failed to get resource roles"},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

// Revoke Role

func (s *HandlersServer) RevokeRoleFunc(ctx *gin.Context) {
	var req struct {
		UserPublicID     string  `json:"user_public_id" binding:"required"`
		RoleName         string  `json:"role_name" binding:"required"`
		ResourceType     *string `json:"resource_type"`
		ResourcePublicID *string `json:"resource_public_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	userPublicID, err := uuid.Parse(req.UserPublicID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "INVALID_UUID", "message": "Invalid user_public_id"},
		})
		return
	}

	profile, err := s.store.GetProfileByPublicID(ctx, userPublicID)
	if err != nil || profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"code": "NOT_FOUND", "message": "User not found"},
		})
		return
	}

	role, err := s.store.GetRoleByName(ctx, req.RoleName)
	if err != nil || role == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"code": "INVALID_ROLE", "message": "Role not found"},
		})
		return
	}

	var resourceID *int64
	if req.ResourceType != nil && req.ResourcePublicID != nil {
		resourcePublicID, err := uuid.Parse(*req.ResourcePublicID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   gin.H{"code": "INVALID_UUID", "message": "Invalid resource_public_id"},
			})
			return
		}
		rid, err := s.resolveResourceID(ctx, *req.ResourceType, resourcePublicID)
		if err != nil || rid == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   gin.H{"code": "NOT_FOUND", "message": "Resource not found"},
			})
			return
		}
		resourceID = &rid
	}

	err = s.store.RevokeUserRole(ctx, db.RevokeUserRoleParams{
		UserID:       int64(profile.UserID),
		RoleID:       role.ID,
		ResourceType: req.ResourceType,
		ResourceID:   resourceID,
	})
	if err != nil {
		s.logger.Error("RevokeRoleFunc: failed to revoke role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "INTERNAL_ERROR", "message": "Failed to revoke role"},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role revoked successfully",
	})
}

// Get All Roles

func (s *HandlersServer) GetAllRolesFunc(ctx *gin.Context) {
	roles, err := s.store.GetAllRoles(ctx)
	if err != nil {
		s.logger.Error("GetAllRolesFunc: failed: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"code": "INTERNAL_ERROR", "message": "Failed to get roles"},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

// Resolve resource internal ID from public UUID

func (s *HandlersServer) resolveResourceID(ctx *gin.Context, resourceType string, publicID uuid.UUID) (int64, error) {
	switch resourceType {
	case "tournament":
		tournament, err := s.store.GetTournament(ctx, publicID)
		if err != nil || tournament == nil {
			return 0, err
		}
		return tournament.ID, nil

	case "match":
		match, err := s.store.GetTournamentMatchByMatchID(ctx, publicID)
		if err != nil || match == nil {
			return 0, err
		}
		return match.ID, nil

	case "team":
		team, err := s.store.GetTeamByPublicID(ctx, publicID)
		if err != nil || team == nil {
			return 0, err
		}
		return team.ID, nil

	default:
		return 0, nil
	}
}
