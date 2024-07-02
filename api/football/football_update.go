package football

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FootballUpdateServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewFootballUpdate(store *db.Store, logger *logger.Logger) *FootballUpdateServer {
	return &FootballUpdateServer{store: store, logger: logger}
}

// type addFootballMatchScoreRequest struct {
// 	MatchID      int64 `json:"match_id"`
// 	TournamentID int64 `json:"tournament_id"`
// 	TeamID       int64 `json:"team_id"`
// 	GoalFor      int64 `json:"goal_for"`
// 	GoalAgainst  int64 `json:"goal_against"`
// }

// func (server *Server) addFootballMatchScore(ctx *gin.Context) {

// 	var req addFootballMatchScoreRequest
// 	err := ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		fmt.Errorf("Failed to bind football match score: %v", err)
// 		return
// 	}

// 	// matchIdStr := ctx.Query("match_id")
// 	// matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
// 	// if err != nil {
// 	// 	fmt.Println("Line no 15: ", err)
// 	// 	ctx.Error(err)
// 	// 	return
// 	// }

// 	// tournamentIdStr := ctx.Query("tournament_id")
// 	// tournamentId, err := strconv.ParseInt(tournamentIdStr, 10, 64)
// 	// if err != nil {
// 	// 	fmt.Println("Line no 24: ", err)
// 	// 	ctx.Error(err)
// 	// 	return
// 	// }

// 	// teamIdStr := ctx.Query("team_id")
// 	// teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
// 	// if err != nil {
// 	// 	fmt.Println("Line no 32: ", err)
// 	// 	ctx.Error(err)
// 	// 	return
// 	// }

// 	// goalScoreStr := ctx.Query("goal_score")
// 	// goalScore, err := strconv.ParseInt(goalScoreStr, 10, 64)
// 	// if err != nil {
// 	// 	fmt.Println("Line no 40: ", err)
// 	// 	ctx.Error(err)
// 	// 	return
// 	// }

// 	arg := db.AddFootballMatchScoreParams{
// 		MatchID:      req.MatchID,
// 		TournamentID: req.TournamentID,
// 		TeamID:       req.TeamID,
// 		GoalFor:      req.GoalFor,
// 		GoalAgainst:  req.GoalAgainst,
// 	}

// 	response, err := s.store.AddFootballMatchScore(ctx, arg)
// 	if err != nil {
// 		fmt.Errorf("Failed to add football match score: %v", err)
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return

// }

// func (server *Server) getFootballMatchScore(ctx *gin.Context) {
// 	matchIdStr := ctx.Query("match_id")
// 	matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
// 	if err != nil {
// 		fmt.Errorf("Failed to parse match id: %v", err)
// 		return
// 	}

// 	teamIdStr := ctx.Query("team_id")
// 	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
// 	if err != nil {
// 		fmt.Errorf("Failed to parse team id: %v", err)
// 		return
// 	}

// 	arg := db.GetFootballMatchScoreParams{
// 		MatchID: matchId,
// 		TeamID:  teamID,
// 	}

// 	response, err := s.store.GetFootballMatchScore(ctx, arg)
// 	if err != nil {
// 		fmt.Errorf("Failed to get football match score: %v", err)
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return

// }

type updateFootballMatchScoreRequest struct {
	// TournamentID int64 `json:"tournament_id"`
	// MatchID      int64 `json:"match_id"`
	// TeamID       int64 `json:"team_id"`

	GoalFor     int64 `json:"goal_for"`
	GoalAgainst int64 `json:"goal_against"`
	MatchID     int64 `json:"match_id"`
	TeamID      int64 `json:"team_id"`
}

func (s *FootballUpdateServer) UpdateFootballMatchScoreFunc(ctx *gin.Context) {

	var req updateFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind update football match score: %v", err)
		return
	}

	// tournamentIdStr := ctx.Query("tournament_id")
	// matchIdStr := ctx.Query("match_id")
	// teamIdStr := ctx.Query("team_id")
	// goalScoreStr := ctx.Query("goal_score")
	// matchID, err := strconv.ParseInt(matchIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	// teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	// goalScore, err := strconv.ParseInt(goalScoreStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	// tournamentId, err := strconv.ParseInt(tournamentIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Errorf("unable to parse the match Id")
	// 	ctx.JSON(http.StatusResetContent, err)
	// 	return
	// }

	arg := db.UpdateFootballMatchScoreParams{
		GoalFor:     req.GoalFor,
		GoalAgainst: req.GoalAgainst,
		MatchID:     req.MatchID,
		TeamID:      req.TeamID,
	}

	response, err := s.store.UpdateFootballMatchScore(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to update football match score: %v", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type addFootballteamPlayerScoreRequest struct {
	MatchID       int64     `json:"match_id"`
	TeamID        int64     `json:"team_id"`
	PlayerID      int64     `json:"player_id"`
	TournamentID  int64     `json:"tournament_id"`
	GoalScoreTime time.Time `json:"goal_score_time"`
}

func (s *FootballUpdateServer) AddFootballGoalByPlayerFunc(ctx *gin.Context) {

	var req addFootballteamPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind add football goal: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.AddFootballGoalByPlayerParams{
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		PlayerID:      req.PlayerID,
		TournamentID:  req.TournamentID,
		GoalScoreTime: req.GoalScoreTime,
	}

	response, err := s.store.AddFootballGoalByPlayer(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to add football goal by player : %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}
