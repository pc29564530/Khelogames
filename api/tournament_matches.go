package api

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type createTournamentMatchRequest struct {
	OrganizerID  int64     `json:"organizer_id"`
	TournamentID int64     `json:"tournament_id"`
	Team1ID      int64     `json:"team1_id"`
	Team2ID      int64     `json:"team2_id"`
	DateON       time.Time `json:"date_on"`
	StartAt      time.Time `json:"start_at"`
	Stage        string    `json:"stage"`
	Sports       string    `json:"sports"`
}

func (server *Server) createTournamentMatch(ctx *gin.Context) {
	var req createTournamentMatchRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateMatchParams{
		OrganizerID:  req.OrganizerID,
		TournamentID: req.TournamentID,
		Team1ID:      req.Team1ID,
		Team2ID:      req.Team2ID,
		DateOn:       req.DateON,
		StartAt:      req.StartAt,
		Stage:        req.Stage,
		Sports:       req.Sports,
	}

	response, err := server.store.CreateMatch(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getAllTournamentMatch(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	sports := ctx.Query("sports")
	arg := db.GetTournamentMatchParams{
		TournamentID: tournamentID,
		Sports:       sports,
	}

	response, err := server.store.GetTournamentMatch(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}