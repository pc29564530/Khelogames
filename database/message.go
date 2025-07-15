package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const createNewMessage = `
WITH senderID AS (
	SELECT user_id FROM users
	WHERE public_id = $1;
),
receiverID AS (
	SELECT user_id FROM users
	WHERE public_id = $2;
)
INSERT INTO message (
    sender_id,
    receiver_id,
	content,
    media_url,
    media_type,
    CURRENT_TIMESTAMP,
)
SELECT 
	senderID.id,
    receiverID.id,
	$3,
    $4,
    $5,
    CURRENT_TIMESTAMP,
FROM senderID, receiverID
RETURNING *;
`

type CreateNewMessageParams struct {
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	MediaUrl   string    `json:"media_url"`
	MediaType  string    `json:"media_type"`
}

func (q *Queries) CreateNewMessage(ctx context.Context, arg CreateNewMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, createNewMessage,
		arg.SenderID,
		arg.ReceiverID,
		arg.Content,
		arg.MediaUrl,
		arg.MediaType,
	)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.SenderID,
		&i.ReceiverID,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsSeen,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

const deleteMessage = `
DELETE FROM message
WHERE sender_id=$1 and receiver_id=$2
RETURNING *
`

type DeleteMessageParams struct {
	SenderID string `json:"sender_id"`
	ID       int64  `json:"id"`
}

func (q *Queries) DeleteMessage(ctx context.Context, arg DeleteMessageParams) (models.Message, error) {
	row := q.db.QueryRowContext(ctx, deleteMessage, arg.SenderID, arg.ID)
	var i models.Message
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.SenderID,
		&i.ReceiverID,
		&i.Content,
		&i.MediaUrl,
		&i.MediaType,
		&i.IsSeen,
		&i.IsDeleted,
		&i.CreatedAt,
	)
	return i, err
}

// fetch all the message by the receiver id
const getMessageByReceiver = `
SELECT m.*
FROM message m
JOIN users sender ON sender.id = m.sender_id
JOIN users receiver ON receiver.id = m.receiver_id
WHERE (sender.public_id = $1 AND receiver.public_id = $2)
   OR (sender.public_id = $2 AND receiver.public_id = $1)
ORDER BY m.id ASC;
`

type GetMessageByReceiverParams struct {
	ReceiverID uuid.UUID `json:"receiver_id"`
	SenderID   uuid.UUID `json:"sender_id"`
}

func (q *Queries) GetMessageByReceiver(ctx context.Context, arg GetMessageByReceiverParams) ([]models.Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessageByReceiver, arg.SenderID, arg.ReceiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Message
	for rows.Next() {
		var i models.Message
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.SenderID,
			&i.ReceiverID,
			&i.Content,
			&i.MediaUrl,
			&i.MediaType,
			&i.IsSeen,
			&i.IsDeleted,
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

// get users by message sender
const getUserByMessageSender = `
SELECT DISTINCT receiver_id
FROM message m
JOIN users AS u ON u.id = m.user_id
WHERE u.public_id = $1
`

func (q *Queries) GetUserByMessageSend(ctx context.Context, senderID uuid.UUID) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getUserByMessageSender, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var receiverID int32
		if err := rows.Scan(&receiverID); err != nil {
			return nil, err
		}
		items = append(items, receiverID)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// schedule delete message
const scheduledDeleteMessage = `
DELETE FROM message
WHERE is_deleted = TRUE AND deleted_at < NOW() - INTERVAL '30 days'
RETURNING *
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
			&i.PublicID,
			&i.SenderID,
			&i.ReceiverID,
			&i.Content,
			&i.MediaUrl,
			&i.MediaType,
			&i.IsSeen,
			&i.IsDeleted,
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
