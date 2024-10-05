package tournaments

import (
	"encoding/json"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addTournamentRequest struct {
	TournamentName string `json:"tournament_name"`
	Slug           string `json:"slug"`
	Sports         string `json:"sports"`
	Country        string `json:"country"`
	StatusCode     string `json:"status_code"`
	Level          string `json:"level"`
	StartTimestamp string `json:"start_timestamp"`
	GameID         int64  `json:"game_id"`
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

	slug := util.GenerateSlug(req.TournamentName)
	startTimeStamp, err := util.ConvertTimeStamp(timestamp)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		return
	}

	arg := db.NewTournamentParams{
		TournamentName: req.TournamentName,
		Slug:           slug,
		Sports:         req.Sports,
		Country:        req.Country,
		StatusCode:     req.StatusCode,
		Level:          req.Level,
		StartTimestamp: startTimeStamp,
		GameID:         req.GameID,
	}

	newTournament, err := s.store.NewTournament(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	argAdmin := db.AddAdminParams{
		ContentID: newTournament.ID,
		Admin:     authPayload.Username,
	}

	_, err = s.store.AddAdmin(ctx, argAdmin)
	if err != nil {
		s.logger.Error("Failed to add admin for the tournament: ", err)
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	s.logger.Info("Successfully created the tournament")

	ctx.JSON(http.StatusAccepted, newTournament)
	return
}

type getTournamentTeamRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (s *TournamentServer) GetTournamentTeamCountFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get team count for tournament")
	var req getTournamentTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := s.store.GetTournamentTeamsCount(ctx, req.TournamentID)
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

type getTournamentRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (s *TournamentServer) GetTournamentFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get a tournament")
	var req getTournamentRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := s.store.GetTournament(ctx, req.TournamentID)
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
	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament ID: ", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	startOnStr := ctx.Query("start_on")
	startTimeStamp, err := util.ConvertTimeStamp(startOnStr)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		return
	}

	arg := db.UpdateTournamentDateParams{
		StartTimestamp: startTimeStamp,
		ID:             tournamentID,
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

	arg := db.GetTournamentsByLevelParams{
		Sports: sports,
		Level:  level,
	}

	s.logger.Debug("GetTournamentByLevelParams: ", arg)

	response, err := s.store.GetTournamentsByLevel(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournaments by level: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved tournaments by level: ", response)
	ctx.JSON(http.StatusAccepted, response)
}

func (s *TournamentServer) UpdateTournamentStatusFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse the match id: ", err)
		return
	}

	statusCode := ctx.Query("status_code")

	arg := db.UpdateTournamentStatusParams{
		ID:         tournamentID,
		StatusCode: statusCode,
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
	var tournamentData []map[string]interface{}
	var gameDetail map[string]interface{}
	for _, row := range rows {
		var tournamentDetails map[string]interface{}
		err := json.Unmarshal((row.TournamentData), &tournamentDetails)
		if err != nil {
			s.logger.Error("Failed to unmarshal tournament data: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tournament data"})
			return
		}

		tournamentData = append(tournamentData, tournamentDetails)
		gameDetail = map[string]interface{}{
			"id":          row.ID,
			"name":        row.Name,
			"min_players": row.MinPlayers,
		}
	}

	results = map[string]interface{}{
		"game":       gameDetail,
		"tournament": tournamentData,
	}

	ctx.JSON(http.StatusOK, results)
}
