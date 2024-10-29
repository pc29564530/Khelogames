package tournaments

import (
	"encoding/json"
	"fmt"
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

	//need to response back with group

	rows, err := s.store.GetCricketStanding(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament standing: ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var standings []map[string]interface{}
	for _, row := range rows {
		var data map[string]interface{}

		tt := (row.StandingData).([]byte)

		err := json.Unmarshal(tt, &data)
		if err != nil {
			s.logger.Error("Failed to unmarshal ", err)
			return
		}
		standing := map[string]interface{}{
			"tournament": data["tournament"],
			"group":      data["group"],
			"teams":      data["teams"],
			"id":         data["id"],
			"matches":    data["matches"],
			"wins":       data["wins"],
			"loss":       data["loss"],
			"draw":       data["draw"],
			"points":     data["points"],
		}
		standings = append(standings, standing)

	}

	fmt.Println("Standings: ", standings)

	ctx.JSON(http.StatusAccepted, standings)
	return
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
