package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const loginQuery = `
INSERT INTO login (
    username,
    password
) VALUES (
    $1, $2
) RETURNING username, password
`

func (q *Queries) Login(ctx context.Context, username string, password string) (models.Login, error) {
	var Login models.Login
	row := q.db.QueryRowContext(ctx, loginQuery, username, password)
	err := row.Scan(&Login.Username, &Login.Password)
	return Login, err
}
