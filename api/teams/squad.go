package teams

import (
	db "khelogames/database"
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addPlayerToTeamRequest struct {
	PlayerPublicID uuid.UUID `json:"player_public_id"`
	TeamPublicID   uuid.UUID `json:"team_public_id"`
	JoinDate       string    `json:"join_date"`
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

	checkPlayerExist := s.store.GetTeamPlayer(ctx, req.TeamPublicID, req.PlayerPublicID)
	if checkPlayerExist {
		var leaveData int32
		_, err := s.store.RemovePlayerFromTeam(ctx, req.TeamPublicID, req.PlayerPublicID, leaveData)
		if err != nil {
			s.logger.Error("Failed to update the the leave date: ", err)
			return
		}

		player, err := s.store.GetPlayerByPublicID(ctx, req.PlayerPublicID)
		if err != nil {
			s.logger.Error("Failed to get the player: ", err)
			return
		}
		playerData := map[string]interface{}{
			"id":         player.ID,
			"public_id":  player.PublicID,
			"name":       player.Name,
			"slug":       player.Slug,
			"short_name": player.ShortName,
			"position":   player.Positions,
			"country":    player.Country,
			"media_url":  player.MediaUrl,
			"game_id":    player.GameID,
		}
		ctx.JSON(http.StatusAccepted, playerData)
		return
	} else {
		arg := db.AddTeamPlayersParams{
			TeamPublicID:   req.TeamPublicID,
			PlayerPublicID: req.PlayerPublicID,
			JoinDate:       int32(startDate),
			LeaveDate:      nil,
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
}

func (s *TeamsServer) GetTeamsMemberFunc(ctx *gin.Context) {
	var req struct {
		TeamPublicID uuid.UUID `uri:"team_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	players, err := s.store.GetPlayerByTeam(ctx, req.TeamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team member: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var playerList []map[string]interface{}
	for _, player := range players {
		playerData := map[string]interface{}{
			"id":          player.ID,
			"public_id":   player.PublicID,
			"player_name": player.PlayerName,
			"slug":        player.Slug,
			"short_name":  player.ShortName,
			"position":    player.Positions,
			"country":     player.Country,
			"media_url":   player.MediaUrl,
			"game_id":     player.GameID,
		}
		playerList = append(playerList, playerData)
	}

	s.logger.Info("successfully get club member")

	ctx.JSON(http.StatusAccepted, playerList)
	return
}

type removePlayerFromTeamRequest struct {
	TeamPublicID   uuid.UUID `json:"team_public_id"`
	PlayerPublicID uuid.UUID `json:"player_public_id"`
	LeaveDate      string    `json:"leave_date"`
}

func (s *TeamsServer) RemovePlayerFromTeamFunc(ctx *gin.Context) {
	var req removePlayerFromTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	//var endDate *int64
	endDate, err := util.ConvertTimeStamp(req.LeaveDate)
	if err != nil {
		s.logger.Error("Failed to convert to second")
		return
	}

	leave := int32(endDate)
	var eddPointer *int32 = &leave

	response, err := s.store.RemovePlayerFromTeam(ctx, req.TeamPublicID, req.PlayerPublicID, *eddPointer)
	if err != nil {
		s.logger.Error("Failed to remove player from team: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
