package tournaments

import (
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
		s.logger.Error("Failed to get cricket standing: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var standings []map[string]interface{}
	var standingsData []map[string]interface{}

	for _, row := range *rows {

		if row.StandingData == nil {
			s.logger.Warn("No standings data available")
			ctx.JSON(http.StatusNoContent, gin.H{"message": "No standings data available"})
			return
		}

		dataMap := row.StandingData

		// Safely extract values with type assertions
		tournamentID, _ := dataMap.(map[string]interface{})["tournament_id"].(float64)
		groupID := dataMap.(map[string]interface{})["group_id"].(float64)
		id := dataMap.(map[string]interface{})["id"].(float64)
		matches := dataMap.(map[string]interface{})["matches"].(float64)
		wins := dataMap.(map[string]interface{})["wins"].(float64)
		loss := dataMap.(map[string]interface{})["loss"].(float64)
		draw := dataMap.(map[string]interface{})["draw"].(float64)
		points := dataMap.(map[string]interface{})["point"].(float64)
		standing := map[string]interface{}{
			"tournament":    dataMap.(map[string]interface{})["details"].(map[string]interface{})["tournament"],
			"group":         dataMap.(map[string]interface{})["details"].(map[string]interface{})["group"],
			"teams":         dataMap.(map[string]interface{})["details"].(map[string]interface{})["teams"],
			"tournament_id": int64(tournamentID),
			"group_id":      int64(groupID),
			"id":            int64(id),
			"matches":       int32(matches),
			"wins":          int32(wins),
			"loss":          int32(loss),
			"draw":          int32(draw),
			"points":        int32(points),
		}
		standingsData = append(standingsData, standing)
	}

	groupData := make(map[int64][]map[string]interface{})
	visited := make(map[int]string)

	for _, standing := range standingsData {
		groupID := standing["group_id"]
		grpID := groupID.(int64)
		groupData[grpID] = append(groupData[grpID], map[string]interface{}{
			"teams":   standing["teams"],
			"id":      standing["id"],
			"matches": standing["matches"],
			"wins":    standing["wins"],
			"loss":    standing["loss"],
			"draw":    standing["draw"],
			"points":  standing["points"],
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
