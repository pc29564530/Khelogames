package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

// tournament
type createTournamentRequest struct {
	TournamentName string `json:"tournament_name"`
	SportType      string `json:"sport_type"`
	Format         string `json:"format"`
	TeamsJoined    int64  `json:"teams_joined`
}

func (server *Server) createTournament(ctx *gin.Context) {
	var req createTournamentRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateTournamentParams{
		TournamentName: req.TournamentName,
		SportType:      req.SportType,
		Format:         req.Format,
		TeamsJoined:    req.TeamsJoined,
	}

	fmt.Println("Arg: line 34:L ", arg)

	response, err := server.store.CreateTournament(ctx, arg)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := server.store.GetTeamsCount(ctx, req.TouurnamentID)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.UpdateTeamsJoinedParams{
		TeamsJoined:  req.TeamsJoined,
		TournamentID: req.TournamentID,
	}

	response, err := server.store.UpdateTeamsJoined(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getTournaments(ctx *gin.Context) {

	response, err := server.store.GetTournaments(ctx)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := server.store.GetTournament(ctx, req.TournamentID)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateOrganizerParams{
		OrganizerName: req.OrganizerName,
		TournamentID:  req.TournamentID,
	}

	fmt.Println("Arg: ", arg)
	response, err := server.store.CreateOrganizer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getOrganizerRequest struct {
	TournamentID int64 `uri:"tournament_id"`
}

func (server *Server) getOrganizer(ctx *gin.Context) {
	var req getOrganizerRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	fmt.Println("Id: ", req.TournamentID)
	response, err := server.store.GetOrganizer(ctx, req.TournamentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	fmt.Println("Respoonse: ", response)
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddTeamParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	fmt.Println("Team: params: ", arg)

	response, err := server.store.AddTeam(ctx, arg)
	if err != nil {
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := server.store.GetTeams(ctx, req.TournamentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
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
