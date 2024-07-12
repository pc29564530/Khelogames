package cricket

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addCricketMatchTossRequest struct {
	TournamentID int64  `json:"tournament_id"`
	MatchID      int64  `json:"match_id"`
	TossWon      int64  `json:"toss_won"`
	BatOrBowl    string `json:"bat_or_bowl"`
}

func (s *CricketServer) AddCricketMatchTossFunc(ctx *gin.Context) {
	var req addCricketMatchTossRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	arg := db.AddCricketMatchTossParams{
		TournamentID: req.TournamentID,
		MatchID:      req.MatchID,
		TossWon:      req.TossWon,
		BatOrBowl:    req.BatOrBowl,
	}

	response, err := s.store.AddCricketMatchToss(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add cricket match toss : %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

type getCricketMatchTossRequest struct {
	TournamentID int64 `json:"tournament_id"`
	MatchID      int64 `json:"match_id"`
}

func (s *CricketServer) GetCricketMatchTossFunc(ctx *gin.Context) {
	// var req getCricketMatchTossRequest
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadGateway, err)
	// 	return
	// }

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		return
	}

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match id: %v", err)
		return
	}

	arg := db.GetCricketMatchTossParams{
		TournamentID: tournamentID,
		MatchID:      matchID,
	}

	response, err := s.store.GetCricketMatchToss(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get the cricket match toss: %v", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
