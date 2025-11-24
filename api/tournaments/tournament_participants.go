package tournaments

import (
	"khelogames/core/token"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TournamentServer) AddTournamentParticipantsFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `json:"tournament_public_id"`
		GroupID            int32  `json:"group_id"`
		EntityPublicID     string `json:"entity_public_id"` //team or player
		EntityType         string `json:"entity_type"`      //team or player
		SeedNumber         int    `json:"seed_number"`
		Status             string `json:"status"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse: ", err)
		return
	}

	entityPublicID, err := uuid.Parse(req.EntityPublicID)
	if err != nil {
		s.logger.Error("Failed to parse: ", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Tournamet not found"})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(tournament.ID), authPayload.UserID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Check  failed"})
		return
	}
	if !isExists {
		ctx.JSON(403, gin.H{"error": "You do not own this tournament participants"})
		return
	}

	var groupID *int32

	participants, err := s.store.AddTournamentParticipants(ctx, tournamentPublicID, groupID, entityPublicID, req.EntityType, req.SeedNumber, req.Status)
	if err != nil {
		s.logger.Error("Failed to add tournament participants: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, participants)
}

func (s *TournamentServer) GetTournamentParticipantsFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: ", err)
		return
	}

	tournamentParticipants, err := s.store.GetTournamentParticipants(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament participants: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, tournamentParticipants)
}
