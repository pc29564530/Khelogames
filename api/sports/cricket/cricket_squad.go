package cricket

import (
	"khelogames/database/models"
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
	GameID     int64  `json:"game_id"`
	PlayerID   int32  `json:"player_id"`
}

type MatchSquadRequest struct {
	MatchID *int64   `json:"match_id"`
	TeamID  int64    `json:"team_id"`
	Player  []Player `json:"player"`
	OnBench []int64  `json:"on_bench"`
}

func (s *CricketServer) AddCricketSquadFunc(ctx *gin.Context) {
	var req MatchSquadRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}
	benchMap := make(map[int64]bool)
	for _, onBenchID := range req.OnBench {
		benchMap[onBenchID] = true
	}
	var cricketSquad []map[string]interface{}
	for _, player := range req.Player {
		var squad models.CricketSquad
		var err error
		isBench := benchMap[player.ID]
		squad, err = s.store.AddCricketSquad(ctx, *req.MatchID, req.TeamID, player.ID, player.Position, isBench, false)
		if err != nil {
			s.logger.Error("Failed to add football squad: ", err)
			return
		}
		// }

		cricketSquad = append(cricketSquad, map[string]interface{}{
			"id":         squad.ID,
			"match_id":   squad.MatchID,
			"team_id":    squad.TeamID,
			"player":     player,
			"role":       squad.Role,
			"on_bench":   squad.OnBench,
			"is_captain": squad.IsCaptain,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cricket squad added successfully",
		"squad":   cricketSquad,
	})
}

func (s *CricketServer) GetCricketMatchSquadFunc(ctx *gin.Context) {

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
