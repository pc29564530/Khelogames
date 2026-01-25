package database

import (
	"context"
	"database/sql"
	"khelogames/database/models"
	"time"

	"github.com/google/uuid"
)

type CommunityMessage struct {
	ID          int64     `json:"id"`
	PublicID    uuid.UUID `json:"public_id"`
	CommunityID int32     `json:"community_ud"`
	SenderID    int32     `json:"sender_id"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	MediaUrl    string    `json:"media_url"`
	MediaType   string    `json:"media_type"`
	SentAt      time.Time `json:"sent_at"`
}

const createCommunityMessage = `
WITH senderID AS (
	SELECT * FROM users WHERE public_id = $1
),
communityID AS (
	SELECT * FROM communities WHERE public_id = $2
)
INSERT INTO community_message(
    community_id,
	sender_id,
    name,
    content,
	media_url,
	media_type,
    sent_at
) 
SELECT 
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	CURRENT_TIMESTAMP
FROM senderID, communityID
RETURNING *;
`

type CreateCommunityMessageParams struct {
	CommuntiyPublicID uuid.UUID `json:"community_public_id"`
	SenderPublicID    uuid.UUID `json:"sender_public_id"`
	Name              string    `json:"name"`
	Content           string    `json:"content"`
	MediaUrl          string    `json:"media_url"`
	MediaType         string    `json:"media_type"`
}

func (q *Queries) CreateCommunityMessage(ctx context.Context, arg CreateCommunityMessageParams) (models.CommunityMessage, error) {
	row := q.db.QueryRowContext(ctx, createCommunityMessage, arg.CommuntiyPublicID, arg.SenderPublicID, arg.Name, arg.Content, arg.MediaUrl, arg.MediaType)
	var i models.CommunityMessage
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.CommunityID,
		&i.SenderID,
		&i.Name,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.SentAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.CommunityMessage{}, nil
		}
		return models.CommunityMessage{}, err
	}
	return i, nil
}

const getCommunityByMessage = `
SELECT c.*
FROM communities c
WHERE c.id IN (
    SELECT DISTINCT cm.community_id
    FROM community_message cm
    JOIN join_community jc ON jc.community_id = cm.community_id
    JOIN users u ON u.id = jc.user_id
    WHERE u.public_id = $1
);
`

func (q *Queries) GetCommunityByMessage(ctx context.Context, userPublicID uuid.UUID) ([]models.Communities, error) {
	rows, err := q.db.QueryContext(ctx, getCommunityByMessage, userPublicID)
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

const getCommuntiyMessages = `
	SELECT * FROM community_message cm
	JOIN communities AS c ON c.id = cm.community_id
	JOIN users AS u ON u.id = cm.user_id
	WHERE c.public_id = $1;
`

func (q *Queries) GetCommuntiyMessage(ctx context.Context, communityPublicID uuid.UUID) ([]models.CommunityMessage, error) {
	rows, err := q.db.QueryContext(ctx, getCommuntiyMessages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.CommunityMessage
	for rows.Next() {
		var i models.CommunityMessage
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.CommunityID,
			&i.SenderID,
			&i.Name,
			&i.Content,
			&i.MediaUrl,
			&i.MediaType,
			&i.SentAt,
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
