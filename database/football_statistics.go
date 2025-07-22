package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const createFootballStatistics = `
WITH resolved_ids AS (
  SELECT 
    m.id AS match_id,
    t.id AS team_id
  FROM matches m, teams t
  WHERE m.public_id = $1 AND t.public_id = $2
)
INSERT INTO football_statistics (
    match_id,
    team_id,
    shots_on_target,
    total_shots,
    corner_kicks,
    fouls,
    goalkeeper_saves,
    free_kicks,
    yellow_cards,
    red_cards
)
SELECT 
    r.match_id,
    r.team_id,
    $3, $4, $5, $6, $7, $8, $9, $10
FROM resolved_ids r
RETURNING *;
`

type CreateFootballStatisticsParams struct {
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
	ShotsOnTarget   int32 `json:"shots_on_target"`
	TotalShots      int32 `json:"total_shots"`
	CornerKicks     int32 `json:"corner_kicks"`
	Fouls           int32 `json:"fouls"`
	GoalkeeperSaves int32 `json:"goalkeeper_saves"`
	FreeKicks       int32 `json:"free_kicks"`
	YellowCards     int32 `json:"yellow_cards"`
	RedCards        int32 `json:"red_cards"`
}

func (q *Queries) CreateFootballStatistics(ctx context.Context, arg CreateFootballStatisticsParams) (models.FootballStatistic, error) {
	row := q.db.QueryRowContext(ctx, createFootballStatistics,
		arg.MatchID,
		arg.TeamID,
		arg.ShotsOnTarget,
		arg.TotalShots,
		arg.CornerKicks,
		arg.Fouls,
		arg.GoalkeeperSaves,
		arg.FreeKicks,
		arg.YellowCards,
		arg.RedCards,
	)
	var i models.FootballStatistic
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.ShotsOnTarget,
		&i.TotalShots,
		&i.CornerKicks,
		&i.Fouls,
		&i.GoalkeeperSaves,
		&i.FreeKicks,
		&i.YellowCards,
		&i.RedCards,
	)
	return i, err
}

const getFootballStatistics = `
SELECT * FROM football_statistics fs
JOIN matches m ON m.id = fs.match_id
JOIN teams t ON t.id = fs.team_id
WHERE m.public_id=$1 AND t.public_id=$2
`

func (q *Queries) GetFootballStatistics(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) (models.FootballStatistic, error) {
	row := q.db.QueryRowContext(ctx, getFootballStatistics, matchPublicID, teamPublicID)
	var i models.FootballStatistic
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.ShotsOnTarget,
		&i.TotalShots,
		&i.CornerKicks,
		&i.Fouls,
		&i.GoalkeeperSaves,
		&i.FreeKicks,
		&i.YellowCards,
		&i.RedCards,
	)
	return i, err
}

const updateFootballStatistics = `
UPDATE football_statistics
SET 
    shots_on_target = shots_on_target + $2,
    total_shots = total_shots + $3,
    corner_kicks = corner_kicks + $4,
    fouls = fouls + $5,
    goalkeeper_saves = goalkeeper_saves + $6,
    free_kicks = free_kicks + $7,
    yellow_cards = yellow_cards + $9,
    red_cards = red_cards + $10
FROM matches m, teams t
WHERE m.public_id = $1 AND t.public_id = $2 AND m.id = fs.match_id AND t.id = fs.team_id
RETURNING *
`

type UpdateFootballStatisticsParams struct {
	MatchPublicID   uuid.UUID `json:"match_public_id"`
	TeamPublicID    uuid.UUID `json:"team_public_id"`
	ShotsOnTarget   int32     `json:"shots_on_target"`
	TotalShots      int32     `json:"total_shots"`
	CornerKicks     int32     `json:"corner_kicks"`
	Fouls           int32     `json:"fouls"`
	GoalkeeperSaves int32     `json:"goalkeeper_saves"`
	FreeKicks       int32     `json:"free_kicks"`
	YellowCards     int32     `json:"yellow_cards"`
	RedCards        int32     `json:"red_cards"`
}

func (q *Queries) UpdateFootballStatistics(ctx context.Context, arg UpdateFootballStatisticsParams) (models.FootballStatistic, error) {
	row := q.db.QueryRowContext(ctx, updateFootballStatistics,
		arg.MatchPublicID,
		arg.TeamPublicID,
		arg.ShotsOnTarget,
		arg.TotalShots,
		arg.CornerKicks,
		arg.Fouls,
		arg.GoalkeeperSaves,
		arg.FreeKicks,
		arg.YellowCards,
		arg.RedCards,
	)
	var i models.FootballStatistic
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.ShotsOnTarget,
		&i.TotalShots,
		&i.CornerKicks,
		&i.Fouls,
		&i.GoalkeeperSaves,
		&i.FreeKicks,
		&i.YellowCards,
		&i.RedCards,
	)
	return i, err
}
