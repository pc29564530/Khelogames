package tournaments

import (
	"encoding/json"
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getFootballStandingRequest struct {
	TournamentPublicID string `uri:"tournament_public_id"`
}

func (s *TournamentServer) GetFootballStandingFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get football standing")
	var req getFootballStandingRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"tournament_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	rows, err := s.store.GetFootballStanding(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get football standing: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get football standing",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Check if rows or rows.StandingData is nil
	if rows == nil || rows.StandingData == nil {
		s.logger.Warn("No standings data available")
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []interface{}{},
			"message": "No standings data available",
		})
		return
	}

	var standings []map[string]interface{}
	var standingsData []map[string]interface{}
	var standingData []interface{}

	err = json.Unmarshal(rows.StandingData.([]byte), &standingData)
	if err != nil {
		s.logger.Error("Failed to unmarshal standing data: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "DATA_PARSE_ERROR",
				"message": "Failed to process standing data",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Validate we have data to process
	if len(standingData) == 0 {
		s.logger.Warn("Empty standings data after unmarshal")
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []interface{}{},
			"message": "No standings data available",
		})
		return
	}

	for _, data := range standingData {
		// Ensure data is of type map[string]interface{}
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			s.logger.Error("Invalid data format in standing data")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "DATA_PARSE_ERROR",
					"message": "Invalid standing data format",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		standingsData = append(standingsData, map[string]interface{}{
			"tournament":      dataMap["tournament"],
			"group":           dataMap["group"],
			"teams":           dataMap["teams"],
			"tournament_id":   dataMap["torunament_id"],
			"group_id":        dataMap["group_id"],
			"id":              dataMap["id"],
			"public_id":       dataMap["public_id"],
			"matches":         dataMap["matches"],
			"wins":            dataMap["wins"],
			"loss":            dataMap["loss"],
			"draw":            dataMap["draw"],
			"goal_for":        dataMap["goal_for"],
			"goal_against":    dataMap["goal_against"],
			"goal_difference": dataMap["goal_difference"],
			"points":          dataMap["points"],
		})
	}

	groupData := make(map[int64][]map[string]interface{})
	visited := make(map[int]string)

	for _, standing := range standingsData {
		if standing["group_id"] == nil {
			ind := -1
			groupData[int64(ind)] = append(groupData[int64(ind)], map[string]interface{}{
				"teams":           standing["teams"],
				"id":              standing["id"],
				"public_id":       standing["public_id"],
				"matches":         standing["matches"],
				"wins":            standing["wins"],
				"loss":            standing["loss"],
				"draw":            standing["draw"],
				"goal_for":        standing["goal_for"],
				"goal_against":    standing["goal_against"],
				"goal_difference": standing["goal_difference"],
				"points":          standing["points"],
			})
		} else {
			groupID := standing["group_id"]
			grpID, ok := groupID.(float64)
			if !ok {
				s.logger.Error("Invalid group_id type")
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "DATA_PARSE_ERROR",
						"message": "Invalid group ID format",
					},
					"request_id": ctx.GetString("request_id"),
				})
				return
			}

			// Append standings data to groupData by groupID
			groupData[int64(grpID)] = append(groupData[int64(grpID)], map[string]interface{}{
				"teams":           standing["teams"],
				"id":              standing["id"],
				"public_id":       standing["public_id"],
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
			if _, exists := visited[int(int64(grpID))]; !exists {
				if groupMap, ok := standing["group"].(map[string]interface{}); ok {
					if groupName, ok := groupMap["name"].(string); ok {
						visited[int(int64(grpID))] = groupName
					}
				}
			}
		}
	}

	// Add grouped standings to the final standings slice
	if len(standingsData) > 0 {
		standings = append(standings, map[string]interface{}{
			"tournament": standingsData[0]["tournament"],
		})
	}

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

	s.logger.Info("Successfully retrieved football standing")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    standings,
	})
}
