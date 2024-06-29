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

type TournamentOrganizerServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewTournamentOrganizerServer(store *db.Store, logger *logger.Logger) *TournamentOrganizerServer {
	return &TournamentOrganizerServer{store: store, logger: logger}
}

type createTournamentOrganizationRequest struct {
	TournamentID    int64     `json:"tournament_id"`
	PlayerCount     int64     `json:"player_count"`
	TeamCount       int64     `json:"team_count"`
	GroupCount      int64     `json:"group_count"`
	AdvancedTeam    int64     `json:"advanced_team"`
	TournamentStart time.Time `json:"tournament_start"`
}

func (s *TournamentOrganizerServer) CreateTournamentOrganizationFunc(ctx *gin.Context) {

	var req createTournamentOrganizationRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateTournamentOrganizationParams{
		TournamentID:    req.TournamentID,
		TournamentStart: req.TournamentStart,
		PlayerCount:     req.PlayerCount,
		TeamCount:       req.TeamCount,
		GroupCount:      req.GroupCount,
		AdvancedTeam:    req.AdvancedTeam,
	}

	response, err := s.store.CreateTournamentOrganization(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentOrganizerServer) GetTournamentOrganizationFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil || tournamentIDStr == " " {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response, err := s.store.GetTournamentOrganization(ctx, tournamentID)
	if err != nil {
		fmt.Errorf("Failed to get tournament organization: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
