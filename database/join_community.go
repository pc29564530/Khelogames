package database

import (
	"context"
	"khelogames/database/models"
)

const addJoinCommunity = `
INSERT INTO join_community (
    community_name,
    username
) VALUES (
    $1, $2
) RETURNING id, community_name, username
`

type AddJoinCommunityParams struct {
	CommunityName string `json:"community_name"`
	Username      string `json:"username"`
}

func (q *Queries) AddJoinCommunity(ctx context.Context, arg AddJoinCommunityParams) (models.JoinCommunity, error) {
	row := q.db.QueryRowContext(ctx, addJoinCommunity, arg.CommunityName, arg.Username)
	var i models.JoinCommunity
	err := row.Scan(&i.ID, &i.CommunityName, &i.Username)
	return i, err
}

const getCommunityByUser = `
SELECT id, community_name, username FROM join_community
WHERE username=$1
ORDER BY id
`

func (q *Queries) GetCommunityByUser(ctx context.Context, username string) ([]models.JoinCommunity, error) {
	rows, err := q.db.QueryContext(ctx, getCommunityByUser, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.JoinCommunity
	for rows.Next() {
		var i models.JoinCommunity
		if err := rows.Scan(&i.ID, &i.CommunityName, &i.Username); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByCommunity = `
SELECT DISTINCT username FROM join_community
WHERE community_name=$1
`

func (q *Queries) GetUserByCommunity(ctx context.Context, communityName string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getUserByCommunity, communityName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		items = append(items, username)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const inActiveUserFromCommunity = `
UPDATE join_community
SET is_active = FALSE
WHERE username = $1 AND id = $2
RETURNING id, community_name, username
`

type InActiveUserFromCommunityParams struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

func (q *Queries) InActiveUserFromCommunity(ctx context.Context, arg InActiveUserFromCommunityParams) (models.JoinCommunity, error) {
	row := q.db.QueryRowContext(ctx, inActiveUserFromCommunity, arg.Username, arg.ID)
	var i models.JoinCommunity
	err := row.Scan(&i.ID, &i.CommunityName, &i.Username)
	return i, err
}

const removeUserFromCommunity = `
DELETE FROM join_community
WHERE id=$1 AND username=$2
RETURNING id, community_name, username
`

type RemoveUserFromCommunityParams struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (q *Queries) RemoveUserFromCommunity(ctx context.Context, arg RemoveUserFromCommunityParams) (models.JoinCommunity, error) {
	row := q.db.QueryRowContext(ctx, removeUserFromCommunity, arg.ID, arg.Username)
	var i models.JoinCommunity
	err := row.Scan(&i.ID, &i.CommunityName, &i.Username)
	return i, err
}
