package tournaments

import (
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

	err := ctx.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	entityPublicID, err := uuid.Parse(req.EntityPublicID)
	if err != nil {
		s.logger.Error("Failed to parse: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// NOTE: Commented out authorization check - uncomment if needed
	// authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	// tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	// if err != nil {
	// 	s.logger.Error("Failed to get tournament: ", err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"success": false,
	// 		"code":    "DATABASE_ERROR",
	// 		"message": "Failed to get tournament",
	// 	})
	// 	return
	// }
	// isExists, err := s.store.GetTournamentUserRole(ctx, int32(tournament.ID), authPayload.UserID)
	// if err != nil {
	// 	ctx.JSON(404, gin.H{"error": "Check  failed"})
	// 	return
	// }
	// if !isExists {
	// 	ctx.JSON(403, gin.H{"error": "You do not own this tournament participants"})
	// 	return
	// }

	var groupID *int32
	if req.GroupID != 0 {
		gid := req.GroupID
		groupID = &gid
	}

	participants, err := s.store.AddTournamentParticipants(ctx, tournamentPublicID, groupID, entityPublicID, req.EntityType, req.SeedNumber, req.Status)
	if err != nil {
		s.logger.Error("Failed to add tournament participants: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add tournament participants",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    participants,
	})
}

func (s *TournamentServer) GetTournamentParticipantsFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	tournamentParticipants, err := s.store.GetTournamentParticipants(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament participants: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournament pariticipants",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    tournamentParticipants,
	})
}
