package football

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FootballPlayerServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewFootballPlayer(store *db.Store, logger *logger.Logger) *FootballPlayerServer {
	return &FootballPlayerServer{store: store, logger: logger}
}

type addFootballteamPlayerScoreRequest struct {
	MatchID       int64     `json:"match_id"`
	TeamID        int64     `json:"team_id"`
	PlayerID      int64     `json:"player_id"`
	TournamentID  int64     `json:"tournament_id"`
	GoalScoreTime time.Time `json:"goal_score_time"`
}

func (server *FootballPlayerServer) addFootballGoalByPlayer(ctx *gin.Context) {

	var req addFootballteamPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		server.logger.Error("Failed to bind add football goal: %v", err)
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

	response, err := server.store.AddFootballGoalByPlayer(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to add football goal by player : %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}
