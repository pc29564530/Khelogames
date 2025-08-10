package tournaments

import (
	"encoding/json"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"
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

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

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
		Country:        req.Country,
		Status:         req.Status,
		Level:          req.Level,
		StartTimestamp: startTimeStamp,
		GameID:         req.GameID,
		GroupCount:     req.GroupCount,
		MaxGroupTeams:  req.MaxGroupTeams,
		Stage:          req.Stage,
		HasKnockout:    req.HasKnockout,
	}

	newTournament, err := s.store.NewTournament(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	// authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// argAdmin := db.AddAdminParams{
	// 	ContentID: newTournament.ID,
	// 	Admin:     authPayload.UserID,
	// }

	// _, err = s.store.AddAdmin(ctx, argAdmin)
	// if err != nil {
	// 	s.logger.Error("Failed to add admin for the tournament: ", err)
	// }

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
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
