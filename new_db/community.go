package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const createCommunity = `
INSERT INTO communities (
    owner,
    communities_name,
    description,
    community_type
) VALUES (
    $1, $2, $3, $4
) RETURNING id, owner, communities_name, description, community_type, created_at
`

type CreateCommunityParams struct {
	Owner           string `json:"owner"`
	CommunitiesName string `json:"communities_name"`
	Description     string `json:"description"`
	CommunityType   string `json:"community_type"`
}

func (q *Queries) CreateCommunity(ctx context.Context, arg CreateCommunityParams) (models.Community, error) {
	row := q.db.QueryRowContext(ctx, createCommunity,
		arg.Owner,
		arg.CommunitiesName,
		arg.Description,
		arg.CommunityType,
	)
	var i models.Community
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.CommunitiesName,
		&i.Description,
		&i.CommunityType,
		&i.CreatedAt,
	)
	return i, err
}

const getAllCommunities = `
SELECT id, owner, communities_name, description, community_type, created_at FROM communities
ORDER BY id
`

func (q *Queries) GetAllCommunities(ctx context.Context) ([]models.Community, error) {
	rows, err := q.db.QueryContext(ctx, getAllCommunities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Community
	for rows.Next() {
		var i models.Community
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.CommunitiesName,
			&i.Description,
			&i.CommunityType,
			&i.CreatedAt,
		); err != nil {
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

const getCommunitiesMember = `
SELECT users.username FROM users
JOIN communities ON users.username = communities.owner
WHERE communities.communities_name=$1
`

func (q *Queries) GetCommunitiesMember(ctx context.Context, communitiesName string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getCommunitiesMember, communitiesName)
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

const getCommunity = `
SELECT id, owner, communities_name, description, community_type, created_at FROM communities
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCommunity(ctx context.Context, id int64) (models.Community, error) {
	row := q.db.QueryRowContext(ctx, getCommunity, id)
	var i models.Community
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.CommunitiesName,
		&i.Description,
		&i.CommunityType,
		&i.CreatedAt,
	)
	return i, err
}

const getCommunityByCommunityName = `
SELECT id, owner, communities_name, description, community_type, created_at FROM communities
WHERE communities_name=$1
`

func (q *Queries) GetCommunityByCommunityName(ctx context.Context, communitiesName string) (models.Community, error) {
	row := q.db.QueryRowContext(ctx, getCommunityByCommunityName, communitiesName)
	var i models.Community
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.CommunitiesName,
		&i.Description,
		&i.CommunityType,
		&i.CreatedAt,
	)
	return i, err
}

const updateCommunityDescription = `
UPDATE communities
SET description=$1
WHERE id=$2
RETURNING id, owner, communities_name, description, community_type, created_at
`

type UpdateCommunityDescriptionParams struct {
	Description string `json:"description"`
	ID          int64  `json:"id"`
}

func (q *Queries) UpdateCommunityDescription(ctx context.Context, arg UpdateCommunityDescriptionParams) (models.Community, error) {
	row := q.db.QueryRowContext(ctx, updateCommunityDescription, arg.Description, arg.ID)
	var i models.Community
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.CommunitiesName,
		&i.Description,
		&i.CommunityType,
		&i.CreatedAt,
	)
	return i, err
}

const updateCommunityName = `-- name: UpdateCommunityName :one
UPDATE communities
SET communities_name=$1
WHERE id=$2
RETURNING id, owner, communities_name, description, community_type, created_at
`

type UpdateCommunityNameParams struct {
	CommunitiesName string `json:"communities_name"`
	ID              int64  `json:"id"`
}

func (q *Queries) UpdateCommunityName(ctx context.Context, arg UpdateCommunityNameParams) (models.Community, error) {
	row := q.db.QueryRowContext(ctx, updateCommunityName, arg.CommunitiesName, arg.ID)
	var i models.Community
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.CommunitiesName,
		&i.Description,
		&i.CommunityType,
		&i.CreatedAt,
	)
	return i, err
}
