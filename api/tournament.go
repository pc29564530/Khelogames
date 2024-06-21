package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

func (server *Server) createTournament(ctx *gin.Context) {
	var req createTournamentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
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

	response, err := server.store.CreateTournament(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to create tournament: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTournamentTeamRequest struct {
	TouurnamentID int64 `uri:"tournament_id"`
}

func (server *Server) getTournamentTeamCount(ctx *gin.Context) {
	var req getTournamentTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := server.store.GetTeamsCount(ctx, req.TouurnamentID)
	if err != nil {
		server.logger.Error("Failed to get team count: %v", err)
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

func (server *Server) updateTeamsJoined(ctx *gin.Context) {
	var req updateTeamsJoinedRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.UpdateTeamsJoinedParams{
		TeamsJoined:  req.TeamsJoined,
		TournamentID: req.TournamentID,
	}

	response, err := server.store.UpdateTeamsJoined(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update team joined: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournaments(ctx *gin.Context) {

	response, err := server.store.GetTournaments(ctx)
	if err != nil {
		server.logger.Error("Failed to get tournament: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournamentsBySport(ctx *gin.Context) {

	sport := ctx.Param("sport")

	response, err := server.store.GetTournamentsBySport(ctx, sport)
	if err != nil {
		server.logger.Error("Failed to get tournament by sport: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getTournamentRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (server *Server) getTournament(ctx *gin.Context) {
	var req getTournamentRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := server.store.GetTournament(ctx, req.TournamentID)
	if err != nil {
		server.logger.Error("Failed to get tournament: %v", err)
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

func (server *Server) createOrganizer(ctx *gin.Context) {
	var req createOrganizerRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateOrganizerParams{
		OrganizerName: req.OrganizerName,
		TournamentID:  req.TournamentID,
	}

	server.logger.Info("organizer arg: %v", arg)

	response, err := server.store.CreateOrganizer(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to create organizer: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getOrganizer(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		server.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	organizerName := ctx.Query("organizer_name")

	arg := db.GetOrganizerParams{
		TournamentID:  tournamentID,
		OrganizerName: organizerName,
	}

	response, err := server.store.GetOrganizer(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to get organizer: %v", err)
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

func (server *Server) addTeam(ctx *gin.Context) {
	var req addTournamentTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddTeamParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	response, err := server.store.AddTeam(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to add team: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type getTeamsRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (server *Server) getTeams(ctx *gin.Context) {
	var req getTeamsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := server.store.GetTeams(ctx, req.TournamentID)
	if err != nil {
		server.logger.Error("Failed to get teams: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

type getTeamRequest struct {
	TeamID int64 `uri:"team_id"`
}

func (server *Server) getTeam(ctx *gin.Context) {
	var req getTeamRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	response, err := server.store.GetTeam(ctx, req.TeamID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (server *Server) updateTournamentDate(ctx *gin.Context) {
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
		fmt.Println("Error parsing date:", err)
		return
	}

	endOnStr := ctx.Query("end_on")

	endOn, err := time.Parse(layout, endOnStr)
	if err != nil {
		server.logger.Error("Failed to parse data: %v", err)
		return
	}

	arg := db.UpdateTournamentDateParams{
		StartOn:      startOn,
		EndOn:        endOn,
		TournamentID: tournamentID,
	}

	response, err := server.store.UpdateTournamentDate(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to update tournament date: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournamentByLevel(ctx *gin.Context) {
	sport := ctx.Param("sport")
	category := ctx.Query("category")
	fmt.Println("Category: ", category)
	arg := db.GetTournamentByLevelParams{
		SportType: sport,
		Category:  category,
	}

	server.logger.Info("Tournament by level arg: %v", arg)

	response, err := server.store.GetTournamentByLevel(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to get tournament by level: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
