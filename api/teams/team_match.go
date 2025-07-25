package teams

import (
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TeamsServer) GetTournamentbyTeamFunc(ctx *gin.Context) {
	var req struct {
		TeamPublicID uuid.UUID `uri:"team_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	response, err := s.store.GetTournamentsByTeam(ctx, req.TeamPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament by team id: ", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *TeamsServer) GetMatchByTeamFunc(ctx *gin.Context) {
	var req struct {
		TeamPublicID uuid.UUID `uri:"team_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	sport := ctx.Param("sport")

	matches, err := s.store.GetMatchByTeam(ctx, req.TeamPublicID)
	if err != nil {
		s.logger.Error("Failed to get match by team id: ", err)
		return
	}
	var matchesDetails []map[string]interface{}
	clubMatchDetails := s.getMatchScore(ctx, matches, sport, matchesDetails)
	ctx.JSON(http.StatusAccepted, clubMatchDetails)
	return
}

func (s *TeamsServer) GetMatchesByTeamFunc(ctx *gin.Context) {
	var req struct {
		TeamPublicID uuid.UUID `uri:"team_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	sport := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get the game: ", err)
		return
	}

	matches, err := s.store.GetMatchesByTeam(ctx, req.TeamPublicID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by team: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, matches)
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

		homeScore, err := s.store.GetCricketScore(ctx, int32(match.MatchID), match.HomeTeamID)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for home team:", err)

		}
		awayScore, err := s.store.GetCricketScore(ctx, int32(match.MatchID), match.AwayTeamID)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for away team:", err)
		}

		homeTeam, err := s.store.GetTeamByID(ctx, int64(match.HomeTeamID))
		if err != nil {
			s.logger.Error("Failed to get home team:", err)
			return nil
		}

		awayTeam, err := s.store.GetTeamByID(ctx, int64(match.AwayTeamID))
		if err != nil {
			s.logger.Error("Failed to get away team:", err)
			return nil
		}

		tournament, err := s.store.GetTournamentByID(ctx, int64(match.TournamentID))
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
				"public_id":       tournament.PublicID,
				"user_id":         tournament.UserID,
				"name":            tournament.Name,
				"slug":            tournament.Slug,
				"country":         tournament.Country,
				"status_code":     tournament.Status,
				"level":           tournament.Level,
				"start_timestamp": tournament.StartTimestamp,
				"game_id":         tournament.GameID,
				"group_count":     tournament.GroupCount,
				"max_group_team":  tournament.MaxGroupTeam,
			},
			"homeTeam":       map[string]interface{}{"id": homeTeam.ID, "public_id": homeTeam.PublicID, "user_id": homeTeam.UserID, "name": homeTeam.Name, "slug": homeTeam.Slug, "shortName": homeTeam.Shortname, "gender": homeTeam.Gender, "national": homeTeam.National, "country": homeTeam.Country, "type": homeTeam.Type},
			"homeScore":      homeScoreMap,
			"awayTeam":       map[string]interface{}{"id": awayTeam.ID, "public_id": awayTeam.PublicID, "user_id": awayTeam.UserID, "name": awayTeam.Name, "slug": awayTeam.Slug, "shortName": awayTeam.Shortname, "gender": awayTeam.Gender, "national": awayTeam.National, "country": awayTeam.Country, "type": awayTeam.Type},
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
		homeTeam, err := s.store.GetTeamByID(ctx, int64(match.HomeTeamID))
		if err != nil {
			s.logger.Error("Failed to get home team:", err)
			return nil
		}

		awayTeam, err := s.store.GetTeamByID(ctx, int64(match.AwayTeamID))
		if err != nil {
			s.logger.Error("Failed to get away team:", err)
			return nil
		}

		tournament, err := s.store.GetTournamentByID(ctx, int64(match.TournamentID))
		if err != nil {
			s.logger.Error("Failed to get tournament: ", err)
			return nil
		}

		homeTeamArg := db.GetFootballScoreParams{MatchID: match.MatchID, TeamID: int64(match.HomeTeamID)}
		awayTeamArg := db.GetFootballScoreParams{MatchID: match.MatchID, TeamID: int64(match.AwayTeamID)}
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
				"public_id":       tournament.PublicID,
				"name":            tournament.Name,
				"slug":            tournament.Slug,
				"country":         tournament.Country,
				"status_code":     tournament.Status,
				"level":           tournament.Level,
				"start_timestamp": tournament.StartTimestamp,
				"game_id":         tournament.GameID,
				"group_count":     tournament.GroupCount,
				"max_group_team":  tournament.MaxGroupTeam,
			},
			"homeTeam":       map[string]interface{}{"id": homeTeam.ID, "public_id": homeTeam.PublicID, "name": homeTeam.Name, "slug": homeTeam.Slug, "shortName": homeTeam.Shortname, "gender": homeTeam.Gender, "national": homeTeam.National, "country": homeTeam.Country, "type": homeTeam.Type},
			"homeScore":      homeScoreMap,
			"awayTeam":       map[string]interface{}{"id": awayTeam.ID, "public_id": awayTeam.PublicID, "name": awayTeam.Name, "slug": awayTeam.Slug, "shortName": awayTeam.Shortname, "gender": awayTeam.Gender, "national": awayTeam.National, "country": awayTeam.Country, "type": awayTeam.Type},
			"awayScore":      awayScoreMap,
			"startTimeStamp": match.StartTimestamp,
			"status":         match.StatusCode,
			"type":           match.Type,
		}

		matchesDetails = append(matchesDetails, matchDetail)

	}
	return matchesDetails
}
