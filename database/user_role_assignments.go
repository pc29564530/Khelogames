package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"
)

// Create User Role
const assignUserRole = `
	INSERT INTO user_role_assignments (
		user_id,
		role_id,
		resource_type,
		resource_id,
		assigned_by
	) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (user_id, role_id, resource_type, resource_id)
	DO UPDATE SET is_active = TRUE
	RETURNING id, user_id, role_id, resource_type, resource_id, assigned_by, is_active, created_at
`

type AssignUserRoleParams struct {
	UserID       int64
	RoleID       int64
	ResourceType *string // "tournament" | "match" | "team" | nil
	ResourceID   *int64
	AssignedBy   *int64
}

func (q *Queries) AssignUserRole(ctx context.Context, arg AssignUserRoleParams) (*models.UserRoleAssignment, error) {
	row := q.db.QueryRowContext(ctx, assignUserRole,
		arg.UserID,
		arg.RoleID,
		arg.ResourceType,
		arg.ResourceID,
		arg.AssignedBy,
	)
	var i models.UserRoleAssignment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RoleID,
		&i.ResourceType,
		&i.ResourceID,
		&i.AssignedBy,
		&i.IsActive,
		&i.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("AssignUserRole: failed to scan: %w", err)
	}
	return &i, nil
}

// Check Permission

const hasPermissionQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM user_role_assignments ura
		JOIN roles r ON r.id = ura.role_id
		WHERE ura.user_id = $1
		  AND r.name = $2
		  AND ura.is_active = TRUE
		  AND (
		      ura.resource_type IS NULL
		      OR (ura.resource_type = $3 AND ura.resource_id = $4)
		  )
	) AS has_permission
`

type HasPermissionParams struct {
	UserID       int64
	RoleName     string
	ResourceType *string
	ResourceID   *int64
}

func (q *Queries) HasRolePermission(ctx context.Context, arg HasPermissionParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, hasPermissionQuery,
		arg.UserID,
		arg.RoleName,
		arg.ResourceType,
		arg.ResourceID,
	)
	var hasPermission bool
	err := row.Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("HasRolePermission: failed to scan: %w", err)
	}
	return hasPermission, nil
}

// Get User Roles

const getUserRoles = `
	SELECT
		ura.id,
		ura.user_id,
		r.name   AS role_name,
		r.scope  AS role_scope,
		ura.resource_type,
		ura.resource_id,
		ura.assigned_by,
		ura.is_active,
		ura.created_at
	FROM user_role_assignments ura
	JOIN roles r ON r.id = ura.role_id
	WHERE ura.user_id = $1
	  AND ura.is_active = TRUE
	ORDER BY ura.created_at DESC
`

func (q *Queries) GetUserRoles(ctx context.Context, userID int64) ([]*models.UserRoleAssignmentWithRole, error) {
	rows, err := q.db.QueryContext(ctx, getUserRoles, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserRoles: query failed: %w", err)
	}
	defer rows.Close()

	var result []*models.UserRoleAssignmentWithRole
	for rows.Next() {
		var i models.UserRoleAssignmentWithRole
		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RoleName,
			&i.RoleScope,
			&i.ResourceType,
			&i.ResourceID,
			&i.AssignedBy,
			&i.IsActive,
			&i.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetUserRoles: failed to scan row: %w", err)
		}
		result = append(result, &i)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUserRoles: rows error: %w", err)
	}
	return result, nil
}

// Get Roles for a specific Resource

const getResourceUserRoles = `
	SELECT
		ura.id,
		ura.user_id,
		r.name   AS role_name,
		r.scope  AS role_scope,
		ura.resource_type,
		ura.resource_id,
		ura.assigned_by,
		ura.is_active,
		ura.created_at
	FROM user_role_assignments ura
	JOIN roles r ON r.id = ura.role_id
	WHERE ura.resource_type = $1
	  AND ura.resource_id   = $2
	  AND ura.is_active = TRUE
`

func (q *Queries) GetResourceUserRoles(ctx context.Context, resourceType string, resourceID int64) ([]*models.UserRoleAssignmentWithRole, error) {
	rows, err := q.db.QueryContext(ctx, getResourceUserRoles, resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("GetResourceUserRoles: query failed: %w", err)
	}
	defer rows.Close()

	var result []*models.UserRoleAssignmentWithRole
	for rows.Next() {
		var i models.UserRoleAssignmentWithRole
		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RoleName,
			&i.RoleScope,
			&i.ResourceType,
			&i.ResourceID,
			&i.AssignedBy,
			&i.IsActive,
			&i.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetResourceUserRoles: failed to scan row: %w", err)
		}
		result = append(result, &i)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetResourceUserRoles: rows error: %w", err)
	}
	return result, nil
}

// Revoke Role

const revokeUserRole = `
	UPDATE user_role_assignments
	SET is_active = FALSE
	WHERE user_id     = $1
	  AND role_id     = $2
	  AND (resource_type = $3 OR (resource_type IS NULL AND $3 IS NULL))
	  AND (resource_id   = $4 OR (resource_id   IS NULL AND $4 IS NULL))
`

type RevokeUserRoleParams struct {
	UserID       int64
	RoleID       int64
	ResourceType *string
	ResourceID   *int64
}

func (q *Queries) RevokeUserRole(ctx context.Context, arg RevokeUserRoleParams) error {
	_, err := q.db.ExecContext(ctx, revokeUserRole,
		arg.UserID,
		arg.RoleID,
		arg.ResourceType,
		arg.ResourceID,
	)
	if err != nil {
		return fmt.Errorf("RevokeUserRole: failed: %w", err)
	}
	return nil
}

// Get Role by Name

const getRoleByName = `
	SELECT id, name, scope, description, created_at
	FROM roles
	WHERE name = $1
`

func (q *Queries) GetRoleByName(ctx context.Context, name string) (*models.Roles, error) {
	row := q.db.QueryRowContext(ctx, getRoleByName, name)
	var r models.Roles
	err := row.Scan(&r.ID, &r.Name, &r.Scope, &r.Description, &r.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetRoleByName: failed to scan: %w", err)
	}
	return &r, nil
}

// Get All Roles

const getAllRoles = `
	SELECT id, name, scope, description, created_at
	FROM roles
	ORDER BY scope, name
`

func (q *Queries) GetAllRoles(ctx context.Context) ([]*models.Roles, error) {
	rows, err := q.db.QueryContext(ctx, getAllRoles)
	if err != nil {
		return nil, fmt.Errorf("GetAllRoles: query failed: %w", err)
	}
	defer rows.Close()

	var result []*models.Roles
	for rows.Next() {
		var r models.Roles
		err := rows.Scan(&r.ID, &r.Name, &r.Scope, &r.Description, &r.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("GetAllRoles: failed to scan: %w", err)
		}
		result = append(result, &r)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllRoles: rows error: %w", err)
	}
	return result, nil
}
