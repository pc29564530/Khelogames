package tournaments

import (
	"encoding/json"
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
	Country        string `json:"country"`
	Status         string `json:"status"`
	Level          string `json:"level"`
	StartTimestamp string `json:"start_timestamp"`
	GameID         *int64 `json:"game_id"`
	GroupCount     *int32 `json:"group_count"`
	MaxGroupTeams  *int32 `json:"max_group_teams"`
	Stage          string `json:"stage"`
	HasKnockout    bool   `json:"has_knockout"`
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

	arg := db.NewTournamentParams{
		UserPublicID:   authPayload.PublicID,
		Name:           req.Name,
		Slug:           slug,
		Description:    "",
		Country:        req.Country,
		Status:         req.Status,
		Season:         1,
		Level:          req.Level,
		StartTimestamp: startTimeStamp,
		GameID:         req.GameID,
		GroupCount:     req.GroupCount,
		MaxGroupTeams:  req.MaxGroupTeams,
		Stage:          req.Stage,
		HasKnockout:    req.HasKnockout,
		IsPublic:       true,
	}

	newTournament, err := s.txStore.AddNewTournamentTx(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	s.logger.Info("Successfully created the tournament")

	ctx.JSON(http.StatusAccepted, newTournament)
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

	rows, err := s.store.GetTournamentsBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get the tournaments: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tournaments"})
		return
	}

	var results map[string]interface{}
	var gameDetail map[string]interface{}
	var tournaments []map[string]interface{}
	for _, row := range rows {
		gameDetail = map[string]interface{}{
			"id":          row.ID,
			"name":        row.Name,
			"min_players": row.MinPlayers,
		}
		var tournament map[string]interface{}
		tt := (row.Tournament).([]byte)
		err := json.Unmarshal(tt, &tournament)
		if err != nil {
			s.logger.Error("Failed to unmarshal tournament data: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tournament data"})
			return
		}

		tournaments = append(tournaments, tournament)
	}

	results = map[string]interface{}{
		"id":          gameDetail["id"],
		"name":        gameDetail["name"],
		"min_players": gameDetail["min_players"],
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
