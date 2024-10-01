package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const createFootballStatistics = `
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
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id, match_id, team_id, shots_on_target, total_shots, corner_kicks, fouls, goalkeeper_saves, free_kicks, yellow_cards, red_cards
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
SELECT id, match_id, team_id, shots_on_target, total_shots, corner_kicks, fouls, goalkeeper_saves, free_kicks, yellow_cards, red_cards FROM football_statistics
WHERE match_id=$1 AND team_id=$2
`

type GetFootballStatisticsParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetFootballStatistics(ctx context.Context, arg GetFootballStatisticsParams) (models.FootballStatistic, error) {
	row := q.db.QueryRowContext(ctx, getFootballStatistics, arg.MatchID, arg.TeamID)
	var i models.FootballStatistic
	err := row.Scan(
		&i.ID,
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
    shots_on_target = shots_on_target + $1,
    total_shots = total_shots + $2,
    corner_kicks = corner_kicks + $3,
    fouls = fouls + $4,
    goalkeeper_saves = goalkeeper_saves + $5,
    free_kicks = free_kicks + $6,
    yellow_cards = yellow_cards + $7,
    red_cards = red_cards + $8
WHERE match_id = $9 AND team_id = $10
RETURNING id, match_id, team_id, shots_on_target, total_shots, corner_kicks, fouls, goalkeeper_saves, free_kicks, yellow_cards, red_cards
`

type UpdateFootballStatisticsParams struct {
	ShotsOnTarget   int32 `json:"shots_on_target"`
	TotalShots      int32 `json:"total_shots"`
	CornerKicks     int32 `json:"corner_kicks"`
	Fouls           int32 `json:"fouls"`
	GoalkeeperSaves int32 `json:"goalkeeper_saves"`
	FreeKicks       int32 `json:"free_kicks"`
	YellowCards     int32 `json:"yellow_cards"`
	RedCards        int32 `json:"red_cards"`
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
}

func (q *Queries) UpdateFootballStatistics(ctx context.Context, arg UpdateFootballStatisticsParams) (models.FootballStatistic, error) {
	row := q.db.QueryRowContext(ctx, updateFootballStatistics,
		arg.ShotsOnTarget,
		arg.TotalShots,
		arg.CornerKicks,
		arg.Fouls,
		arg.GoalkeeperSaves,
		arg.FreeKicks,
		arg.YellowCards,
		arg.RedCards,
		arg.MatchID,
		arg.TeamID,
	)
	var i models.FootballStatistic
	err := row.Scan(
		&i.ID,
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
