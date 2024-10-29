package tournaments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *TournamentServer) GetFootballStandingFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	rows, err := s.store.GetFootballStanding(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament standing: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var standingsData []map[string]interface{}
	for _, row := range rows {
		var data map[string]interface{}
		err := json.Unmarshal(row.StandingData.([]byte), &data)
		if err != nil {
			s.logger.Error("Failed to unmarshal ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal standings data"})
			return
		}

		standing := map[string]interface{}{
			"tournament":      data["tournament"],
			"group":           data["group"],
			"teams":           data["teams"],
			"tournament_id":   row.TournamentID,
			"group_id":        row.GroupID,
			"id":              row.ID,
			"matches":         row.Matches,
			"wins":            row.Wins,
			"loss":            row.Loss,
			"draw":            row.Draw,
			"goal_for":        row.GoalFor,
			"goal_against":    row.GoalAgainst,
			"goal_difference": row.GoalDifference,
			"points":          row.Points,
		}

		standingsData = append(standingsData, standing)
	}

	groupData := make(map[int64][]map[string]interface{})
	visited := make(map[int]string) // Initialize the visited map to prevent nil map errors
	var standings []map[string]interface{}

	for _, standing := range standingsData {
		groupID := standing["group_id"].(*int64)
		grpID := *groupID

		// Append standings data to groupData by groupID
		groupData[grpID] = append(groupData[grpID], map[string]interface{}{
			"teams":           standing["teams"],
			"id":              standing["id"],
			"matches":         standing["matches"],
			"wins":            standing["wins"],
			"loss":            standing["loss"],
			"draw":            standing["draw"],
			"goal_for":        standing["goal_for"],
			"goal_against":    standing["goal_against"],
			"goal_difference": standing["goal_difference"],
			"points":          standing["points"],
		})

		// Set the group name if not already visited
		if _, ok := visited[int(grpID)]; !ok {
			vis, ok := standing["group"].(map[string]interface{})["name"].(string)
			if ok {
				visited[int(grpID)] = vis
			}
		}
	}

	// Add grouped standings to the final standings slice
	standings = append(standings, map[string]interface{}{
		"tournament": standingsData[0]["tournament"],
	})
	for grpID, grpData := range groupData {
		groupName := visited[int(grpID)]
		standings = append(standings, map[string]interface{}{
			"group_name": groupName,
			"team_row":   grpData,
		})
	}
	ctx.JSON(http.StatusAccepted, standings)
}

type updateFootballStandingRequest struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) UpdateFootballStandingFunc(ctx *gin.Context) {
	var req updateFootballStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("bind the request: ", req)

	response, err := s.store.UpdateFootballStanding(ctx, req.TournamentID, req.TeamID)
	if err != nil {
		s.logger.Error("Failed to update tournament standing: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("successfully tournament standing: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
