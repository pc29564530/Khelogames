package tournaments

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TournamentServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewTournamentServer(store *db.Store, logger *logger.Logger) *TournamentServer {
	return &TournamentServer{store: store, logger: logger}
}

// tournament
type createTournamentRequest struct {
	TournamentName string    `json:"tournament_name"`
	SportType      string    `json:"sport_type"`
	Format         string    `json:"format"`
	TeamsJoined    int64     `json:"teams_joined"`
	StartOn        time.Time `json:"start_on"`
	EndOn          time.Time `json:"end_on"`
	Category       string    `json:"category"`
}

func (s *TournamentServer) CreateTournamentFunc(ctx *gin.Context) {
	s.logger.Info("Received request to create a tournament")
	var req createTournamentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Request data: %v", req)

	arg := db.CreateTournamentParams{
		TournamentName: req.TournamentName,
		SportType:      req.SportType,
		Format:         req.Format,
		TeamsJoined:    req.TeamsJoined,
		StartOn:        req.StartOn,
		EndOn:          req.EndOn,
		Category:       req.Category,
	}

	response, err := s.store.CreateTournament(ctx, arg)
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

	response, err := s.store.GetTeamsCount(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("Failed to get team count: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved team count: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateTeamsJoinedRequest struct {
	TeamsJoined  int64 `json:"teams_joined"`
	TournamentID int64 `json:"tournament_id"`
}

func (s *TournamentServer) UpdateTeamsJoinedFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update teams joined")
	var req updateTeamsJoinedRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.UpdateTeamsJoinedParams{
		TeamsJoined:  req.TeamsJoined,
		TournamentID: req.TournamentID,
	}

	response, err := s.store.UpdateTeamsJoined(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update teams joined: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully updated teams joined: %v", response)

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

type addTournamentTeamRequest struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) AddTeamFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add a team")
	var req addTournamentTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddTeamParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	response, err := s.store.AddTeam(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add team: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully added team: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTeamsRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (s *TournamentServer) GetTeamsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get teams for a tournament")
	var req getTeamsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := s.store.GetTeams(ctx, req.TournamentID)
	if err != nil {
		s.logger.Error("Failed to get teams: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not found"})
		return
	}
	s.logger.Info("Successfully retrieved teams: %v", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTeamRequest struct {
	TeamID int64 `uri:"team_id"`
}

func (s *TournamentServer) GetTeamFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get a team")
	var req getTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := s.store.GetTeam(ctx, req.TeamID)
	if err != nil {
		s.logger.Error("Failed to get team: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved team: %v", response)

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
	layout := "2006-01-02"
	startOn, err := time.Parse(layout, startOnStr)
	if err != nil {
		s.logger.Error("Failed to parse start date: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	endOnStr := ctx.Query("end_on")
	endOn, err := time.Parse(layout, endOnStr)
	if err != nil {
		s.logger.Error("Failed to parse end date: %v", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	arg := db.UpdateTournamentDateParams{
		StartOn:      startOn,
		EndOn:        endOn,
		TournamentID: tournamentID,
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

	sport := ctx.Param("sport")
	category := ctx.Query("category")
	s.logger.Debug("Category: %v", category)

	arg := db.GetTournamentByLevelParams{
		SportType: sport,
		Category:  category,
	}

	s.logger.Debug("GetTournamentByLevelParams: %v", arg)

	response, err := s.store.GetTournamentByLevel(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournaments by level: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved tournaments by level: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
