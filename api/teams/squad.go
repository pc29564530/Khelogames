package teams

import (
	db "khelogames/database"
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addPlayerToTeamRequest struct {
	PlayerID int64  `json:"player_id"`
	TeamID   int64  `json:"team_id"`
	JoinDate string `json:"join_date"`
}

func (s *TeamsServer) AddTeamsMemberFunc(ctx *gin.Context) {
	var req addPlayerToTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	//convert the date to second to insert into the teamplayer table
	startDate, err := util.ConvertTimeStamp(req.JoinDate)
	if err != nil {
		s.logger.Error("Failed to convert the timestamp to second ", err)
	}
	arg := db.AddTeamPlayersParams{
		TeamID:    req.TeamID,
		PlayerID:  req.PlayerID,
		JoinDate:  int32(startDate),
		LeaveDate: nil,
	}

	members, err := s.store.AddTeamPlayers(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add club member: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("successfully added member to the club")
	ctx.JSON(http.StatusAccepted, members)
}

func (s *TeamsServer) GetTeamsMemberFunc(ctx *gin.Context) {
	teamIDStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse team id string: ", err)
		return
	}
	s.logger.Debug("get club id from reqeust:", teamID)

	players, err := s.store.GetPlayerByTeam(ctx, teamID)
	if err != nil {
		s.logger.Error("Failed to get team member: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var playerList []map[string]interface{}
	for _, player := range players {
		playerData := map[string]interface{}{
			"id":          player.ID,
			"player_name": player.PlayerName,
			"slug":        player.Slug,
			"short_name":  player.ShortName,
			"position":    player.Positions,
			"country":     player.Country,
			"sports":      player.Sports,
			"media_url":   player.MediaUrl,
		}
		playerList = append(playerList, playerData)
	}

	s.logger.Info("successfully get club member")

	ctx.JSON(http.StatusAccepted, playerList)
	return
}

type removePlayerFromTeamRequest struct {
	TeamID    int64  `json:"team_id"`
	PlayerID  int64  `json:"player_id"`
	LeaveDate string `json:"leave_date"`
}

func (s *TeamsServer) RemovePlayerFromTeamFunc(ctx *gin.Context) {
	var req removePlayerFromTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	endDate, err := util.ConvertTimeStamp(req.LeaveDate)
	if err != nil {
		s.logger.Error("Failed to convert to second")
		return
	}

	response, err := s.store.RemovePlayerFromTeam(ctx, req.TeamID, req.PlayerID, int32(endDate))
	if err != nil {
		s.logger.Error("Failed to remove player from team: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
