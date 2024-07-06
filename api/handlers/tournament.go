package handlers

import (
	"fmt"
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
	TeamsJoined    int64     `json:"teams_joined`
	StartOn        time.Time `json:"start_on"`
	EndOn          time.Time `json:"end_on"`
	Category       string    `json:"category"`
}

func (s *TournamentServer) CreateTournamentFunc(ctx *gin.Context) {
	var req createTournamentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

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
		fmt.Errorf("Failed to create tournament: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTournamentTeamRequest struct {
	TouurnamentID int64 `uri:"tournament_id"`
}

func (s *TournamentServer) GetTournamentTeamCountFunc(ctx *gin.Context) {
	var req getTournamentTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := s.store.GetTeamsCount(ctx, req.TouurnamentID)
	if err != nil {
		fmt.Errorf("Failed to get team count: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type updateTeamsJoinedRequest struct {
	TeamsJoined  int64 `json:"teams_joined"`
	TournamentID int64 `json:"tournament_id"`
}

func (s *TournamentServer) UpdateTeamsJoinedFunc(ctx *gin.Context) {
	var req updateTeamsJoinedRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.UpdateTeamsJoinedParams{
		TeamsJoined:  req.TeamsJoined,
		TournamentID: req.TournamentID,
	}

	response, err := s.store.UpdateTeamsJoined(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update team joined: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentsFunc(ctx *gin.Context) {

	response, err := s.store.GetTournaments(ctx)
	if err != nil {
		fmt.Errorf("Failed to get tournament: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentsBySportFunc(ctx *gin.Context) {

	sport := ctx.Param("sport")

	response, err := s.store.GetTournamentsBySport(ctx, sport)
	if err != nil {
		fmt.Errorf("Failed to get tournament by sport: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTournamentRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (s *TournamentServer) GetTournamentFunc(ctx *gin.Context) {
	var req getTournamentRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := s.store.GetTournament(ctx, req.TournamentID)
	if err != nil {
		fmt.Errorf("Failed to get tournament: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type createOrganizerRequest struct {
	OrganizerName string `json:"organizer_name"`
	TournamentID  int64  `json:"tournament_id"`
}

func (s *TournamentServer) CreateOrganizerFunc(ctx *gin.Context) {
	var req createOrganizerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateOrganizerParams{
		OrganizerName: req.OrganizerName,
		TournamentID:  req.TournamentID,
	}

	s.logger.Debug("organizer arg: %v", arg)

	response, err := s.store.CreateOrganizer(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create organizer: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetOrganizerFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
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
		fmt.Errorf("Failed to get organizer: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type addTournamentTeamRequest struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) AddTeamFunc(ctx *gin.Context) {
	var req addTournamentTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddTeamParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	response, err := s.store.AddTeam(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to add team: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type getTeamsRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (s *TournamentServer) GetTeamsFunc(ctx *gin.Context) {
	var req getTeamsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := s.store.GetTeams(ctx, req.TournamentID)
	if err != nil {
		fmt.Errorf("Failed to get teams: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type getTeamRequest struct {
	TeamID int64 `uri:"team_id"`
}

func (s *TournamentServer) GetTeamFunc(ctx *gin.Context) {
	var req getTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	response, err := s.store.GetTeam(ctx, req.TeamID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (s *TournamentServer) UpdateTournamentDateFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	startOnStr := ctx.Query("start_on")
	layout := "2000-05-05"
	startOn, err := time.Parse(layout, startOnStr)
	if err != nil {
		s.logger.Debug("Error parsing date:", err)
		return
	}

	endOnStr := ctx.Query("end_on")

	endOn, err := time.Parse(layout, endOnStr)
	if err != nil {
		fmt.Errorf("Failed to parse data: %v", err)
		return
	}

	arg := db.UpdateTournamentDateParams{
		StartOn:      startOn,
		EndOn:        endOn,
		TournamentID: tournamentID,
	}

	response, err := s.store.UpdateTournamentDate(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update tournament date: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentByLevelFunc(ctx *gin.Context) {
	sport := ctx.Param("sport")
	category := ctx.Query("category")
	s.logger.Debug("Category: ", category)
	arg := db.GetTournamentByLevelParams{
		SportType: sport,
		Category:  category,
	}

	s.logger.Debug("Tournament by level arg: %v", arg)

	response, err := s.store.GetTournamentByLevel(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to get tournament by level: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
