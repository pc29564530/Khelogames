package cricket

import (
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Player struct {
	ID         int64     `json:"id"`
	PublicID   uuid.UUID `json:"public_id"`
	PlayerName string    `json:"player_name"`
	ShortName  string    `json:"short_name"`
	Slug       string    `json:"slug"`
	Country    string    `json:"country"`
	Position   string    `json:"position"`
	MediaURL   string    `json:"media_url"`
	GameID     int64     `json:"game_id"`
	ProfileID  int32     `json:"profile_id"`
}

type MatchSquadRequest struct {
	MatchPublicID uuid.UUID `json:"match_public_id"`
	TeamPublicID  uuid.UUID `json:"team_public_id"`
	Player        []Player  `json:"player"`
	OnBench       []int64   `json:"on_bench"`
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
		squad, err = s.store.AddCricketSquad(ctx, req.MatchPublicID, req.TeamPublicID, player.PublicID, player.Position, isBench, false)
		if err != nil {
			s.logger.Error("Failed to add cricket squad: ", err)
			return
		}

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

	var req struct {
		MatchPublicID uuid.UUID `json: "match_public_id"`
		TeamPublicID  uuid.UUID `json:"team_public_id"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		return
	}

	cricketSquad, err := s.store.GetCricketMatchSquad(ctx, req.MatchPublicID, req.TeamPublicID)
	if err != nil {
		s.logger.Error("Failed to get cricket squad: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, cricketSquad)
}
