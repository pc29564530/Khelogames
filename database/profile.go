package database

import (
	"context"
	"khelogames/database/models"
)

const createProfile = `
INSERT INTO profile (
    owner,
    full_name,
    bio,
    avatar_url,
    created_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP
) RETURNING id, owner, full_name, bio, avatar_url, created_at
`

type CreateProfileParams struct {
	Owner     string `json:"owner"`
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

func (q *Queries) CreateProfile(ctx context.Context, arg CreateProfileParams) (models.Profile, error) {
	row := q.db.QueryRowContext(ctx, createProfile,
		arg.Owner,
		arg.FullName,
		arg.Bio,
		arg.AvatarUrl,
	)
	var profile models.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Owner,
		&profile.FullName,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.CreatedAt,
	)
	return profile, err
}

const editProfile = `
UPDATE profile
SET full_name=$1, avatar_url=$2, bio=$3
WHERE id=$4
RETURNING id, owner, full_name, bio, avatar_url, created_at
`

type EditProfileParams struct {
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
	Bio       string `json:"bio"`
	ID        int64  `json:"id"`
}

func (q *Queries) EditProfile(ctx context.Context, arg EditProfileParams) (models.Profile, error) {
	row := q.db.QueryRowContext(ctx, editProfile,
		arg.FullName,
		arg.AvatarUrl,
		arg.Bio,
		arg.ID,
	)
	var profile models.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Owner,
		&profile.FullName,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.CreatedAt,
	)
	return profile, err
}

const getProfile = `
SELECT id, owner, full_name, bio, avatar_url, created_at FROM profile
WHERE owner=$1
`

func (q *Queries) GetProfile(ctx context.Context, owner string) (models.Profile, error) {
	row := q.db.QueryRowContext(ctx, getProfile, owner)
	var profile models.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Owner,
		&profile.FullName,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.CreatedAt,
	)
	return profile, err
}

const updateAvatar = `
UPDATE profile
SET avatar_url=$1
WHERE owner=$2
RETURNING id, owner, full_name, bio, avatar_url, created_at
`

type UpdateAvatarParams struct {
	AvatarUrl string `json:"avatar_url"`
	Owner     string `json:"owner"`
}

func (q *Queries) UpdateAvatar(ctx context.Context, arg UpdateAvatarParams) (models.Profile, error) {
	row := q.db.QueryRowContext(ctx, updateAvatar, arg.AvatarUrl, arg.Owner)
	var profile models.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Owner,
		&profile.FullName,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.CreatedAt,
	)
	return profile, err
}

const updateBio = `
UPDATE profile
SET bio=$1
WHERE owner=$2
RETURNING id, owner, full_name, bio, avatar_url, created_at
`

type UpdateBioParams struct {
	Bio   string `json:"bio"`
	Owner string `json:"owner"`
}

func (q *Queries) UpdateBio(ctx context.Context, arg UpdateBioParams) (models.Profile, error) {
	row := q.db.QueryRowContext(ctx, updateBio, arg.Bio, arg.Owner)
	var profile models.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Owner,
		&profile.FullName,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.CreatedAt,
	)
	return profile, err
}

const updateFullName = `
UPDATE profile
SET full_name=$1
WHERE owner=$2
RETURNING id, owner, full_name, bio, avatar_url, created_at
`

func (q *Queries) UpdateFullName(ctx context.Context, fullName string, owner string) (models.Profile, error) {
	row := q.db.QueryRowContext(ctx, updateFullName, fullName, owner)
	var profile models.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Owner,
		&profile.FullName,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.CreatedAt,
	)
	return profile, err
}
