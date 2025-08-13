package transcation_setup

import "github.com/google/uuid"

// CreateEmailSignUpParams represents parameters for email signup
type CreateEmailSignUpParams struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}

// CreateProfileParams represents parameters for creating user profile
type CreateProfileParams struct {
	UserID    int32  `json:"user_id"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
}

// CreateCommunityParams represents parameters for creating the community
type CreateCommunityParams struct {
	UserPublicID  uuid.UUID `json:"user_public_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CommunityType string    `json:"community_type"`
	AvatarUrl     string    `json:"avatar_url"`
	CoverImageUrl string    `json:"cover_image_url"`
}
