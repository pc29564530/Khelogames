package teams

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addPlayerToTeamRequest struct {
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
}

func (s *TeamsServer) AddTeamsMemberFunc(ctx *gin.Context) {
	var req addPlayerToTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	//authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.AddTeamPlayersParams{
		TeamID:   req.TeamID,
		PlayerID: req.PlayerID,
	}

	members, err := s.store.AddTeamPlayers(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add club member: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("successfully added member to the club")
	ctx.JSON(http.StatusAccepted, members)
	return
}

// this function should return the player logo full name position etch
func (s *TeamsServer) GetTeamsMemberFunc(ctx *gin.Context) {
	teamIDStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse team id string: %v", err)
		return
	}
	s.logger.Debug("get club id from reqeust: %v", teamID)

	membersID, err := s.store.GetTeamPlayers(ctx, teamID)
	if err != nil {
		s.logger.Error("Failed to get club member: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var playerList []map[string]interface{}
	for _, memberID := range membersID {
		player, err := s.store.GetPlayer(ctx, memberID.PlayerID)
		if err != nil {
			s.logger.Error("failed to get player data: ", err)
			return
		}
		playerData := map[string]interface{}{
			"player_name": player.Name,
			"slug":        player.Slug,
			"short_name":  player.ShortName,
			"position":    player.Positions,
			"country":     player.Country,
			"sports":      player.Sports,
			"player_logo": player.MediaUrl,
		}
		playerList = append(playerList, playerData)
	}

	s.logger.Info("successfully get club member")

	ctx.JSON(http.StatusAccepted, playerList)
	return
}
