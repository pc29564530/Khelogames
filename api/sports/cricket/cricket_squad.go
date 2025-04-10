package cricket

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Player struct {
	ID         int64  `json:"id"`
	PlayerName string `json:"player_name"`
	ShortName  string `json:"short_name"`
	Slug       string `json:"slug"`
	Country    string `json:"country"`
	Position   string `json:"position"`
	MediaURL   string `json:"media_url"`
	Sports     string `json:"sports"`
}

type MatchSquadRequest struct {
	MatchID *int64   `json:"match_id"`
	TeamID  int64    `json:"team_id"`
	Player  []Player `json:"player"`
	Role    string   `json:"role"`
	OnBench bool     `json:"on_bench"`
}

func (s *CricketServer) AddCricketSquadFunc(ctx *gin.Context) {
	var req MatchSquadRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}
	var cricketSquad []map[string]interface{}
	fmt.Println("Player: ", req.Player)
	for _, player := range req.Player {
		squad, err := s.store.AddCricketSquad(ctx, *req.MatchID, req.TeamID, player.ID, player.Position, req.OnBench)
		if err != nil {
			s.logger.Error("Failed to add football squad: ", err)
			return
		}

		cricketSquad = append(cricketSquad, map[string]interface{}{
			"id":       squad.ID,
			"match_id": squad.MatchID,
			"team_id":  squad.TeamID,
			"player":   player,
			"role":     squad.Role,
			"on_bench": squad.OnBench,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cricket squad added successfully",
		"squad":   cricketSquad,
	})
}

func (s *CricketServer) GetCricketMatchSquadFunc(ctx *gin.Context) {
	// var req struct {
	// 	MatchID *int64 `json:"match_id"`
	// 	TeamID  int64  `json:"team_id"`
	// }

	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	s.logger.Error("failed to bind: ", err)
	// 	return
	// }

	teamIDString := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse team id: ", err)
		return
	}

	matchIDString := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match id: ", err)
		return
	}

	cricketSquad, err := s.store.GetCricketMatchSquad(ctx, matchID, teamID)
	if err != nil {
		s.logger.Error("Failed to get cricket squad: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, cricketSquad)
}
