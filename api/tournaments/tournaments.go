package tournaments

import (
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getTournamentPublicIDRequest struct {
	TournamentPublicID string `uri:"tournament_public_id"`
}

type addTournamentRequest struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	Level          string `json:"level"`
	StartTimestamp string `json:"start_timestamp"`
	GameID         *int64 `json:"game_id"`
	GroupCount     *int32 `json:"group_count"`
	MaxGroupTeams  *int32 `json:"max_group_teams"`
	Stage          string `json:"stage"`
	HasKnockout    bool   `json:"has_knockout"`
	City           string `json:"city"`
	State          string `json:"state"`
	Country        string `json:"country"`
}

func (s *TournamentServer) AddTournamentFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create a tournament")
	var req addTournamentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	s.logger.Debug("Bind data: ", req)
	timestamp := req.StartTimestamp

	// Sanitize user text inputs to prevent XSS
	sanitizedName := util.SanitizeString(req.Name)
	sanitizedCity := util.SanitizeString(req.City)
	sanitizedState := util.SanitizeString(req.State)
	sanitizedCountry := util.SanitizeString(req.Country)

	slug := util.GenerateSlug(sanitizedName)
	startTimeStamp, err := util.ConvertTimeStamp(timestamp)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "TIMESTAMP_ERROR",
			"message": "Failed to convert time to second",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.txStore.AddNewTournamentTx(ctx,
		authPayload,
		sanitizedName,    // Sanitized
		slug,
		req.Status,
		req.Level,
		startTimeStamp,
		req.GameID,
		req.GroupCount,
		req.MaxGroupTeams,
		req.Stage,
		req.HasKnockout,
		sanitizedCity,    // Sanitized
		sanitizedState,   // Sanitized
		sanitizedCountry, // Sanitized
	)
	if err != nil {
		s.logger.Error("Failed to create tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to create tournament",
		})
		return
	}

	s.logger.Info("Successfully created the tournament")

	ctx.JSON(http.StatusAccepted, tournament)
	return
}

func (s *TournamentServer) GetTournamentTeamCountFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get team count for tournament")
	var req getTournamentPublicIDRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetTournamentTeamsCount(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get team count: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get team count",
		})
		return
	}
	s.logger.Info("Successfully retrieved team count: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get tournaments")

	response, err := s.store.GetTournaments(ctx)
	if err != nil {
		s.logger.Error("Failed to get tournaments: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}
	s.logger.Info("Successfully retrieved tournaments: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get a tournament")
	var req getTournamentPublicIDRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}
	s.logger.Info("Successfully retrieved tournament: %v", response)

	ctx.JSON(http.StatusAccepted, response)
}

func (s *TournamentServer) UpdateTournamentDateFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update tournament dates")
	var req getTournamentPublicIDRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	startOnStr := ctx.Query("start_on")
	startTimeStamp, err := util.ConvertTimeStamp(startOnStr)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "TIMESTAMP_ERROR",
			"message": "Invalid timestamp format",
		})
		return
	}

	arg := db.UpdateTournamentDateParams{
		TournamentPublicID: tournamentPublicID,
		StartTimestamp:     startTimeStamp,
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}

	if tournament.UserID != authPayload.UserID {
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"code":    "AUTHORIZATION_ERROR",
			"message": "You are not allowed to make change",
		})
		return
	}

	response, err := s.store.UpdateTournamentDate(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update tournament dates: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to update tournament dates",
		})
		return
	}
	s.logger.Info("Successfully updated tournament dates: ", response)

	ctx.JSON(http.StatusAccepted, response)
}

func (s *TournamentServer) GetTournamentByLevelFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get tournaments by level")

	sports := ctx.Param("sport")
	level := ctx.Query("category")
	s.logger.Debug("Category: %v", level)

	game, err := s.store.GetGamebyName(ctx, sports)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get game",
		})
		return
	}

	response, err := s.store.GetTournamentsByLevel(ctx, game.ID, level)
	if err != nil {
		s.logger.Error("Failed to get tournaments by level: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament by level",
		})
		return
	}
	s.logger.Info("Successfully retrieved tournaments by level: ", response)
	ctx.JSON(http.StatusAccepted, response)
}

func (s *TournamentServer) UpdateTournamentStatusFunc(ctx *gin.Context) {
	var req getTournamentPublicIDRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	statusCode := ctx.Query("status_code")

	arg := db.UpdateTournamentStatusParams{
		TournamentPublicID: tournamentPublicID,
		Status:             statusCode,
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}

	if tournament.UserID != authPayload.UserID {
		s.logger.Error("Failed to match tournament user id with current user: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Failed to match user_id",
		})
		return
	}

	updatedMatchData, err := s.store.UpdateTournamentStatus(ctx, arg)
	if err != nil {
		s.logger.Error("unable to update the tournament status: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to update tournament status",
		})
		return
	}

	s.logger.Info("successfully updated the tournament status")
	ctx.JSON(http.StatusAccepted, updatedMatchData)
}

type getTournamentByGameIdRequest struct {
	GameID int64 `uri:"game_id"`
}

func (s *TournamentServer) GetTournamentsBySportFunc(ctx *gin.Context) {
	var req getTournamentByGameIdRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get game",
		})
		return
	}

	tournaments, err := s.store.GetTournamentsBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get the tournaments by sport: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament by sport",
		})
		return
	}

	var results map[string]interface{}

	results = map[string]interface{}{
		"id":          game.ID,
		"name":        game.Name,
		"min_players": game.MinPlayers,
		"tournament":  tournaments,
	}

	ctx.JSON(http.StatusOK, results)
}

func (s *TournamentServer) AddTournamentUserRolesFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `json:"tournament_public_id"`
		Role               string `json:"role"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind add tournament user roles: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse uuid: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}

	newTournamentUserRoles, err := s.store.AddTournamentUserRoles(ctx, int32(tournament.ID), authPayload.UserID, req.Role)
	if err != nil {
		s.logger.Error("Failed to add new tournament user roles: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get user role",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, newTournamentUserRoles)
}

func (s *TournamentServer) GetTournamentByLocationFunc(ctx *gin.Context) {
	// var req struct {
	// 	City    string `json:"city"`
	// 	State   string `json:"state"`
	// 	Country string `json:"country"`
	// }
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	s.logger.Error("Failed to bind: ", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
	// 	return
	// }

	city := ctx.Query("city")
	state := ctx.Query("state")
	country := ctx.Query("country")

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get game",
		})
		return
	}

	tournament, err := s.store.GetTournamentByLocation(ctx, game.ID, city, state, country)
	if err != nil {
		s.logger.Error("Failed to get the tournaments: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}

	var results map[string]interface{}

	results = map[string]interface{}{
		"id":          game.ID,
		"name":        game.Name,
		"min_players": game.MinPlayers,
		"tournament":  tournament,
	}

	ctx.JSON(http.StatusOK, results)
}
