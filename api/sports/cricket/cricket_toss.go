package cricket

import (
	"khelogames/core/token"
	"khelogames/pkg"
	"net/http"

	errorhandler "khelogames/error_handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addCricketTossRequest struct {
	MatchPublicID string `json:"match_public_id"`
	TossDecision  string `json:"toss_decision"`
	TossWin       string `json:"toss_win"`
}

func (s *CricketServer) AddCricketTossFunc(ctx *gin.Context) {
	var req addCricketTossRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
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

	tossWin, err := uuid.Parse(req.TossWin)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
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

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user tournament role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	if !isExists {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN_ERROR",
				"message": "You are not allowed to add toss details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	inningResponse, tossWinTeam, tossDecision, err := s.txStore.AddCricketTossTx(ctx, matchPublicID, req.TossDecision, tossWin)
	if err != nil {
		s.logger.Error("Failed to add cricket toss: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add cricket toss",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"inning": gin.H{
			"id":                  inningResponse.ID,
			"public_id":           inningResponse.PublicID,
			"match_id":            inningResponse.MatchID,
			"team_id":             inningResponse.TeamID,
			"inning_number":       inningResponse.InningNumber,
			"score":               inningResponse.Score,
			"wickets":             inningResponse.Wickets,
			"overs":               inningResponse.Overs,
			"run_rate":            inningResponse.RunRate,
			"target_run_rate":     inningResponse.TargetRunRate,
			"follow_on":           inningResponse.FollowOn,
			"is_inning_completed": inningResponse.IsInningCompleted,
			"declared":            inningResponse.Declared,
			"inning_status":       inningResponse.InningStatus,
		},
		"team":          tossWinTeam,
		"toss_decision": tossDecision,
	})
	return
}

type getTossRequest struct {
	MatchPublicID string `uri:"match_public_id"`
}

func (s *CricketServer) GetCricketTossFunc(ctx *gin.Context) {

	var req getTossRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind cricket toss : ", err)
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
	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
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

	cricketToss, err := s.store.GetCricketToss(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get the cricket match toss: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get cricket match toss",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	tossWonTeamMap, ok := cricketToss["toss_won_team"].(map[string]interface{})
	if !ok {
		s.logger.Error("Invalid toss_won_team format")
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid toss_won_team format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	publicIDStr, ok := tossWonTeamMap["public_id"].(string)
	if !ok {
		s.logger.Error("Invalid public_id format")
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid public_id format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	publicID, err := uuid.Parse(publicIDStr)
	if err != nil {
		s.logger.Error("Failed to parse public_id as UUID: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format for public_id",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	tossWonTeam, err := s.store.GetTeamByPublicID(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get team: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	tossDecision, ok := cricketToss["toss_decision"].(string)
	if !ok {
		s.logger.Error("Invalid toss_decision format")
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid toss_decision format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	tossDetails := map[string]interface{}{

		"tossWonTeam": map[string]interface{}{
			"id":           tossWonTeam.ID,
			"public_id":    tossWonTeam.PublicID,
			"name":         tossWonTeam.Name,
			"slug":         tossWonTeam.Slug,
			"shortName":    tossWonTeam.Shortname,
			"gender":       tossWonTeam.Gender,
			"national":     tossWonTeam.National,
			"country":      tossWonTeam.Country,
			"type":         tossWonTeam.Type,
			"player_count": tossWonTeam.PlayerCount,
			"game_id":      tossWonTeam.GameID,
		},
		"tossDecision": tossDecision,
	}

	s.logger.Debug("toss won team details: ", tossDetails)

	ctx.JSON(http.StatusAccepted, tossDetails)
}
