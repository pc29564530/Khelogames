package transactions

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

type NewTournamentParams struct {
	UserPublicID   uuid.UUID `json:"user_public_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	Country        string    `json:"country"`
	Status         string    `json:"status"`
	Season         int       `json:"season"`
	Level          string    `json:"level"`
	StartTimestamp int64     `json:"start_timestamp"`
	GameID         *int64    `json:"game_id"`
	GroupCount     *int32    `json:"group_count"`
	MaxGroupTeams  *int32    `json:"max_group_teams"`
	Stage          string    `json:"stage"`
	HasKnockout    bool      `json:"has_knockout"`
	IsPublic       bool      `json:"is_public"`
}

func (s *SQLStore) AddNewTournamentTx(ctx context.Context, arg database.NewTournamentParams) (*models.Tournament, error) {
	var newTournament models.Tournament
	err := s.execTx(ctx, func(q *database.Queries) error {
		var err error
		newTournament, err = q.NewTournament(ctx, arg)
		if err != nil {
			s.logger.Error("Failed to create tournament: ", err)
			return err
		}

		organizer := "organizer"

		_, err = q.AddTournamentUserRoles(ctx, int32(newTournament.ID), newTournament.UserID, organizer)
		if err != nil {
			s.logger.Error("Failed to add tournament user roles: ", err)
			return err
		}
		return nil
	})

	return &newTournament, err
}
