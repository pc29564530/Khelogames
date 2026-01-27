package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
	is_active,
	avatar_url,
	cover_image_url,
	created_at,
	updated_at
) 
SELECT
	userID.id,
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
	AvatarUrl     string    `json:"avatar_url"`
	CoverImageUrl string    `json:"cover_image_url"`
}

func (q *Queries) CreateCommunity(ctx context.Context, arg CreateCommunityParams) (models.Communities, error) {
	row := q.db.QueryRowContext(ctx, createCommunity,
		arg.UserPublicID,
		arg.Name,
		arg.Slug,
		arg.Description,
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
			if err == sql.ErrNoRows {
				return nil, nil
			}
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
	JSON_BUILD_OBJECT('id', p.id, 'public_id',p.public_id, 'user_id',p.user_id,  'username',u.username,  'full_name',u.full_name,  'bio',p.bio,  'avatar_url',p.avatar_url,  'created_at',p.created_at )
FROM join_community jc
JOIN communities c ON jc.community_id = c.id
JOIN users u ON u.id = jc.user_id
JOIN user_profiles p ON p.user_id = jc.user_id
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

func (q *Queries) GetCommunitiesMember(ctx context.Context, communityPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getCommunitiesMember, communityPublicID)
	if err != nil {
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
			return nil, fmt.Errorf("Failed to unmarshal: ", err)
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

func (q *Queries) GetCommunity(ctx context.Context, publicID uuid.UUID) (*models.Communities, error) {
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
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Unable to get community: ", err)
	}
	return &i, nil
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

// increment community count
func (q *Queries) IncrementCommunityMemberCount(ctx context.Context, communityPublicID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, `
		UPDATE communities
		SET member_count = member_count + 1,
		    updated_at = NOW()
		WHERE public_id = $1
	`, communityPublicID)
	return err
}
