package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const createUserConnectionQuery = `
WITH userID AS (
	SELECT * FROM users WHERE public_id = $1
),
targetUserID AS (
	SELECT * FROM user_profiles WHERE public_id = $2
)
INSERT INTO users_connections (
    user_id,
    target_user_id
) 
SELECT
	userID.id,
	targetUserID.user_id
FROM userID, targetUserID
RETURNING *;
`

// CreateUserConnections
func (q *Queries) CreateUserConnections(ctx context.Context, userPublicID, targetPublicID uuid.UUID) (*models.UsersConnections, error) {
	row := q.db.QueryRowContext(ctx, createUserConnectionQuery, userPublicID, targetPublicID)
	var i models.UsersConnections
	err := row.Scan(
		&i.UserID,
		&i.TargetUserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan user connection: %w", err)
	}
	return &i, err
}

const deleteUsersConnectionsQuery = `
DELETE FROM users_connections
JOIN users ON u.user_id = users.id
JOIN users ON tu.users_id = users.id
WHERE u.public_id = $1 AND tu.public_id = $2
`

func (q *Queries) DeleteUsersConnections(ctx context.Context, userID, targetUserID uuid.UUID) error {
	row := q.db.QueryRowContext(ctx, deleteUsersConnectionsQuery, userID, targetUserID)
	var i models.UsersConnections
	err := row.Scan(
		&i.UserID,
		&i.TargetUserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("Failed to get users connection: %w", err)
	}
	return err
}

const getAllFollower = `
SELECT 
    JSON_BUILD_OBJECT(
        'user_public_id', tu.public_id,
        'profile', JSON_BUILD_OBJECT(
            'id', up.id,
            'public_id', up.public_id,
            'username', tu.username,
            'full_name', tu.full_name,
            'avatar_url', up.avatar_url
        )
    )
FROM users_connections uc
JOIN users u ON u.id = uc.target_user_id
JOIN users tu ON tu.id = uc.user_id
JOIN user_profiles up ON up.user_id = tu.id
WHERE u.public_id = $1;
`

func (q *Queries) GetAllFollower(ctx context.Context, targetPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getAllFollower, targetPublicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		var item map[string]interface{}
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}
		err := json.Unmarshal(jsonByte, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllFollowing = `
SELECT 
    JSON_BUILD_OBJECT(
        'user_public_id', tu.public_id,
        'profile', JSON_BUILD_OBJECT(
            'id', up.id,
            'public_id', up.public_id,
            'username', tu.username,
            'full_name', tu.full_name,
            'avatar_url', up.avatar_url
        )
    )
FROM users_connections uc
JOIN users u ON u.id = uc.user_id
JOIN users tu ON tu.id = uc.target_user_id
JOIN user_profiles up ON up.user_id = tu.id
WHERE u.public_id = $1;
`

func (q *Queries) GetAllFollowing(ctx context.Context, userPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getAllFollowing, userPublicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		var item map[string]interface{}
		if err := rows.Scan(&jsonByte); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		err := json.Unmarshal(jsonByte, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isFollowingCheck = `
SELECT COUNT(*) > 0
FROM users_connections uc
JOIN users follower ON follower.id = uc.user_id
JOIN user_profiles following ON following.user_id = uc.target_user_id
WHERE follower.public_id = $1
  AND following.public_id = $2;
`

func (q *Queries) IsFollowingF(ctx context.Context, followerPublicID, followingPublicID uuid.UUID) (bool, error) {
	var isFollowingUser bool
	err := q.db.QueryRowContext(ctx, isFollowingCheck, followerPublicID, followingPublicID).Scan(&isFollowingUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return isFollowingUser, nil
}
