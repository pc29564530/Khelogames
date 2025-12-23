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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	s.logger.Debug("Bind data: ", req)
	timestamp := req.StartTimestamp

	slug := util.GenerateSlug(req.Name)
	startTimeStamp, err := util.ConvertTimeStamp(timestamp)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.txStore.AddNewTournamentTx(ctx,
		authPayload,
		req.Name,
		slug,
		req.Status,
		req.Level,
		startTimeStamp,
		req.GameID,
		req.GroupCount,
		req.MaxGroupTeams,
		req.Stage,
		req.HasKnockout,
		req.City,
		req.State,
		req.Country,
	)
	if err != nil {
		s.logger.Error("Failed to create tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response, err := s.store.GetTournamentTeamsCount(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get team count: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	startOnStr := ctx.Query("start_on")
	startTimeStamp, err := util.ConvertTimeStamp(startOnStr)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		return
	}

	arg := db.UpdateTournamentDateParams{
		TournamentPublicID: tournamentPublicID,
		StartTimestamp:     startTimeStamp,
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Match not found"})
		return
	}

	if tournament.UserID != authPayload.UserID {
		ctx.JSON(403, gin.H{"error": "You do not own this match"})
		return
	}

	response, err := s.store.UpdateTournamentDate(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update tournament dates: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		return
	}

	response, err := s.store.GetTournamentsByLevel(ctx, game.ID, level)
	if err != nil {
		s.logger.Error("Failed to get tournaments by level: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
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
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
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
		ctx.JSON(404, gin.H{"error": "Match not found"})
		return
	}

	if tournament.UserID != authPayload.UserID {
		ctx.JSON(403, gin.H{"error": "You do not own this match"})
		return
	}

	updatedMatchData, err := s.store.UpdateTournamentStatus(ctx, arg)
	if err != nil {
		s.logger.Error("unable to update the tournament status: ", err)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}

	tournaments, err := s.store.GetTournamentsBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get the tournaments: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
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
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse uuid: ", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		return
	}

	newTournamentUserRoles, err := s.store.AddTournamentUserRoles(ctx, int32(tournament.ID), authPayload.UserID, req.Role)
	if err != nil {
		s.logger.Error("Failed to add new tournament user roles: ", err)
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
		return
	}

	tournament, err := s.store.GetTournamentByLocation(ctx, game.ID, city, state, country)
	if err != nil {
		s.logger.Error("Failed to get the tournaments: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
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
