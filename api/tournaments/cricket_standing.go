package tournaments

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TournamentServer) GetCricketStandingFunc(ctx *gin.Context) {

	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	rows, err := s.store.GetCricketStanding(ctx, tournamentPublicID)
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

		standingsData = append(standingsData, map[string]interface{}{
			"tournament":    dataMap["tournament"],
			"group":         dataMap["group"],
			"teams":         dataMap["teams"],
			"tournament_id": dataMap["torunament_id"],
			"group_id":      dataMap["group_id"],
			"id":            dataMap["id"],
			"public_id":     dataMap["public_id"],
			"matches":       dataMap["matches"],
			"wins":          dataMap["wins"],
			"loss":          dataMap["loss"],
			"draw":          dataMap["draw"],
			"points":        dataMap["points"],
		})
	}

	groupData := make(map[int64][]map[string]interface{})
	visited := make(map[int]string) // Initialize the visited map to prevent nil map errors

	for _, standing := range standingsData {

		if standing["group_id"] == nil {
			ind := -1
			groupData[int64(ind)] = append(groupData[int64(ind)], map[string]interface{}{
				"teams":     standing["teams"],
				"id":        standing["id"],
				"public_id": standing["public_id"],
				"matches":   standing["matches"],
				"wins":      standing["wins"],
				"loss":      standing["loss"],
				"draw":      standing["draw"],
				"points":    standing["points"],
			})
		} else {
			groupID := standing["group_id"]
			grpID := groupID.(float64)

			// Append standings data to groupData by groupID
			groupData[int64(grpID)] = append(groupData[int64(grpID)], map[string]interface{}{
				"teams":     standing["teams"],
				"id":        standing["id"],
				"public_id": standing["public_id"],
				"matches":   standing["matches"],
				"wins":      standing["wins"],
				"loss":      standing["loss"],
				"draw":      standing["draw"],
				"points":    standing["points"],
			})

			// Set the group name if not already visited
			if _, ok := visited[int(int64(grpID))]; !ok {
				vis, ok := standing["group"].(map[string]interface{})["name"].(string)
				if ok {
					visited[int(int64(grpID))] = vis
				}
			}
		}
	}

	// Add grouped standings to the final standings slice
	standings = append(standings, map[string]interface{}{
		"tournament": standingsData[0]["tournament"],
	})
	for grpID, grpData := range groupData {

		var groupName string
		if visited[int(grpID)] != "" {
			groupName = visited[int(grpID)]
		} else {
			groupName = "League"
		}

		standings = append(standings, map[string]interface{}{
			"group_name": groupName,
			"team_row":   grpData,
		})
	}

	ctx.JSON(http.StatusAccepted, standings)
}
