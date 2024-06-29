package handlers

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TournamentGroupServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewTournamentGroup(store *db.Store, logger *logger.Logger) *TournamentGroupServer {
	return &TournamentGroupServer{store: store, logger: logger}
}

type createTournamentGroupRequest struct {
	GroupName     string `json:"group_name"`
	TournamentID  int64  `json:"tournament_id"`
	GroupStrength int64  `json:"group_strength"`
}

func (s *TournamentGroupServer) CreateTournamentGroupFunc(ctx *gin.Context) {
	var req createTournamentGroupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateTournamentGroupParams{
		GroupName:     req.GroupName,
		TournamentID:  req.TournamentID,
		GroupStrength: req.GroupStrength,
	}

	response, err := s.store.CreateTournamentGroup(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create tournament group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentGroupServer) GetTournamentGroupFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse group id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTournamentGroupParams{
		GroupID:      groupID,
		TournamentID: tournamentID,
	}

	response, err := s.store.GetTournamentGroup(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to get tournament group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentGroupServer) GetTournamentGroupsFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	response, err := s.store.GetTournamentGroups(ctx, tournamentID)
	if err != nil {
		fmt.Errorf("Failed to get tournament group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
