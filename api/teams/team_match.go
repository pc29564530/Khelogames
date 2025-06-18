package teams

import (
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *TeamsServer) GetTournamentbyTeamFunc(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	response, err := s.store.GetTournamentsByTeam(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get tournament by team id: ", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TeamsServer) GetMatchByTeamFunc(ctx *gin.Context) {
	teamIDStr := ctx.Query("id")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse team id: ", err)
		return
	}

	sport := ctx.Param("sport")

	matches, err := s.store.GetMatchByTeam(ctx, teamID)
	if err != nil {
		s.logger.Error("Failed to get match by team id: ", err)
		return
	}
	var matchesDetails []map[string]interface{}
	clubMatchDetails := s.getMatchScore(ctx, matches, sport, matchesDetails)
	ctx.JSON(http.StatusAccepted, clubMatchDetails)
	return
}

func (s *TeamsServer) getMatchScore(ctx *gin.Context, matches []db.GetMatchByTeamRow, sport string, matchesDetails []map[string]interface{}) []map[string]interface{} {
	switch sport {
	case "cricket":
		return s.getCricketMatchScore(ctx, matches, matchesDetails)
	case "football":
		return s.getFootballMatchScore(ctx, matches, matchesDetails)
	default:
		s.logger.Error("Unsupported sport type:", matches)
		return nil
	}
}

func (s *TeamsServer) getCricketMatchScore(ctx *gin.Context, matches []db.GetMatchByTeamRow, matchesDetails []map[string]interface{}) []map[string]interface{} {

	for _, match := range matches {
		homeArg := db.GetCricketScoreParams{MatchID: match.MatchID, TeamID: match.HomeTeamID}
		awayArg := db.GetCricketScoreParams{MatchID: match.MatchID, TeamID: match.AwayTeamID}

		homeScore, err := s.store.GetCricketScore(ctx, homeArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for home team:", err)

		}
		awayScore, err := s.store.GetCricketScore(ctx, awayArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for away team:", err)
		}

		homeTeam, err := s.store.GetTeam(ctx, match.HomeTeamID)
		if err != nil {
			s.logger.Error("Failed to get home team:", err)
			return nil
		}

		awayTeam, err := s.store.GetTeam(ctx, match.AwayTeamID)
		if err != nil {
			s.logger.Error("Failed to get away team:", err)
			return nil
		}

		tournament, err := s.store.GetTournament(ctx, match.TournamentID)
		if err != nil {
			s.logger.Error("Failed to get tournament: ", err)
			return nil
		}

		var awayScoreMap map[string]interface{}
		var homeScoreMap map[string]interface{}
		var emptyScore models.CricketScore
		if awayScore != emptyScore {
			awayScoreMap = map[string]interface{}{"id": awayScore.ID, "score": awayScore.Score, "wickets": homeScore.Wickets, "overs": awayScore.Overs, "inning_number": awayScore.InningNumber}
		}

		if homeScore != emptyScore {
			homeScoreMap = map[string]interface{}{"id": homeScore.ID, "score": homeScore.Score, "wickets": homeScore.Wickets, "overs": homeScore.Overs, "inning_number": homeScore.InningNumber}
		}

		matchDetail := map[string]interface{}{
			"matchId": match.MatchID,
			"tournament": map[string]interface{}{
				"id":              tournament.ID,
				"name":            tournament.Name,
				"slug":            tournament.Slug,
				"country":         tournament.Country,
				"sports":          tournament.Sports,
				"status_code":     tournament.StatusCode,
				"level":           tournament.Level,
				"start_timestamp": tournament.StartTimestamp,
				"game_id":         tournament.GameID,
				"group_count":     tournament.GroupCount,
				"max_group_team":  tournament.MaxGroupTeam,
			},
			"homeTeam":       map[string]interface{}{"id": homeTeam.ID, "name": homeTeam.Name, "slug": homeTeam.Slug, "shortName": homeTeam.Shortname, "gender": homeTeam.Gender, "national": homeTeam.National, "country": homeTeam.Country, "type": homeTeam.Type},
			"homeScore":      homeScoreMap,
			"awayTeam":       map[string]interface{}{"id": awayTeam.ID, "name": awayTeam.Name, "slug": awayTeam.Slug, "shortName": awayTeam.Shortname, "gender": awayTeam.Gender, "national": awayTeam.National, "country": awayTeam.Country, "type": awayTeam.Type},
			"awayScore":      awayScoreMap,
			"startTimeStamp": match.StartTimestamp,
			"status":         match.StatusCode,
			"type":           match.Type,
		}

		matchesDetails = append(matchesDetails, matchDetail)
	}

	return matchesDetails
}

