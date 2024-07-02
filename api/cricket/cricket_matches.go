package cricket

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CricketMatchServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCricketMatch(store *db.Store, logger *logger.Logger) *CricketMatchServer {
	return &CricketMatchServer{store: store, logger: logger}
}

func (s *CricketMatchServer) GetCricketTournamentMatchesFunc(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		fmt.Errorf("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	sports := ctx.Query("sports")

	arg := db.GetTournamentMatchParams{
		TournamentID: tournamentID,
		Sports:       sports,
	}
	matches, err := s.store.GetTournamentMatch(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	tournament, err := s.store.GetTournament(ctx, tournamentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	var matchDetails []map[string]interface{}

	for _, matchData := range matches {
		team1Name, err1 := s.store.GetClub(ctx, matchData.Team1ID)
		if err1 != nil {
			fmt.Errorf("Failed to get club details for team1: %v", err1)
			continue
		}
		team2Name, err2 := s.store.GetClub(ctx, matchData.Team2ID)
		if err2 != nil {
			fmt.Errorf("Failed to get club details for team2: %v", err2)
			continue
		}

		argTeam1 := db.GetCricketMatchScoreParams{
			MatchID: matchData.MatchID,
			TeamID:  matchData.Team1ID,
		}
		argTeam2 := db.GetCricketMatchScoreParams{
			MatchID: matchData.MatchID,
			TeamID:  matchData.Team2ID,
		}

		matchScore1, err3 := s.store.GetCricketMatchScore(ctx, argTeam1)
		if err3 != nil {
			fmt.Errorf("Failed to get score for team1: %v", err3)
			continue
		}
		matchScore2, err4 := s.store.GetCricketMatchScore(ctx, argTeam2)
		if err4 != nil {
			fmt.Errorf("Failed to get score for team2: %v", err4)
			continue
		}

		matchDetail := map[string]interface{}{
			"tournament_id":   matchData.TournamentID,
			"tournament_name": tournament.TournamentName,
			"match_id":        matchData.MatchID,
			"team1_name":      team1Name.ClubName,
			"team2_name":      team2Name.ClubName,
			"team1_score":     matchScore1.Score,
			"team2_score":     matchScore2.Score,
			"team1_wickets":   matchScore1.Wickets,
			"team2_wickets":   matchScore2.Wickets,
			"team1_extras":    matchScore1.Extras,
			"team2_extras":    matchScore1.Extras,
			"team1_innings":   matchScore1.Innings,
			"team2_innings":   matchScore1.Innings,
			"team1_overs":     matchScore1.Overs,
			"team2_overs":     matchScore1.Overs,
			"start_time":      matchData.StartTime,
			"end_time":        matchData.EndTime,
			"date_on":         matchData.DateOn,
			"sports":          matchData.Sports,
		}
		matchDetails = append(matchDetails, matchDetail)
	}
	ctx.JSON(http.StatusOK, matchDetails)
}
