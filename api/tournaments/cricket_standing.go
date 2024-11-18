package tournaments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *TournamentServer) GetCricketStandingFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")

	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: ", err)
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	rows, err := s.store.GetCricketStanding(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament standing: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	// Check if rows or rows.StandingData is nil
	if rows == nil || rows.StandingData == nil {
		s.logger.Warn("No standings data available")
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No standings data available"})
		return
	}

	var standings []map[string]interface{}
	var standingsData []map[string]interface{}
	var standingData []interface{}

	err = json.Unmarshal(rows.StandingData.([]byte), &standingData)
	if err != nil {
		s.logger.Error("Failed to unmarshal ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal standings data"})
		return
	}
	for _, data := range standingData {

		// Ensure data is of type map[string]interface{}
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			s.logger.Error("Invalid data format")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format"})
			return
		}

		// Safely extract and convert numeric fields
		tournamentID, _ := dataMap["tournament_id"].(float64)
		groupID, _ := dataMap["group_id"].(float64)
		id, _ := dataMap["id"].(float64)
		matches, _ := dataMap["matches"].(float64)
		wins, _ := dataMap["wins"].(float64)
		loss, _ := dataMap["loss"].(float64)
		draw, _ := dataMap["draw"].(float64)
		goalFor, _ := dataMap["goal_for"].(float64)
		goalAgainst, _ := dataMap["goal_against"].(float64)
		goalDifference, _ := dataMap["goal_difference"].(float64)
		points, _ := dataMap["points"].(float64)

		standing := map[string]interface{}{
			"tournament":      dataMap["tournament"],
			"group":           dataMap["group"],
			"teams":           dataMap["teams"],
			"tournament_id":   int64(tournamentID),
			"group_id":        int64(groupID),
			"id":              int64(id),
			"matches":         int64(matches),
			"wins":            int64(wins),
			"loss":            int64(loss),
			"draw":            int64(draw),
			"goal_for":        int64(goalFor),
			"goal_against":    int64(goalAgainst),
			"goal_difference": int64(goalDifference),
			"points":          int64(points),
		}
		standingsData = append(standingsData, standing)
	}

	groupData := make(map[int64][]map[string]interface{})
	visited := make(map[int]string) // Initialize the visited map to prevent nil map errors

	for _, standing := range standingsData {
		groupID := standing["group_id"]
		grpID := groupID

		// Append standings data to groupData by groupID
		groupData[grpID.(int64)] = append(groupData[grpID.(int64)], map[string]interface{}{
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
		if _, ok := visited[int(grpID.(int64))]; !ok {
			vis, ok := standing["group"].(map[string]interface{})["name"].(string)
			if ok {
				visited[int(grpID.(int64))] = vis
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

type updateCricketStandingRequest struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *TournamentServer) UpdateCricketStandingFunc(ctx *gin.Context) {
	var req updateCricketStandingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("bind the request: ", req)

	response, err := s.store.UpdateCricketStanding(ctx, req.TournamentID, req.TeamID)
	if err != nil {
		s.logger.Error("Failed to update tournament standing: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("successfully tournament standing: ", response)
	ctx.JSON(http.StatusAccepted, response)
	return
}
