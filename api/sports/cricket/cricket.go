package cricket

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
)

type CricketServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCricketScore(store *db.Store, logger *logger.Logger) *CricketServer {
	return &CricketServer{store: store, logger: logger}
}

// type addCricketMatchScoreRequest struct {
// 	MatchID      int64 `json:"match_id"`
// 	TournamentID int64 `json:"tournament_id"`
// 	TeamID       int64 `json:"team_id"`
// 	Score        int64 `json:"score"`
// 	Wickets      int64 `json:"wickets"`
// 	Overs        int64 `json:"overs"`
// 	Extras       int64 `json:"extras"`
// 	Innings      int64 `json:"innings"`
// }

// func (s *CricketServer) AddCricketMatchScoreFunc(ctx *gin.Context) {

// 	var req addCricketMatchScoreRequest
// 	err := ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	arg := db.CreateCricketMatchScoreParams{
// 		MatchID:      req.MatchID,
// 		TournamentID: req.TournamentID,
// 		TeamID:       req.TeamID,
// 		Score:        req.Score,
// 		Wickets:      req.Wickets,
// 		Overs:        req.Overs,
// 		Extras:       req.Extras,
// 		Innings:      req.Innings,
// 	}

// 	response, err := s.store.CreateCricketMatchScore(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusNoContent, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return

// }
