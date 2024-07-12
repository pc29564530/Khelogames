package clubs

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *ClubServer) GetTournamentsByClubFunc(ctx *gin.Context) {
	clubName := ctx.Query("club_name")
	response, err := s.store.GetTournamentsByClub(ctx, clubName)
	if err != nil {
		s.logger.Error("Failed to get tournament by club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetMatchByClubNameFunc(ctx *gin.Context) {
	clubIdStr := ctx.Query("id")
	clubID, err := strconv.ParseInt(clubIdStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse club id: %v", err)
		return
	}

	matches, err := s.store.GetMatchByClubName(ctx, clubID)
	if err != nil {
		s.logger.Error("Failed to get match by clubname: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	var matchDetails []map[string]interface{}

	for _, match := range matches {
		matchScoreData := s.getMatchScore(ctx, match)
		s.logger.Debug("Match Score Data: ", matchScoreData)
		s.logger.Debug("matches: ", match)
		matchDetails = append(matchDetails, matchScoreData)

	}
	ctx.JSON(http.StatusAccepted, matchDetails)
	return
}

func (s *ClubServer) getMatchScore(ctx *gin.Context, match db.GetMatchByClubNameRow) map[string]interface{} {
	clubMatchDetail := map[string]interface{}{
		"tournament_id":   match.TournamentID,
		"tournament_name": match.TournamentName,
		"match_id":        match.MatchID,
		"team1_id":        match.Team1ID,
		"team2_id":        match.Team2ID,
		"team1_name":      match.Team1Name,
		"team2_name":      match.Team2Name,
		"date_on":         match.DateOn,
		"start_time":      match.StartTime,
		"end_time":        match.EndTime,
	}

	switch match.Sports {
	case "Cricket":
		return s.getCricketMatchScore(ctx, match, clubMatchDetail)
	case "Football":
		return s.getFootballMatchScore(ctx, match, clubMatchDetail)
	default:
		s.logger.Error("Unsupported sport type:", match.Sports)
		return nil
	}
}

func (s *ClubServer) getCricketMatchScore(ctx *gin.Context, match db.GetMatchByClubNameRow, clubMatchDetail map[string]interface{}) map[string]interface{} {
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

	clubMatchDetail["team1_score"] = matchScoreData1.Score
	clubMatchDetail["team1_wickets"] = matchScoreData1.Wickets
	clubMatchDetail["team1_extras"] = matchScoreData1.Extras
	clubMatchDetail["team1_overs"] = matchScoreData1.Overs
	clubMatchDetail["team1_innings"] = matchScoreData1.Innings
	clubMatchDetail["team2_score"] = matchScoreData2.Score
	clubMatchDetail["team2_wickets"] = matchScoreData2.Wickets
	clubMatchDetail["team2_extras"] = matchScoreData2.Extras
	clubMatchDetail["team2_overs"] = matchScoreData2.Overs
	clubMatchDetail["team2_innings"] = matchScoreData2.Innings

	return clubMatchDetail
}

func (s *ClubServer) getFootballMatchScore(ctx *gin.Context, match db.GetMatchByClubNameRow, clubMatchDetail map[string]interface{}) map[string]interface{} {
	arg1 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team1ID}
	arg2 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team2ID}

	matchScoreData1, err := s.store.GetFootballMatchScore(ctx, arg1)
	if err != nil {
		s.logger.Error("Failed to get football match score for team 1:", err)
		return nil
	}
	matchScoreData2, err := s.store.GetFootballMatchScore(ctx, arg2)
	if err != nil {
		s.logger.Error("Failed to get football match score for team 2:", err)
		return nil
	}

	clubMatchDetail["team1_score"] = matchScoreData1.GoalFor
	clubMatchDetail["team2_score"] = matchScoreData2.GoalFor

	return clubMatchDetail
}
