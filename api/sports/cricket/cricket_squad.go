package cricket

import (
	"net/http"

	errorhandler "khelogames/error_handler"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type Player struct {
	ID         int64  `json:"id"`
	PublicID   string `json:"public_id"`
	UserID     int32  `json:"user_id"`
	PlayerName string `json:"player_name"`
	ShortName  string `json:"short_name"`
	Slug       string `json:"slug"`
	Country    string `json:"country"`
	Position   string `json:"position"`
	MediaURL   string `json:"media_url"`
	GameID     int64  `json:"game_id"`
	OnBench    bool   `json:"on_bench"`
}

type MatchSquadRequest struct {
	MatchPublicID string   `json:"match_public_id"`
	TeamPublicID  string   `json:"team_public_id"`
	Player        []Player `json:"player"`
}

func (s *CricketServer) AddCricketSquadFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add cricket squad")
	var req MatchSquadRequest

	err := ctx.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid team UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	var cricketSquad []map[string]interface{}
	for _, player := range req.Player {
		var err error
		playerPublicID, err := uuid.Parse(player.PublicID)
		if err != nil {
			s.logger.Error("Invalid player UUID format: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "Invalid player UUID format",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		squad, err := s.store.AddCricketSquad(ctx, matchPublicID, teamPublicID, playerPublicID, player.Position, player.OnBench, false)
		if err != nil {
			s.logger.Error("Failed to add cricket squad: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to add cricket squad",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		cricketSquad = append(cricketSquad, map[string]interface{}{
			"id":         squad.ID,
			"public_id":  squad.PublicID,
			"match_id":   squad.MatchID,
			"team_id":    squad.TeamID,
			"player":     player,
			"role":       squad.Role,
			"on_bench":   squad.OnBench,
			"is_captain": squad.IsCaptain,
		})
	}

	s.logger.Info("Successfully added cricket squad")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cricketSquad,
	})
}

func (s *CricketServer) GetCricketMatchSquadFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get cricket match squad")
	matchPublicIDStr := ctx.Query("match_public_id")
	teamPublicIDStr := ctx.Query("team_public_id")

	fieldErrors := make(map[string]string)

	if matchPublicIDStr == "" {
		fieldErrors["match_public_id"] = "Match public ID is required"
	}

	if teamPublicIDStr == "" {
		fieldErrors["team_public_id"] = "Team public ID is required"
	}

	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(matchPublicIDStr)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(teamPublicIDStr)
	if err != nil {
		s.logger.Error("Invalid team UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	cricketSquad, err := s.store.GetCricketMatchSquad(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get cricket squad: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get cricket squad",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully retrieved cricket match squad")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cricketSquad,
	})
}
