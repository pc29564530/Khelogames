package tournaments

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addTournamentTeamRequest struct {
	TournamentPublicID string `json:"tournament_public_id"`
	TeamPublicID       string `json:"team_public_id"`
}

func (s *TournamentServer) AddTournamentTeamFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add a team")
	var req struct {
		TournamentPublicID string `json:"tournament_public_id"`
		TeamPublicID       string `json:"team_public_id"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public ID: ", err)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Failed to parse team public ID: ", err)
		return
	}

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}

	teamPlayer, err := s.store.GetPlayerByTeam(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team player: ", err)
		return
	}

	teamPlayerCount := len(teamPlayer)

	if game.MinPlayers > int32(teamPlayerCount) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Team strength does not satisfied"})
		return
	}

	_, err = s.store.NewTournamentTeam(ctx, tournamentPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to add team: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	teamByTournament, err := s.store.GetTournamentTeam(ctx, tournamentPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get Tournament Team: ", err)
		return
	}

	s.logger.Info("Successfully added team: ", teamByTournament.TeamData)
	var teamData map[string]interface{}
	err = json.Unmarshal([]byte(teamByTournament.TeamData), &teamData)
	if err != nil {
		s.logger.Error("Failed to unmarshal team data: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, teamData)
	return
}

type getTournamentTeamsRequest struct {
	TournamentPublicID string `uri:"tournament_public_id"`
}

func (s *TournamentServer) GetTournamentTeamsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get teams for a tournament")
	var req getTournamentTeamsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament public ID: ", err)
		return
	}

	tournamentTeamsData, err := s.store.GetTournamentTeams(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get teams: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not found"})
		return
	}

	var teamsData []map[string]interface{}
	for _, team := range tournamentTeamsData {
		var data map[string]interface{}
		err = json.Unmarshal(team.TeamData, &data)
		if err != nil {
			s.logger.Error("Failed to unmarshal: ", err)
			return
		}
		teamsData = append(teamsData, data)
	}

	s.logger.Info("Successfully retrieved teams: ", teamsData)

	ctx.JSON(http.StatusAccepted, teamsData)
}
