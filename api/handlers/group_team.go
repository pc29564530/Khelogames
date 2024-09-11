package handlers

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupTeamServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewGroupTeamServer(store *db.Store, logger *logger.Logger) *GroupTeamServer {
	return &GroupTeamServer{store: store, logger: logger}
}

type addGroupTeamRequest struct {
	GroupID      int64 `json:"group_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *GroupTeamServer) AddGroupTeamFunc(ctx *gin.Context) {
	var req addGroupTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("bind the request: ", req)

	arg := db.AddGroupTeamParams{
		GroupID:      req.GroupID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}
	s.logger.Debug("params arg: ", arg)

	response, err := s.store.AddGroupTeam(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add group team: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully add group team: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *GroupTeamServer) GetTeamsByGroupFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	s.logger.Debug("tournament id parse: ", tournamentID)
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to group id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	s.logger.Debug("group id parse: ", groupID)

	arg := db.GetTeamByGroupParams{
		TournamentID: tournamentID,
		GroupID:      groupID,
	}
	s.logger.Debug("params arg: ", arg)
	response, err := s.store.GetTeamByGroup(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get team by group: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully get team by group: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
