package api

import (
	"fmt"
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

func (server *Server) addCricketMatchToss(ctx *gin.Context) {
	var req addCricketMatchTossRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err)
		return
	}

	arg := db.AddCricketMatchTossParams{
		TournamentID: req.TournamentID,
		MatchID:      req.MatchID,
		TossWon:      req.TossWon,
		BatOrBowl:    req.BatOrBowl,
	}

	response, err := server.store.AddCricketMatchToss(ctx, arg)
	if err != nil {
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

func (server *Server) getCricketMatchToss(ctx *gin.Context) {
	// var req getCricketMatchTossRequest
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadGateway, err)
	// 	return
	// }

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Println("Lien no 130: ", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		fmt.Println("Lien no 138: ", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}

	arg := db.GetCricketMatchTossParams{
		TournamentID: tournamentID,
		MatchID:      matchID,
	}

	fmt.Println("arg: ", arg)

	response, err := server.store.GetCricketMatchToss(ctx, arg)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}