func (s *TeamsServer) getFootballMatchScore(ctx *gin.Context, matches []db.GetMatchByTeamRow, matchesDetails []map[string]interface{}) []map[string]interface{} {
	for _, match := range matches {
		homeTeam, err := s.store.GetTeam(ctx, match.HomeTeamID)
		if err != nil {
			s.logger.Error("Failed to get home team:", err)
			return nil
		}

		awayTeam, err := s.store.GetTeam(ctx, match.AwayTeamID)
		if err != nil {
			s.logger.Error("Failed to get away team:", err)
			return nil
		}

		tournament, err := s.store.GetTournament(ctx, match.TournamentID)
		if err != nil {
			s.logger.Error("Failed to get tournament: ", err)
			return nil
		}

		homeTeamArg := db.GetFootballScoreParams{MatchID: match.MatchID, TeamID: match.HomeTeamID}
		awayTeamArg := db.GetFootballScoreParams{MatchID: match.MatchID, TeamID: match.AwayTeamID}
		homeScore, err := s.store.GetFootballScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get football match score for home team:", err)
		}
		awayScore, err := s.store.GetFootballScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get fooball match score for away team: ", err)
		}

		var emptyScore models.FootballScore
		var homeScoreMap map[string]interface{}
		if homeScore != emptyScore {
			homeScoreMap = map[string]interface{}{
				"homeScore": map[string]interface{}{
					"id":         homeScore.ID,
					"score":      homeScore.Goals,
					"firstHalf":  homeScore.FirstHalf,
					"secondHalf": homeScore.SecondHalf},
			}
		}
		var awayScoreMap map[string]interface{}
		if awayScore != emptyScore {
			awayScoreMap = map[string]interface{}{
				"awayScore": map[string]interface{}{
					"id":         awayScore.ID,
					"score":      awayScore.Goals,
					"firstHalf":  awayScore.FirstHalf,
					"secondHalf": awayScore.SecondHalf,
				},
			}
		}

		matchDetail := map[string]interface{}{
			"matchId": match.MatchID,
			"tournament": map[string]interface{}{
				"id":              tournament.ID,
				"name":            tournament.Name,
				"slug":            tournament.Slug,
				"country":         tournament.Country,
				"sports":          tournament.Sports,
				"status_code":     tournament.StatusCode,
				"level":           tournament.Level,
				"start_timestamp": tournament.StartTimestamp,
				"game_id":         tournament.GameID,
				"group_count":     tournament.GroupCount,
				"max_group_team":  tournament.MaxGroupTeam,
			},
			"homeTeam":       map[string]interface{}{"id": homeTeam.ID, "name": homeTeam.Name, "slug": homeTeam.Slug, "shortName": homeTeam.Shortname, "gender": homeTeam.Gender, "national": homeTeam.National, "country": homeTeam.Country, "type": homeTeam.Type},
			"homeScore":      homeScoreMap,
			"awayTeam":       map[string]interface{}{"id": awayTeam.ID, "name": awayTeam.Name, "slug": awayTeam.Slug, "shortName": awayTeam.Shortname, "gender": awayTeam.Gender, "national": awayTeam.National, "country": awayTeam.Country, "type": awayTeam.Type},
			"awayScore":      awayScoreMap,
			"startTimeStamp": match.StartTimestamp,
			"status":         match.StatusCode,
			"type":           match.Type,
		}

		matchesDetails = append(matchesDetails, matchDetail)

	}
	return matchesDetails
}
