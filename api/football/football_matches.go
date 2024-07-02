package football

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FootballMatches struct {
	store  *db.Store
	logger *logger.Logger
}

func NewFootballMatches(store *db.Store, logger *logger.Logger) *FootballMatches {
	return &FootballMatches{store: store, logger: logger}
}

type addFootballMatchScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
	GoalFor      int64 `json:"goal_for"`
	GoalAgainst  int64 `json:"goal_against"`
}

func (s *FootballMatches) AddFootballMatchScoreFunc(ctx *gin.Context) {

	var req addFootballMatchScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind football match score: %v", err)
		return
	}

	// matchIdStr := ctx.Query("match_id")
	// matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 15: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	// tournamentIdStr := ctx.Query("tournament_id")
	// tournamentId, err := strconv.ParseInt(tournamentIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 24: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	// teamIdStr := ctx.Query("team_id")
	// teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 32: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	// goalScoreStr := ctx.Query("goal_score")
	// goalScore, err := strconv.ParseInt(goalScoreStr, 10, 64)
	// if err != nil {
	// 	fmt.Println("Line no 40: ", err)
	// 	ctx.Error(err)
	// 	return
	// }

	arg := db.AddFootballMatchScoreParams{
		MatchID:      req.MatchID,
		TournamentID: req.TournamentID,
		TeamID:       req.TeamID,
		GoalFor:      req.GoalFor,
		GoalAgainst:  req.GoalAgainst,
	}

	response, err := s.store.AddFootballMatchScore(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to add football match score: %v", err)
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (s *FootballMatches) GetFootballTournamentMatchesFunc(ctx *gin.Context) {
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	tournament, err := s.store.GetTournament(ctx, tournamentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
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

		argTeam1 := db.GetFootballMatchScoreParams{
			MatchID:      matchData.MatchID,
			TeamID:       matchData.Team1ID,
			TournamentID: matchData.TournamentID,
		}
		argTeam2 := db.GetFootballMatchScoreParams{
			MatchID:      matchData.MatchID,
			TeamID:       matchData.Team2ID,
			TournamentID: matchData.TournamentID,
		}

		matchScore1, err3 := s.store.GetFootballMatchScore(ctx, argTeam1)
		if err3 != nil {
			fmt.Errorf("Failed to get score for team1: %v", err3)
			continue
		}
		matchScore2, err4 := s.store.GetFootballMatchScore(ctx, argTeam2)
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
			"team1_score":     matchScore1.GoalFor,
			"team2_score":     matchScore2.GoalFor,
			"start_time":      matchData.StartTime,
			"end_time":        matchData.EndTime,
			"date_on":         matchData.DateOn,
			"sports":          matchData.Sports,
		}
		// fmt.Println("matchDetails: ", matchDetails)
		matchDetails = append(matchDetails, matchDetail)
	}
	ctx.JSON(http.StatusOK, matchDetails)
}
