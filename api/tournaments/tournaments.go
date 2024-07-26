package tournaments

import (
	db "khelogames/db/sqlc"
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
}

func (s *TournamentServer) AddTournamentFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create a tournament")
	var req addTournamentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request data: %v", req)
	timestamp := req.StartTimestamp

	slug := util.GenerateSlug(req.TournamentName)
	startTimeStamp, err := util.ConvertTimeStamp(timestamp)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		return
	}
	// fmt.Println("req.StatusCode", req.StatusCode)

	// statusCodeJSON, err := json.Marshal(req.StatusCode)
	// if err != nil {
	// 	s.logger.Error("failed to marshal status code: %v", err)
	// }

	arg := db.NewTournamentParams{
		TournamentName: req.TournamentName,
		Slug:           slug,
		Sports:         req.Sports,
		Country:        req.Country,
		StatusCode:     req.StatusCode,
		Level:          req.Level,
		StartTimestamp: startTimeStamp,
	}

	response, err := s.store.NewTournament(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create tournament: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully created the tournament: %v", response)

	ctx.JSON(http.StatusAccepted, response)
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

func (s *TournamentServer) GetTournamentsBySportFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get tournaments by sport")

	sport := ctx.Param("sport")

	response, err := s.store.GetTournamentsBySport(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get tournaments by sport: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved tournaments by sport: %v", response)

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

type createOrganizerRequest struct {
	OrganizerName string `json:"organizer_name"`
	TournamentID  int64  `json:"tournament_id"`
}

func (s *TournamentServer) CreateOrganizerFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create an organizer")
	var req createOrganizerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateOrganizerParams{
		OrganizerName: req.OrganizerName,
		TournamentID:  req.TournamentID,
	}

	s.logger.Debug("Organizer arg: %v", arg)

	response, err := s.store.CreateOrganizer(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create organizer: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully created organizer: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetOrganizerFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get an organizer")

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	organizerName := ctx.Query("organizer_name")

	arg := db.GetOrganizerParams{
		TournamentID:  tournamentID,
		OrganizerName: organizerName,
	}

	response, err := s.store.GetOrganizer(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get organizer: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved organizer: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) UpdateTournamentDateFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update tournament dates")
	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament ID: %v", err)
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
		s.logger.Error("Failed to update tournament dates: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated tournament dates: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
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

	s.logger.Debug("GetTournamentByLevelParams: %v", arg)

	response, err := s.store.GetTournamentsByLevel(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournaments by level: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved tournaments by level: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
