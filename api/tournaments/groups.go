package tournaments

import (
	db "khelogames/database"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createTournamentGroupRequest struct {
	Name         string `json:"name"`
	TournamentID int64  `json:"tournament_id"`
	Strength     int32  `json:"strength"`
}

func (s *TournamentServer) CreateTournamentGroupFunc(ctx *gin.Context) {
	var req createTournamentGroupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("bind request: ", req)
	arg := db.CreateTournamentGroupParams{
		Name:         req.Name,
		TournamentID: req.TournamentID,
		Strength:     req.Strength,
	}
	s.logger.Debug("params arg: ", arg)

	response, err := s.store.CreateTournamentGroup(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create tournament group: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully created tournament group: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTournamentGroupFunc(ctx *gin.Context) {
	tournamentIDStr := ctx.Query("tournament_id")
	groupIDStr := ctx.Query("id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse group id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	arg := db.GetTournamentGroupParams{
		ID:           groupID,
		TournamentID: tournamentID,
	}

	response, err := s.store.GetTournamentGroup(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournament group: ", err)
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
		s.logger.Error("Failed to parse tournament id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}
	s.logger.Debug("parse the tournamend id: ", tournamentID)

	response, err := s.store.GetTournamentGroups(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament group: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully get tournament groups: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

type addGroupTeamRequest struct {
	GroupID      int64 `json:"group_id" form:"group_id"`
	TournamentID int64 `json:"tournament_id" form:"tournament_id"`
	TeamID       int64 `json:"team_id" form:"team_id"`
}

func (s *TournamentServer) AddGroupTeamFunc(ctx *gin.Context) {
	var req addGroupTeamRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		s.logger.Error("Failed to bind add group team: ", err)
		return
	}
	s.logger.Debug("bind the request: ", req)

	arg := db.CreateGroupTeamsParams{
		GroupID:      req.GroupID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}
	s.logger.Debug("params arg: ", arg)

	response, err := s.store.CreateGroupTeams(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add group team: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully add group team: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetTeamsByGroupFunc(ctx *gin.Context) {
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

	arg := db.GetGroupTeamsParams{
		TournamentID: tournamentID,
		GroupID:      groupID,
	}
	s.logger.Debug("params arg: ", arg)
	response, err := s.store.GetGroupTeams(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get team by group: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully get team by group: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TournamentServer) GetGroupsFunc(ctx *gin.Context) {

	response, err := s.store.GetGroups(ctx)
	if err != nil {
		s.logger.Error("Failed to get groups: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("successfully get group: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
