package new_db

import (
	"context"
	"khelogames/new_db/models"
	"time"
)

const createNewMessage = `
INSERT INTO message (
    content,
    is_seen,
    sender_username,
    receiver_username,
    media_url,
    media_type,
    sent_at,
    is_deleted,
    deleted_at
) VALUES (
    $1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP, $7, $8
) RETURNING id, content, is_seen, sender_username, receiver_username, sent_at, media_url, media_type, is_deleted, deleted_at
`

type CreateNewMessageParams struct {
	Content          string    `json:"content"`
	IsSeen           bool      `json:"is_seen"`
	SenderUsername   string    `json:"sender_username"`
	ReceiverUsername string    `json:"receiver_username"`
	MediaUrl         string    `json:"media_url"`
	MediaType        string    `json:"media_type"`
	IsDeleted        bool      `json:"is_deleted"`
	DeletedAt        time.Time `json:"deleted_at"`
}

func (q *Queries) CreateNewMessage(ctx context.Context, arg CreateNewMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, createNewMessage,
		arg.Content,
		arg.IsSeen,
		arg.SenderUsername,
		arg.ReceiverUsername,
		arg.MediaUrl,
		arg.MediaType,
		arg.IsDeleted,
		arg.DeletedAt,
	)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.IsSeen,
		&i.SenderUsername,
		&i.ReceiverUsername,
		&i.SentAt,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsDeleted,
		&i.DeletedAt,
	)
	return i, err
}

const deleteMessage = `
DELETE FROM message
WHERE sender_username=$1 and id=$2
RETURNING id, content, is_seen, sender_username, receiver_username, sent_at, media_url, media_type, is_deleted, deleted_at
`

type DeleteMessageParams struct {
	SenderUsername string `json:"sender_username"`
	ID             int64  `json:"id"`
}

func (q *Queries) DeleteMessage(ctx context.Context, arg DeleteMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, deleteMessage, arg.SenderUsername, arg.ID)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.IsSeen,
		&i.SenderUsername,
		&i.ReceiverUsername,
		&i.SentAt,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsDeleted,
		&i.DeletedAt,
	)
	return i, err
}

const getMessageByReceiver = `
SELECT id, content, is_seen, sender_username, receiver_username, sent_at, media_url, media_type, is_deleted, deleted_at FROM message
WHERE (sender_username=$1 AND receiver_username=$2) OR (receiver_username=$1 AND sender_username=$2)
ORDER BY id ASC
`

type GetMessageByReceiverParams struct {
	SenderUsername   string `json:"sender_username"`
	ReceiverUsername string `json:"receiver_username"`
}

func (q *Queries) GetMessageByReceiver(ctx context.Context, arg GetMessageByReceiverParams) ([]models.Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessageByReceiver, arg.SenderUsername, arg.ReceiverUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Message
	for rows.Next() {
		var i models.Message
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.IsSeen,
			&i.SenderUsername,
			&i.ReceiverUsername,
			&i.SentAt,
			&i.MediaUrl,
			&i.MediaType,
			&i.IsDeleted,
			&i.DeletedAt,
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

const getUserByMessageSend = `
SELECT DISTINCT receiver_username
FROM message
WHERE sender_username = $1
`

func (q *Queries) GetUserByMessageSend(ctx context.Context, senderUsername string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getUserByMessageSend, senderUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var receiver_username string
		if err := rows.Scan(&receiver_username); err != nil {
			return nil, err
		}
		items = append(items, receiver_username)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const scheduledDeleteMessage = `
DELETE FROM message
WHERE is_deleted = TRUE AND deleted_at < NOW() - INTERVAL '30 days'
RETURNING id, content, is_seen, sender_username, receiver_username, sent_at, media_url, media_type, is_deleted, deleted_at
`

func (q *Queries) ScheduledDeleteMessage(ctx context.Context) ([]models.Message, error) {
	rows, err := q.db.QueryContext(ctx, scheduledDeleteMessage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Message
	for rows.Next() {
		var i models.Message
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.IsSeen,
			&i.SenderUsername,
			&i.ReceiverUsername,
			&i.SentAt,
			&i.MediaUrl,
			&i.MediaType,
			&i.IsDeleted,
			&i.DeletedAt,
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

const updateDeletedMessage = `-- name: UpdateDeletedMessage :one
UPDATE message
SET is_deleted=true AND deleted_at=NOW()
WHERE sender_username=$1 and id=$2
RETURNING id, content, is_seen, sender_username, receiver_username, sent_at, media_url, media_type, is_deleted, deleted_at
`

type UpdateDeletedMessageParams struct {
	SenderUsername string `json:"sender_username"`
	ID             int64  `json:"id"`
}

func (q *Queries) UpdateDeletedMessage(ctx context.Context, arg UpdateDeletedMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, updateDeletedMessage, arg.SenderUsername, arg.ID)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.IsSeen,
		&i.SenderUsername,
		&i.ReceiverUsername,
		&i.SentAt,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsDeleted,
		&i.DeletedAt,
	)
	return i, err
}

const updateSoftDeleteMessage = `
UPDATE message
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1 AND sender_username = $2
RETURNING id, content, is_seen, sender_username, receiver_username, sent_at, media_url, media_type, is_deleted, deleted_at
`

type UpdateSoftDeleteMessageParams struct {
	ID             int64  `json:"id"`
	SenderUsername string `json:"sender_username"`
}

func (q *Queries) UpdateSoftDeleteMessage(ctx context.Context, arg UpdateSoftDeleteMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, updateSoftDeleteMessage, arg.ID, arg.SenderUsername)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.IsSeen,
		&i.SenderUsername,
		&i.ReceiverUsername,
		&i.SentAt,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsDeleted,
		&i.DeletedAt,
	)
	return i, err
}
