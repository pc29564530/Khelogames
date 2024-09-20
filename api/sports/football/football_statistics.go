package football

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addFootballStatisticsRequest struct {
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
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
		MatchID:         req.MatchID,
		TeamID:          req.TeamID,
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
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (s *FootballServer) GetFootballStatisticsFunc(ctx *gin.Context) {
	var req getFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.GetFootballStatisticsParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := s.store.GetFootballStatistics(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get the football statistics: ", err)
	}

	ctx.JSON(http.StatusAccepted, response)
}

type StatisticsUpdate struct {
	ShotsOnTarget   *int32
	TotalShots      *int32
	CornerKicks     *int32
	Fouls           *int32
	GoalkeeperSaves *int32
	FreeKicks       *int32
	YellowCards     *int32
	RedCards        *int32
}

func GetStatisticsUpdateFromIncident(incidentType string) StatisticsUpdate {
	switch incidentType {
	case "goal":
		return StatisticsUpdate{
			ShotsOnTarget: ptrInt32(1),
			TotalShots:    ptrInt32(1),
		}
	case "fouls":
		return StatisticsUpdate{
			Fouls:     ptrInt32(1),
			FreeKicks: ptrInt32(1),
		}
	case "yellow_cards":
		return StatisticsUpdate{
			YellowCards: ptrInt32(1),
		}
	case "red_cards":
		return StatisticsUpdate{
			RedCards: ptrInt32(1),
		}
	case "goalkeeper_saves":
		return StatisticsUpdate{
			GoalkeeperSaves: ptrInt32(1),
			ShotsOnTarget:   ptrInt32(1),
			TotalShots:      ptrInt32(1),
		}
	case "corner_kick":
		return StatisticsUpdate{
			CornerKicks: ptrInt32(1),
		}
	case "total_shots":
		return StatisticsUpdate{
			TotalShots: ptrInt32(1),
		}
	default:
		return StatisticsUpdate{}
	}
}

func ptrInt32(i int32) *int32 {
	return &i
}
