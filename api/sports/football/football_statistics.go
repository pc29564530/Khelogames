package football

import (
	db "khelogames/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addFootballStatisticsRequest struct {
	MatchID         int32 `json:"match_id"`
	TeamID          int32 `json:"team_id"`
	ShotsOnTarget   int32 `json:"shots_on_target"`
	TotalShots      int32 `json:"total_shots"`
	CornerKicks     int32 `json:"corner_kicks"`
	Fouls           int32 `json:"fouls"`
	GoalkeeperSaves int32 `json:"goalkeeper_saves"`
	FreeKicks       int32 `json:"free_kicks"`
	YellowCards     int32 `json:"yellow_cards"`
	RedCards        int32 `json:"red_cards"`
}

func (s *FootballServer) AddFootballStatisticsFunc(ctx *gin.Context) {
	var req addFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.CreateFootballStatisticsParams{
		MatchID:         int32(req.MatchID),
		TeamID:          int32(req.TeamID),
		ShotsOnTarget:   req.ShotsOnTarget,
		TotalShots:      req.TotalShots,
		CornerKicks:     req.CornerKicks,
		Fouls:           req.Fouls,
		GoalkeeperSaves: req.GoalkeeperSaves,
		FreeKicks:       req.FreeKicks,
		YellowCards:     req.YellowCards,
		RedCards:        req.RedCards,
	}

	response, err := s.store.CreateFootballStatistics(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the football statistics: ", err)
	}

	ctx.JSON(http.StatusAccepted, response)

}

type getFootballStatisticsRequest struct {
	MatchPublicID uuid.UUID `json:"match_public_id"`
	TeamPublicID  uuid.UUID `json:"team_public_id"`
}

func (s *FootballServer) GetFootballStatisticsFunc(ctx *gin.Context) {
	var req getFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	response, err := s.store.GetFootballStatistics(ctx, req.MatchPublicID, req.TeamPublicID)
	if err != nil {
		s.logger.Error("Failed to get the football statistics: ", err)
	}

	ctx.JSON(http.StatusAccepted, response)
}

type StatisticsUpdate struct {
	Penalty         int32
	ShotsOnTarget   int32
	TotalShots      int32
	CornerKicks     int32
	Fouls           int32
	GoalkeeperSaves int32
	FreeKicks       int32
	YellowCards     int32
	RedCards        int32
}

func GetStatisticsUpdateFromIncident(incidentType string) StatisticsUpdate {
	switch incidentType {
	case "goal":
		return StatisticsUpdate{
			ShotsOnTarget: 1,
			TotalShots:    1,
		}
	case "fouls":
		return StatisticsUpdate{
			Fouls:     1,
			FreeKicks: 1,
		}
	case "yellow_cards":
		return StatisticsUpdate{
			YellowCards: 1,
		}
	case "red_cards":
		return StatisticsUpdate{
			RedCards: 1,
		}
	case "goalkeeper_saves":
		return StatisticsUpdate{
			GoalkeeperSaves: 1,
			ShotsOnTarget:   1,
			TotalShots:      1,
		}
	case "corner_kicks":
		return StatisticsUpdate{
			CornerKicks: 1,
		}
	case "total_shots":
		return StatisticsUpdate{
			TotalShots: 1,
		}
	case "shots_on_target":
		return StatisticsUpdate{
			ShotsOnTarget: 1,
		}
	case "penalty":
		return StatisticsUpdate{
			ShotsOnTarget: 1,
			TotalShots:    1,
		}
	case "missed_penalty":
		return StatisticsUpdate{
			ShotsOnTarget:   1,
			TotalShots:      1,
			GoalkeeperSaves: 1,
		}
	default:
		return StatisticsUpdate{}
	}
}
