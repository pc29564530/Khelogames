package football

import (
	"context"
	db "khelogames/database"
	"khelogames/database/models"
	"strings"

	"github.com/google/uuid"
)

func (s *FootballServer) GetFootballScore(matches []db.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
	ctx := context.Background()

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
	}

	var matchDetail []map[string]interface{}
	groupMatches := []map[string]interface{}{}
	knockoutMatches := map[string][]map[string]interface{}{
		"final":       {},
		"semifinal":   {},
		"quaterfinal": {},
		"round_16":    {},
		"round_32":    {},
		"round_64":    {},
		"round_128":   {},
	}
	leagueMatches := []map[string]interface{}{}

	for _, match := range matches {
		homeTeamArg := db.GetFootballScoreParams{MatchID: match.ID, TeamID: int64(match.HomeTeamID)}
		awayTeamArg := db.GetFootballScoreParams{MatchID: match.ID, TeamID: int64(match.AwayTeamID)}
		homeScore, err := s.store.GetFootballScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get football match score for home team:", err)
		}
		awayScore, err := s.store.GetFootballScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get fooball match score for away team: ", err)
		}

		var emptyScore models.FootballScore
		var hScore map[string]interface{}
		if homeScore != emptyScore {
			hScore = map[string]interface{}{
				"public_id":        homeScore.PublicID,
				"first_half":       homeScore.FirstHalf,
				"second_half":      homeScore.SecondHalf,
				"goals":            homeScore.Goals,
				"penalty_shootout": homeScore.PenaltyShootOut,
			}
		}
		var aScore map[string]interface{}
		if awayScore != emptyScore {
			aScore = map[string]interface{}{
				"public_id":        awayScore.PublicID,
				"first_half":       awayScore.FirstHalf,
				"second_half":      awayScore.SecondHalf,
				"goals":            awayScore.Goals,
				"penalty_shootout": awayScore.PenaltyShootOut,
			}
		}

		game, err := s.store.GetGame(ctx, match.HomeGameID)
		if err != nil {
			s.logger.Error("Failed to get the game: ", err)
		}

		matchMap := map[string]interface{}{
			"id":                match.ID,
			"public_id":         match.PublicID,
			"homeTeam":          map[string]interface{}{"id": match.HomeTeamID, "public_id": match.HomeTeamPublicID, "name": match.HomeTeamName, "slug": match.HomeTeamSlug, "short_name": match.HomeTeamShortname, "gender": match.HomeTeamGender, "national": match.HomeTeamNational, "country": match.HomeTeamCountry, "type": match.HomeTeamType, "player_count": match.HomeTeamPlayerCount, "media_url": match.HomeTeamMediaUrl},
			"homeScore":         hScore,
			"awayTeam":          map[string]interface{}{"id": match.AwayTeamID, "public_id": match.AwayTeamPublicID, "name": match.AwayTeamName, "slug": match.AwayTeamSlug, "short_name": match.AwayTeamShortname, "gender": match.AwayTeamGender, "national": match.AwayTeamNational, "country": match.AwayTeamCountry, "type": match.AwayTeamType, "player_count": match.AwayTeamPlayerCount, "media_url": match.AwayTeamMediaUrl},
			"awayScore":         aScore,
			"start_timestamp":   match.StartTimestamp,
			"end_timestamp":     match.EndTimestamp,
			"type":              match.Type,
			"status_code":       match.StatusCode,
			"game":              game,
			"result":            match.Result,
			"stage":             match.Stage,
			"knockout_level_id": match.KnockoutLevelID,
		}

		if match.Stage == nil {
			// skip matches with no stage set
		} else if strings.EqualFold(*match.Stage, "group") {
			groupMatches = append(groupMatches, matchMap)
		} else if strings.EqualFold(*match.Stage, "knockout") {
			switch *match.KnockoutLevelID {
			case 1:
				knockoutMatches["final"] = append(knockoutMatches["final"], matchMap)
			case 2:
				knockoutMatches["semifinal"] = append(knockoutMatches["semifinal"], matchMap)
			case 3:
				knockoutMatches["quaterfinal"] = append(knockoutMatches["quaterfinal"], matchMap)
			case 4:
				knockoutMatches["round_16"] = append(knockoutMatches["round_16"], matchMap)
			case 5:
				knockoutMatches["round_32"] = append(knockoutMatches["round_32"], matchMap)
			case 6:
				knockoutMatches["round_64"] = append(knockoutMatches["round_64"], matchMap)
			case 7:
				knockoutMatches["round_128"] = append(knockoutMatches["round_128"], matchMap)
			}
		} else if strings.EqualFold(*match.Stage, "league") {
			leagueMatches = append(leagueMatches, matchMap)
		}
	}
	matchDetail = append(matchDetail, map[string]interface{}{
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
		"group_stage":    groupMatches,
		"league_stage":   leagueMatches,
		"knockout_stage": knockoutMatches,
	})

	return matchDetail
}
