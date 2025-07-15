package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const createCommunity = `
WITH userID AS (
	SELECT * FROM users
	WHERE public_id=$1
)
INSERT INTO communities (
    user_id,
    name,
	slug,
    description,
    community_type,
	is_active,
	avatar_url,
	cover_image_url,
	created_at,
	updated_at
) 
SELECT
	userID.id,
	$1,
	$2,
	$3,
	$4,
	true,
	$5,
	$6,
	CURRENT_TIMESTAMP,
	CURRENT_TIMESTAMP
FROM userID
RETURNING *
`

type CreateCommunityParams struct {
	UserPublicID  uuid.UUID `json:"user_public_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CommunityType string    `json:"community_type"`
	AvatarUrl     string    `json:"avatar_url"`
	CoverImageUrl string    `json:"cover_image_url"`
}

func (q *Queries) CreateCommunity(ctx context.Context, arg CreateCommunityParams) (models.Communities, error) {
	row := q.db.QueryRowContext(ctx, createCommunity,
		arg.UserPublicID,
		arg.Name,
		arg.Slug,
		arg.Description,
		arg.CommunityType,
		arg.AvatarUrl,
		arg.CoverImageUrl,
	)
	var i models.Communities
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.CommunityType,
		&i.IsActive,
		&i.MemberCount,
		&i.AvatarUrl,
		&i.CoverImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllCommunities = `
SELECT * FROM communities
ORDER BY id
`

func (q *Queries) GetAllCommunities(ctx context.Context) ([]models.Communities, error) {
	rows, err := q.db.QueryContext(ctx, getAllCommunities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Communities
	for rows.Next() {
		var i models.Communities
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.Name,
			&i.Slug,
			&i.Description,
			&i.CommunityType,
			&i.IsActive,
			&i.MemberCount,
			&i.AvatarUrl,
			&i.CoverImageUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
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
SELECT
	u.id 
	up.public_id AS profile_public_id,
	u.public_id AS user_public_id,
	u.username AS username
	up.full_name,
	up.bio,
	up.avatar_url
FROM join_community jc
JOIN communities c ON jc.community_id = c.id
JOIN users_profile up ON up.user_id = jc.user_id
WHERE c.public_id = $1;
`

type CommunityMember struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	UserPublicID    uuid.UUID `json:"user_public_id"`
	ProfilePublicID uuid.UUID `json:"profile_public_id"`
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Bio             string    `json:"bio"`
	AvatarURL       string    `json:"avatar_url"`
}

func (q *Queries) GetCommunitiesMember(ctx context.Context, communityPublicID uuid.UUID) ([]CommunityMember, error) {
	rows, err := q.db.QueryContext(ctx, getCommunitiesMember, communityPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CommunityMember
	for rows.Next() {
		var item CommunityMember
		if err := rows.Scan(&item); err != nil {
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

const getCommunity = `
SELECT * FROM communities
WHERE public_id = $1
`

func (q *Queries) GetCommunity(ctx context.Context, publicID uuid.UUID) (models.Communities, error) {
	row := q.db.QueryRowContext(ctx, getCommunity, publicID)
	var i models.Communities
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.CommunityType,
		&i.IsActive,
		&i.MemberCount,
		&i.AvatarUrl,
		&i.CoverImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCommunityByCommunityName = `
SELECT * FROM communities
WHERE name=$1
`

func (q *Queries) GetCommunityByCommunityName(ctx context.Context, name string) (models.Communities, error) {
	row := q.db.QueryRowContext(ctx, getCommunityByCommunityName, name)
	var i models.Communities
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.CommunityType,
		&i.IsActive,
		&i.MemberCount,
		&i.AvatarUrl,
		&i.CoverImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCommunityDescription = `
	UPDATE communities
	SET description=$2
	WHERE publicID=$1
	RETURNING *
`

func (q *Queries) UpdateCommunityDescription(ctx context.Context, publicID uuid.UUID, description string) (models.Communities, error) {
	row := q.db.QueryRowContext(ctx, updateCommunityDescription, publicID, description)
	var i models.Communities
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.CommunityType,
		&i.IsActive,
		&i.MemberCount,
		&i.AvatarUrl,
		&i.CoverImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCommunityName = `
	UPDATE communities
	SET name=$2
	WHERE public_id=$1
	RETURNING *
`

func (q *Queries) UpdateCommunityName(ctx context.Context, publicID uuid.UUID, name string) (models.Communities, error) {
	row := q.db.QueryRowContext(ctx, updateCommunityName, publicID, name)
	var i models.Communities
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.CommunityType,
		&i.IsActive,
		&i.MemberCount,
		&i.AvatarUrl,
		&i.CoverImageUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
