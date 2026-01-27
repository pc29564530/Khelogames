package tournaments

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createTournamentStandingRequest struct {
	GroupID            int32  `json:"group_id"`
	TeamPublicID       string `json:"team_public_id"`
	TournamentPublicID string `json:"tournament_public_id"`
}

func (s *TournamentServer) CreateTournamentStandingFunc(ctx *gin.Context) {
	var req createTournamentStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"tournament_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(tournament.ID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user tournament role: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND_ERROR",
				"message": "Check failed",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	if !isExists {
		s.logger.Error("Tournament user role does not exist: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	game := ctx.Param("sport")

	if game == "football" {
		footballStanding, err := s.store.CreateFootballStanding(ctx, tournamentPublicID, req.GroupID, teamPublicID)
		if err != nil {
			s.logger.Error("Failed to create football standing: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to create football standing",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"success": true,
			"data":    footballStanding,
		})
	} else if game == "cricket" {
		cricketStanding, err := s.store.CreateCricketStanding(ctx, tournamentPublicID, req.GroupID, teamPublicID)
		if err != nil {
			s.logger.Error("Failed to create cricket standing: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to create cricket standing",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"success": true,
			"data":    cricketStanding,
		})
	}
}
