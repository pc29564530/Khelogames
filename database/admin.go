package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addAdmin = `
INSERT INTO content_admin (
    content_id,
    admin
) VALUES ( $1, $2 )
RETURNING id, content_id, admin
`

type AddAdminParams struct {
	ContentID int64     `json:"content_id"`
	Admin     uuid.UUID `json:"admin"`
}

func (q *Queries) AddAdmin(ctx context.Context, arg AddAdminParams) (models.ContentAdmin, error) {
	row := q.db.QueryRowContext(ctx, addAdmin, arg.ContentID, arg.Admin)
	var i models.ContentAdmin
	err := row.Scan(&i.ID, &i.ContentID, &i.Admin)
	return i, err
}

const deleteAdmin = `
DELETE FROM content_admin
WHERE content_id=$1 AND admin=$2
RETURNING id, content_id, admin
`

type DeleteAdminParams struct {
	ContentID int64     `json:"content_id"`
	Admin     uuid.UUID `json:"admin"`
}

func (q *Queries) DeleteAdmin(ctx context.Context, arg DeleteAdminParams) (models.ContentAdmin, error) {
	row := q.db.QueryRowContext(ctx, deleteAdmin, arg.ContentID, arg.Admin)
	var i models.ContentAdmin
	err := row.Scan(&i.ID, &i.ContentID, &i.Admin)
	return i, err
}

const getAdmin = `
SELECT id, content_id, admin FROM content_admin
WHERE content_id=$1
`

func (q *Queries) GetAdmin(ctx context.Context, contentID int64) ([]models.ContentAdmin, error) {
	rows, err := q.db.QueryContext(ctx, getAdmin, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.ContentAdmin
	for rows.Next() {
		var i models.ContentAdmin
		if err := rows.Scan(&i.ID, &i.ContentID, &i.Admin); err != nil {
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

const updateAdmin = `
UPDATE content_admin
SET admin=$1
WHERE content_id=$2 AND admin=$3
RETURNING id, content_id, admin
`

type UpdateAdminParams struct {
	Admin     string `json:"admin"`
	ContentID int64  `json:"content_id"`
	Admin_2   string `json:"admin_2"`
}

func (q *Queries) UpdateAdmin(ctx context.Context, arg UpdateAdminParams) (models.ContentAdmin, error) {
	row := q.db.QueryRowContext(ctx, updateAdmin, arg.Admin, arg.ContentID, arg.Admin_2)
	var i models.ContentAdmin
	err := row.Scan(&i.ID, &i.ContentID, &i.Admin)
	return i, err
}
