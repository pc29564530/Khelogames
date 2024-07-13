package tournaments

import (
	db "khelogames/db/sqlc"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createTournamentGroupRequest struct {
	GroupName     string `json:"group_name"`
	TournamentID  int64  `json:"tournament_id"`
	GroupStrength int64  `json:"group_strength"`
}

func (s *TournamentServer) CreateTournamentGroupFunc(ctx *gin.Context) {
	var req createTournamentGroupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("bind request: %v", req)
	arg := db.CreateTournamentGroupParams{
		GroupName:     req.GroupName,
		TournamentID:  req.TournamentID,
		GroupStrength: req.GroupStrength,
	}
	s.logger.Debug("params arg: %v", arg)

	response, err := s.store.CreateTournamentGroup(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create tournament group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully created tournament group: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentGroupFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse group id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTournamentGroupParams{
		GroupID:      groupID,
		TournamentID: tournamentID,
	}

	response, err := s.store.GetTournamentGroup(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournament group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentGroupsFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	s.logger.Debug("parse the tournamend id: %v", tournamentID)

	response, err := s.store.GetTournamentGroups(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully get tournament groups: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

type addGroupTeamRequest struct {
	GroupID      int64 `json:"group_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) AddGroupTeamFunc(ctx *gin.Context) {
	var req addGroupTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("bind the request: %v", req)

	arg := db.AddGroupTeamParams{
		GroupID:      req.GroupID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}
	s.logger.Debug("params arg: %v", arg)

	response, err := s.store.AddGroupTeam(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add group team: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully add group team: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTeamsByGroupFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("group_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	s.logger.Debug("tournament id parse: %v", tournamentID)
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to group id: %v", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	s.logger.Debug("group id parse: %v", groupID)

	arg := db.GetTeamByGroupParams{
		TournamentID: tournamentID,
		GroupID:      groupID,
	}
	s.logger.Debug("params arg: %v", arg)
	response, err := s.store.GetTeamByGroup(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get team by group: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully get team by group: %v", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
