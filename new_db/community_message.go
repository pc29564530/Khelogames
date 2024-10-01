package new_db

import (
	"context"
	"khelogames/new_db/models"
	"time"
)

const createCommunityMessage = `
INSERT INTO communitymessage(
    community_name,
    sender_username,
    content,
    sent_at
) VALUES ($1,$2, $3, CURRENT_TIMESTAMP )
RETURNING id, community_name, sender_username, content, sent_at
`

type CreateCommunityMessageParams struct {
	CommunityName  string `json:"community_name"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
}

func (q *Queries) CreateCommunityMessage(ctx context.Context, arg CreateCommunityMessageParams) (models.Communitymessage, error) {
	row := q.db.QueryRowContext(ctx, createCommunityMessage, arg.CommunityName, arg.SenderUsername, arg.Content)
	var i models.Communitymessage
	err := row.Scan(
		&i.ID,
		&i.CommunityName,
		&i.SenderUsername,
		&i.Content,
		&i.SentAt,
	)
	return i, err
}

const createMessageMedia = `
INSERT INTO messagemedia (
    message_id,
    media_id
) VALUES (
    $1, $2
) RETURNING message_id, media_id
`

type CreateMessageMediaParams struct {
	MessageID int64 `json:"message_id"`
	MediaID   int64 `json:"media_id"`
}

func (q *Queries) CreateMessageMedia(ctx context.Context, arg CreateMessageMediaParams) (models.Messagemedium, error) {
	row := q.db.QueryRowContext(ctx, createMessageMedia, arg.MessageID, arg.MediaID)
	var i models.Messagemedium
	err := row.Scan(&i.MessageID, &i.MediaID)
	return i, err
}

const createUploadMedia = `
INSERT INTO uploadmedia (
    media_url,
    media_type,
    sent_at
) VALUES ($1, $2, CURRENT_TIMESTAMP)
RETURNING id, media_url, media_type, sent_at
`

type CreateUploadMediaParams struct {
	MediaUrl  string `json:"media_url"`
	MediaType string `json:"media_type"`
}

func (q *Queries) CreateUploadMedia(ctx context.Context, arg CreateUploadMediaParams) (models.Uploadmedium, error) {
	row := q.db.QueryRowContext(ctx, createUploadMedia, arg.MediaUrl, arg.MediaType)
	var i models.Uploadmedium
	err := row.Scan(
		&i.ID,
		&i.MediaUrl,
		&i.MediaType,
		&i.SentAt,
	)
	return i, err
}

const getCommunityByMessage = `
SELECT DISTINCT community_name FROM communitymessage
`

func (q *Queries) GetCommunityByMessage(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getCommunityByMessage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var community_name string
		if err := rows.Scan(&community_name); err != nil {
			return nil, err
		}
		items = append(items, community_name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCommuntiyMessage = `
SELECT cm.id, cm.community_name, cm.sender_username, cm.content, cm.sent_at, um.media_url, um.media_type 
FROM communitymessage cm 
JOIN messagemedia mm ON mm.message_id = cm.id 
JOIN uploadmedia um ON mm.media_id = um.id
`

type GetCommuntiyMessageRow struct {
	ID             int64     `json:"id"`
	CommunityName  string    `json:"community_name"`
	SenderUsername string    `json:"sender_username"`
	Content        string    `json:"content"`
	SentAt         time.Time `json:"sent_at"`
	MediaUrl       string    `json:"media_url"`
	MediaType      string    `json:"media_type"`
}

func (q *Queries) GetCommuntiyMessage(ctx context.Context) ([]GetCommuntiyMessageRow, error) {
	rows, err := q.db.QueryContext(ctx, getCommuntiyMessage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommuntiyMessageRow
	for rows.Next() {
		var i GetCommuntiyMessageRow
		if err := rows.Scan(
			&i.ID,
			&i.CommunityName,
			&i.SenderUsername,
			&i.Content,
			&i.SentAt,
			&i.MediaUrl,
			&i.MediaType,
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
