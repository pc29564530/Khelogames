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

type updateFootballStatisticsRequest struct {
	ShotsOnTarget   int32 `json:"shots_on_target"`
	TotalShots      int32 `json:"total_shots"`
	CornerKicks     int32 `json:"corner_kicks"`
	Fouls           int32 `json:"fouls"`
	GoalkeeperSaves int32 `json:"goalkeeper_saves"`
	FreeKicks       int32 `json:"free_kicks"`
	YellowCards     int32 `json:"yellow_cards"`
	RedCards        int32 `json:"red_cards"`
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
}

func (s *FootballServer) UpdateFootballStatisticsFunc(ctx *gin.Context) {
	var req updateFootballStatisticsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.UpdateFootballStatisticsParams{
		ShotsOnTarget:   req.ShotsOnTarget,
		TotalShots:      req.TotalShots,
		CornerKicks:     req.CornerKicks,
		Fouls:           req.Fouls,
		GoalkeeperSaves: req.GoalkeeperSaves,
		FreeKicks:       req.FreeKicks,
		YellowCards:     req.YellowCards,
		RedCards:        req.RedCards,
		MatchID:         req.MatchID,
		TeamID:          req.TeamID,
	}

	response, err := s.store.UpdateFootballStatistics(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update statistics: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
