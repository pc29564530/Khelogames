package database

import (
	"context"
	"database/sql"
	"fmt"
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

func (q *Queries) UpdateFullName(ctx context.Context, owner string, fullName string) (models.Profile, error) {
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

const getRoles = `
	SELECT * FROM roles;
`

func (q *Queries) GetRoles(ctx context.Context) ([]models.Roles, error) {
	rows, err := q.db.QueryContext(ctx, getRoles)
	if err != nil {
		return nil, fmt.Errorf("Failed to query: ", err)
	}
	var roles []models.Roles

	for rows.Next() {
		var row models.Roles
		err := rows.Scan(
			&row.ID,
			&row.Name,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed to scan the row: ", err)
		}
		roles = append(roles, row)
	}

	return roles, err
}

const addRole = `
	INSERT INTO user_roles (
		profile_id,
		role_id,
		created_at
	) VALUES ($1, $2, CURRENT_TIMESTAMP) RETURNING *;
`

func (q *Queries) AddRole(ctx context.Context, profileID, roleID int64) (models.UserRole, error) {
	row := q.db.QueryRowContext(ctx, addRole,
		profileID,
		roleID,
	)
	var userRole models.UserRole
	err := row.Scan(
		&userRole.ID,
		&userRole.ProfileID,
		&userRole.RoleID,
		&userRole.CreatedAT,
	)
	return userRole, err
}

const addOrganizerVerificationDetails = `
	INSERT INTO organizers (
		profile_id,
		organization_name,
		email,
		phone_number,
		is_verified,
		verified_at,
		created_at
	) VALUES ($1, $2, $3, $4, false, null, CURRENT_TIMESTAMP ) RETURNING *;
`

func (q *Queries) AddOrganizerVerificationDetails(ctx context.Context, profileID int64, organizationName string, email string, phoneNumber string) (*models.Organizations, error) {
	row := q.db.QueryRowContext(ctx, addOrganizerVerificationDetails, profileID, organizationName, email, phoneNumber)
	var organizationDetails models.Organizations
	err := row.Scan(
		&organizationDetails.ID,
		&organizationDetails.ProfileID,
		&organizationDetails.OrganizationName,
		&organizationDetails.Email,
		&organizationDetails.PhoneNumber,
		&organizationDetails.IsVerified,
		&organizationDetails.VerifiedAT,
		&organizationDetails.CreatedAT,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan the row: ", err)
	}
	return &organizationDetails, nil
}

const addDocumentVerificationDetails = `
	INSERT INTO document (
		organizer_id,
		document_type,
		file_path,
		submitted_at,
		status
	) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, 'pending' ) RETURNING *;
`

func (q *Queries) AddDocumentVerificationDetails(ctx context.Context, organizerID int64, documentType string, filePath string) (*models.Document, error) {
	row := q.db.QueryRowContext(ctx, addDocumentVerificationDetails, organizerID, documentType, filePath)
	var documentVerification models.Document
	err := row.Scan(
		&documentVerification.ID,
		&documentVerification.OrganizerID,
		&documentVerification.DocumentType,
		&documentVerification.FilePath,
		&documentVerification.SubmittedAT,
		&documentVerification.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan the row: ", err)
	}
	return &documentVerification, nil
}
