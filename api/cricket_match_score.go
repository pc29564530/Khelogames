package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addCricketMatchScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
	Score        int64 `json:"score"`
	Wickets      int64 `json:"wickets"`
	Overs        int64 `json:"overs"`
	Extras       int64 `json:"extras"`
	Innings      int64 `json:"innings"`
}

func (server *Server) addCricketMatchScore(ctx *gin.Context) {

	var req addCricketMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.CreateCricketMatchScoreParams{
		MatchID:      req.MatchID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
		Score:        req.Score,
		Wickets:      req.Wickets,
		Overs:        req.Overs,
		Extras:       req.Extras,
		Innings:      req.Innings,
	}
	fmt.Println("Line no 42: ", arg)
	response, err := server.store.CreateCricketMatchScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (server *Server) getCricketMatchScore(ctx *gin.Context) {
	matchIdStr := ctx.Query("match_id")
	matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	arg := db.GetCricketMatchScoreParams{
		MatchID: matchId,
		TeamID:  teamID,
	}
	fmt.Println("Arg: ", arg)

	response, err := server.store.GetCricketMatchScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	fmt.Println("Score: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return

}

type updateCricketMatchScoreRequest struct {
	Score   int64 `json:"score"`
	Wickets int64 `json:"wickets"`
	Extras  int64 `json:"extras"`
	Innings int64 `json:"innings"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (server *Server) updateCricketMatchScore(ctx *gin.Context) {

	var req updateCricketMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("unable to get the correct data or order ", err)
		return
	}

	arg := db.UpdateCricketMatchScoreParams{
		Score:   req.Score,
		Wickets: req.Wickets,
		Extras:  req.Extras,
		Innings: req.Innings,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	response, err := server.store.UpdateCricketMatchScore(ctx, arg)
	if err != nil {
		fmt.Errorf("unable to get the response ")
		ctx.JSON(http.StatusResetContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

//cricket player score

type addCricketTeamPlayerScoreRequest struct {
	MatchID          int64  `json:"match_id"`
	TournamentID     int64  `json:"tournament_id"`
	TeamID           int64  `json:"team_id"`
	BattingOrBowling string `json:"batting_or_bowling, omit_empty"`
	Position         string `json:"position"`
	PlayerID         int64  `json:"player_id"`
	RunsScored       string `json:"runs_scored"`
	BallsFaced       string `json:"balls_faced"`
	Fours            string `json:"fours"`
	Sixes            string `json:"sixes"`
	WicketsTaken     string `json:"wickets_taken"`
	OversBowled      string `json:"overs_bowled"`
	RunsConceded     string `json:"runs_conceded"`
	WicketTakenBy    string `json:"wicket_taken_by"`
	WicketOf         string `json:"wicket_of"`
}

func (server *Server) addCricketTeamPlayerScore(ctx *gin.Context) {
	var req addCricketTeamPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Println("unable to get the err: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var position int64
	if req.Position == "" {
		position = 0
	} else {
		position, err = strconv.ParseInt(req.Position, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint position")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var runsScored int64
	if req.RunsScored == "" {
		runsScored = 0
	} else {
		runsScored, err = strconv.ParseInt(req.RunsScored, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint runsScored")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var ballsFaced int64
	if req.BallsFaced == "" {
		ballsFaced = 0
	} else {
		ballsFaced, err = strconv.ParseInt(req.BallsFaced, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint balls faced")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var fours int64
	if req.Fours == "" {
		fours = 0
	} else {
		fours, err = strconv.ParseInt(req.Fours, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint fours")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var sixes int64
	if req.Sixes == "" {
		runsScored = 0
	} else {
		sixes, err = strconv.ParseInt(req.Sixes, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint runs scored")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var wicketsTaken int64
	if req.WicketsTaken == "" {
		wicketsTaken = 0
	} else {
		wicketsTaken, err = strconv.ParseInt(req.WicketsTaken, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint wickets taken")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var runsConceded int64
	if req.RunsConceded == "" {
		runsConceded = 0
	} else {
		runsConceded, err = strconv.ParseInt(req.RunsConceded, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint runs conceded")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var wicketTakenBy int64
	if req.WicketTakenBy == "" {
		wicketTakenBy = 0
	} else {
		wicketTakenBy, err = strconv.ParseInt(req.WicketTakenBy, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint wicket taken by")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var wicketOf int64
	if req.WicketOf == "" {
		runsScored = 0
	} else {
		wicketOf, err = strconv.ParseInt(req.WicketOf, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint runsScored")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}
	arg := db.AddCricketTeamPlayerScoreParams{
		MatchID:          req.MatchID,
		TournamentID:     req.TournamentID,
		TeamID:           req.TeamID,
		BattingOrBowling: req.BattingOrBowling,
		Position:         position,
		PlayerID:         req.PlayerID,
		RunsScored:       runsScored,
		BallsFaced:       ballsFaced,
		Fours:            fours,
		Sixes:            sixes,
		WicketsTaken:     wicketsTaken,
		OversBowled:      req.OversBowled,
		RunsConceded:     runsConceded,
		WicketTakenBy:    wicketTakenBy,
		WicketOf:         wicketOf,
	}

	response, err := server.store.AddCricketTeamPlayerScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (server *Server) getCricketPlayerScore(ctx *gin.Context) {

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerIdStr := ctx.Query("player_id")
	playerID, err := strconv.ParseInt(playerIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketPlayerScoreParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
		TeamID:       teamID,
		PlayerID:     playerID,
	}

	response, err := server.store.GetCricketPlayerScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getCricketTeamPlayerScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (server *Server) getCricketTeamPlayerScore(ctx *gin.Context) {

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketTeamPlayerScoreParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
		TeamID:       teamID,
	}

	response, err := server.store.GetCricketTeamPlayerScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchScoreBattingRequest struct {
	Position      string `json:"position"`
	RunsScored    string `json:"runs_scored"`
	BallsFaced    string `json:"balls_faced"`
	Fours         string `json:"fours"`
	Sixes         string `json:"sixes"`
	WicketTakenBy string `json:"wicket_taken_by"`
	TournamentID  int64  `json:"tournament_id"`
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	PlayerID      int64  `json:"player_id"`
}

func (server *Server) updateCricketMatchScoreBatting(ctx *gin.Context) {
	var req updateCricketMatchScoreBattingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var position int64
	if req.Position == "" {
		position = 0
	} else {
		position, err = strconv.ParseInt(req.Position, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint position")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var runsScored int64
	if req.RunsScored == " " {
		runsScored = 0
	} else {

		runsScored, err = strconv.ParseInt(req.RunsScored, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint runsScored")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var ballsFaced int64
	if req.BallsFaced == "" {
		ballsFaced = 0
	} else {
		ballsFaced, err = strconv.ParseInt(req.BallsFaced, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint ballsFaced")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var fours int64
	if req.Fours == "" {
		fours = 0
	} else {
		fours, err = strconv.ParseInt(req.Fours, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint fours")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var sixes int64
	if req.Sixes == "" {
		sixes = 0
	} else {
		sixes, err = strconv.ParseInt(req.Sixes, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint sixes")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var wicketTakenBy int64
	if req.WicketTakenBy == "" {
		wicketTakenBy = 0
	} else {
		wicketTakenBy, err = strconv.ParseInt(req.WicketTakenBy, 10, 64)
		if err != nil {
			fmt.Errorf("unable to parseint wicket taken by")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	arg := db.UpdateCricketTeamPlayerScoreBattingParams{
		Position:      position,
		RunsScored:    runsScored,
		BallsFaced:    ballsFaced,
		Fours:         fours,
		Sixes:         sixes,
		WicketTakenBy: wicketTakenBy,
		TournamentID:  req.TournamentID,
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		PlayerID:      req.PlayerID,
	}

	response, err := server.store.UpdateCricketTeamPlayerScoreBatting(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketMatchScoreBowlingRequest struct {
	OversBowled  string `json:"overs_bowled"`
	RunsConceded string `json:"runs_conceded"`
	WicketsTaken string `json:"wickets_taken"`
	TournamentID int64  `json:"tournament_id"`
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
	PlayerID     int64  `json:"player_id"`
}

func (server *Server) updateCricketMatchScoreBowling(ctx *gin.Context) {
	var req updateCricketMatchScoreBowlingRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var wicketsTaken int64
	if req.WicketsTaken == "" {
		wicketsTaken = 0
	} else {
		wicketsTaken, err = strconv.ParseInt(req.WicketsTaken, 10, 64)
		if err != nil {
			fmt.Println("unable to parse: 199 ", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	var runsConceded int64
	if req.RunsConceded == "" {
		runsConceded = 0
	} else {
		runsConceded, err = strconv.ParseInt(req.RunsConceded, 10, 64)
		if err != nil {
			fmt.Println("unable to parse: 199 ", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	arg := db.UpdateCricketTeamPlayerScoreBowlingParams{
		OversBowled:  req.OversBowled,
		RunsConceded: runsConceded,
		WicketsTaken: wicketsTaken,
		TournamentID: req.TournamentID,
		MatchID:      req.MatchID,
		TeamID:       req.TeamID,
		PlayerID:     req.PlayerID,
	}

	response, err := server.store.UpdateCricketTeamPlayerScoreBowling(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
