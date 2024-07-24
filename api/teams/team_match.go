package teams

import (
	db "khelogames/db/sqlc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *TeamsServer) GetTournamentsByClubFunc(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	response, err := s.store.GetTournamentsByTeam(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get tournament by club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TeamsServer) GetMatchByClubNameFunc(ctx *gin.Context) {
	teamIDStr := ctx.Query("id")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse club id: %v", err)
		return
	}

	matches, err := s.store.GetMatchByTeam(ctx, teamID)
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

func (s *TeamsServer) getMatchScore(ctx *gin.Context, match db.GetMatchByTeamRow) map[string]interface{} {
	clubMatchDetail := map[string]interface{}{
		"tournament_id":   match.TournamentID,
		"tournament_name": match.TournamentName,
		"match_id":        match.MatchID,
		"team1_id":        match.HomeTeamID,
		"team2_id":        match.AwayTeamID,
		"team1_name":      match.HomeTeamName,
		"team2_name":      match.AwayTeamName,
		"start_time":      match.StartTimestamp,
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

func (s *TeamsServer) getCricketMatchScore(ctx *gin.Context, match db.GetMatchByTeamRow, clubMatchDetail map[string]interface{}) map[string]interface{} {
	arg1 := db.GetCricketMatchScoreParams{MatchID: match.MatchID, TeamID: match.HomeTeamID}
	arg2 := db.GetCricketMatchScoreParams{MatchID: match.MatchID, TeamID: match.AwayTeamID}

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

func (s *TeamsServer) getFootballMatchScore(ctx *gin.Context, match db.GetMatchByTeamRow, clubMatchDetail map[string]interface{}) map[string]interface{} {
	arg1 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.HomeTeamID}
	arg2 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.AwayTeamID}

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
