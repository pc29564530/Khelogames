package teams

import (
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addPlayerToTeamRequest struct {
	TeamPublicID   string `json:"team_public_id"`
	PlayerPublicID string `json:"player_public_id"`
	JoinDate       string `json:"join_date"`
}

func (s *TeamsServer) AddTeamsMemberFunc(ctx *gin.Context) {
	var req addPlayerToTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	playerPublicID, err := uuid.Parse(req.PlayerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	//convert the date to second to insert into the teamplayer table
	startDate, err := util.ConvertTimeStamp(req.JoinDate)
	if err != nil {
		s.logger.Error("Failed to convert the timestamp to second ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid date format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team by public id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get team",
		})
		return
	}

	if team.UserID != authPayload.UserID {
		s.logger.Error("You do not own this team")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "FORBIDDEN_ERROR",
			"message": "You do not own this team",
		})
		return
	}

	checkPlayerExist := s.store.GetTeamPlayer(ctx, teamPublicID, playerPublicID)
	if checkPlayerExist {
		var leaveData int32
		_, err := s.store.RemovePlayerFromTeam(ctx, teamPublicID, playerPublicID, leaveData)
		if err != nil {
			s.logger.Error("Failed to update the the leave date: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update leave date",
			})
			return
		}

		player, err := s.store.GetPlayerByPublicID(ctx, playerPublicID)
		if err != nil {
			s.logger.Error("Failed to get the player: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player",
			})
			return
		}
		playerData := map[string]interface{}{
			"id":         player.ID,
			"public_id":  player.PublicID,
			"user_id":    player.UserID,
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
			TeamPublicID:   teamPublicID,
			PlayerPublicID: playerPublicID,
			JoinDate:       int32(startDate),
			LeaveDate:      nil,
		}

		members, err := s.store.AddTeamPlayers(ctx, arg)
		if err != nil {
			s.logger.Error("Failed to add team member: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add team member",
			})
			return
		}
		s.logger.Info("successfully added member to the team")
		ctx.JSON(http.StatusAccepted, members)
	}
}

func (s *TeamsServer) GetTeamsMemberFunc(ctx *gin.Context) {
	var req struct {
		TeamPublicID string `uri:"team_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get team member",
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	players, err := s.store.GetPlayerByTeam(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team member: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get team member",
		})
		return
	}
	s.logger.Info("successfully get team member")

	ctx.JSON(http.StatusAccepted, players)
	return
}

type removePlayerFromTeamRequest struct {
	TeamPublicID   string `json:"team_public_id"`
	PlayerPublicID string `json:"player_public_id"`
	LeaveDate      string `json:"leave_date"`
}

func (s *TeamsServer) RemovePlayerFromTeamFunc(ctx *gin.Context) {
	var req removePlayerFromTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to remove player from team",
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	playerPublicID, err := uuid.Parse(req.PlayerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	//var endDate *int64
	endDate, err := util.ConvertTimeStamp(req.LeaveDate)
	if err != nil {
		s.logger.Error("Failed to convert to second")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid date format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team by public id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get team",
		})
		return
	}

	if authPayload.UserID != int32(team.UserID) {
		s.logger.Error("You do not own this team")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "FORBIDDEN_ERROR",
			"message": "You do not own this team",
		})
		return
	}

	leave := int32(endDate)
	var eddPointer *int32 = &leave

	response, err := s.store.RemovePlayerFromTeam(ctx, teamPublicID, playerPublicID, *eddPointer)
	if err != nil {
		s.logger.Error("Failed to remove player from team: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to remove player from team",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
