package cricket

import (
	"context"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CricketMatchServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCricketMatchServer(store *db.Store, logger *logger.Logger) *CricketMatchServer {
	return &CricketMatchServer{store: store, logger: logger}
}

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

func (s *CricketMatchServer) AddCricketMatchScoreFunc(ctx *gin.Context) {

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

	response, err := s.store.CreateCricketMatchScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (s *CricketMatchServer) GetCricketMatchScore(matches []db.TournamentMatch, matchDetails []map[string]interface{}) []map[string]interface{} {
	ctx := context.Background()
	for _, match := range matches {
		arg1 := db.GetCricketMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team1ID}
		arg2 := db.GetCricketMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team2ID}
		matchScoreData1, err := s.store.GetCricketMatchScore(ctx, arg1)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 1:", err)
			return nil
		}
		matchScoreData2, err := s.store.GetCricketMatchScore(ctx, arg2)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for team 2:", err)
			return nil
		}
		matchDetail := map[string]interface{}{
			"team1_score":   matchScoreData1.Score,
			"team1_wickets": matchScoreData1.Wickets,
			"team1_extras":  matchScoreData1.Extras,
			"team1_overs":   matchScoreData1.Overs,
			"team1_innings": matchScoreData1.Innings,
			"team2_score":   matchScoreData2.Score,
			"team2_wickets": matchScoreData2.Wickets,
			"team2_extras":  matchScoreData2.Extras,
			"team2_overs":   matchScoreData2.Overs,
			"team2_innings": matchScoreData2.Innings,
		}
		matchDetails = append(matchDetails, matchDetail)
	}
	return matchDetails
}