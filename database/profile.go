package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"
	"time"

	"github.com/google/uuid"
)

const createProfile = `
INSERT INTO user_profiles (
    user_id,
    bio,
    avatar_url,
	location,
	created_at,
	updated_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
) RETURNING *
`

type CreateProfileParams struct {
	UserID    int32  `json:"user_id"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
	Location  string `json:"location"`
}

func (q *Queries) CreateProfile(ctx context.Context, arg CreateProfileParams) (models.UserProfiles, error) {
	row := q.db.QueryRowContext(ctx, createProfile,
		arg.UserID,
		arg.Bio,
		arg.AvatarUrl,
		arg.Location,
	)
	var profile models.UserProfiles
	err := row.Scan(
		&profile.ID,
		&profile.PublicID,
		&profile.UserID,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.Location,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	return profile, err
}

type UpdateUserParams struct {
	PublicID uuid.UUID `db:"public_id"`
	FullName string    `db:"full_name"`
}

// edit profile
const editProfile = `
UPDATE user_profiles
SET avatar_url=$2, bio=$3
WHERE public_id=$1
RETURNING *
`

type EditProfileParams struct {
	PublicID  uuid.UUID `json:"public_id"`
	AvatarUrl string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
}

func (q *Queries) EditProfile(ctx context.Context, arg EditProfileParams) (models.UserProfiles, error) {
	row := q.db.QueryRowContext(ctx, editProfile,
		arg.PublicID,
		arg.AvatarUrl,
		arg.Bio,
	)
	var profile models.UserProfiles
	err := row.Scan(
		&profile.ID,
		&profile.PublicID,
		&profile.UserID,
		&profile.Bio,
		&profile.AvatarUrl,
		&profile.Location,
	)
	return profile, err
}

const getProfile = `
SELECT up.*, u.full_name AS full_name, u.username AS username FROM user_profiles up
LEFT JOIN users AS u ON u.id = up.user_id
WHERE up.public_id = $1
`

type userProfile struct {
	ID        int64     `json:"id"`
	PublicID  uuid.UUID `json:"public_id"`
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

func (q *Queries) GetProfile(ctx context.Context, publicID uuid.UUID) (*userProfile, error) {

	row := q.db.QueryRowContext(ctx, getProfile, publicID)
	var res userProfile
	err := row.Scan(
		&res.ID,
		&res.PublicID,
		&res.UserID,
		&res.Username,
		&res.FullName,
		&res.Bio,
		&res.AvatarUrl,
		&res.CreatedAT,
		&res.UpdatedAT,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, err
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
    WITH userID AS (
		SELECT id FROM users
		WHERE public_id = $1
	)
	INSERT INTO user_roles (
		user_id,
		role_id,
		created_at
	) 
	SELECT 
		userID.id,
		$2,
		CURRENT_TIMESTAMP
	FROM userID
	RETURNING *;
`

func (q *Queries) AddRole(ctx context.Context, userID uuid.UUID, roleID int32) (models.UserRole, error) {
	row := q.db.QueryRowContext(ctx, addRole, userID, roleID)
	var userRole models.UserRole
	err := row.Scan(
		&userRole.ID,
		&userRole.UserID,
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
