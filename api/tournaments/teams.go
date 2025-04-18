package tournaments

import (
	"encoding/json"
	db "khelogames/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addTournamentTeamRequest struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) AddTournamentTeamFunc(ctx *gin.Context) {
	s.logger.Info("Received request to add a team")
	var req addTournamentTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}

	teamPlayer, err := s.store.GetTeamByPlayer(ctx, req.TeamID)
	if err != nil {
		s.logger.Error("Failed to get team player: ", err)
		return
	}

	teamPlayerCount := len(teamPlayer)

	if game.MinPlayers > int32(teamPlayerCount) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Team strength does not satisfied"})
		return
	}

	arg := db.NewTournamentTeamParams{
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
	}

	newTeam, err := s.store.NewTournamentTeam(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add team: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	teamByTournament, err := s.store.GetTournamentTeam(ctx, newTeam.TeamID, newTeam.TournamentID)
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
	TournamentID int64 `uri:"tournament_id"`
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

	tournamentTeamsData, err := s.store.GetTournamentTeams(ctx, req.TournamentID)
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
