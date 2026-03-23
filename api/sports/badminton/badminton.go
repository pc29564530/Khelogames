package badminton

import (
	"context"
	database "khelogames/database"
	errorhandler "khelogames/error_handler"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *BadmintonServer) UpdateBadmintonScoreFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `json:"match_public_id"`
		TeamPublicID  string `json:"team_public_id"`
		SetNumber     int    `json:"set_number"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	matchResult, setScore, err := s.txStore.UpdateBadmintonScoreTx(ctx, matchPublicID, teamPublicID, req.SetNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update badminton score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success":      true,
		"data":         setScore,
		"match_result": matchResult,
	})
}

func (s *BadmintonServer) GetBadmintonScoreFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	score, err := s.store.GetBadmintonMatchSetsScore(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Unable to get sets score: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch badminton score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	var sets []map[string]interface{}

	for _, item := range score {
		points, err := s.store.GetBadmintonSetsPoints(ctx, item.MatchID, item.SetNumber)
		if err != nil {
			s.logger.Error("Unable to get sets points: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to fetch set points",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		set := map[string]interface{}{
			"set_number": item.SetNumber,
			"home_score": item.HomeScore,
			"away_score": item.AwayScore,
			"set_status": item.SetStatus,
			"points":     points,
		}

		sets = append(sets, set)
	}

	response := map[string]interface{}{
		"match_public_id": matchPublicID,
		"sets":            sets,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *BadmintonServer) GetBadmintonScore(matches []database.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
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

		score, err := s.store.GetBadmintonMatchScore(ctx, match.PublicID)
		if err != nil {
			s.logger.Error("Failed to get badminton match score for home team:", err)
		}

		var hScore int
		var aScore int
		if score != nil {
			if score.HomeSetsWon != nil {
				hScore = *score.HomeSetsWon
			}
			if score.AwaySetsWon != nil {
				aScore = *score.AwaySetsWon
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

func (s *BadmintonServer) GetBadmintonSetsScore(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	score, err := s.store.GetBadmintonMatchSetsScore(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get badminton match sets score: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update badminton score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    score,
	})
}

func (s *BadmintonServer) GetBadmintonSetsPointsByTeamFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id"`
		TeamPublicID  string `uri:"team_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid match UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid team UUID format",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	match, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch team",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	matchScores, err := s.store.GetBadmintonMatchSetsScore(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get badminton match sets score: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch match sets score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	var sets []map[string]interface{}
	overAllStreak := 0
	overAllLead := 0
	overAllDeficit := 0
	overAllPoints := 0

	for _, it := range matchScores {
		currentStreak := 0
		maxStreak := 0
		maxLead := 0
		maxDeficit := 0
		setNumber := it.SetNumber

		points, err := s.store.GetBadmintonSetsPointsByTeam(ctx, it.MatchID, setNumber)
		if err != nil {
			s.logger.Error("Failed to get badminton set points: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to fetch set points",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}

		var teamWon int32
		if it.HomeScore > it.AwayScore {
			teamWon = match.HomeTeamID
		} else {
			teamWon = match.AwayTeamID
		}
		isTeamWon := teamWon == int32(team.ID)

		var pointWon int
		if int32(team.ID) == match.HomeTeamID {
			pointWon = it.HomeScore
		} else {
			pointWon = it.AwayScore
		}

		for _, point := range points {
			if point.ScoringTeamID == int32(team.ID) {
				currentStreak++
				if currentStreak > maxStreak {
					maxStreak = currentStreak
				}
			} else {
				currentStreak = 0
			}

			if int32(team.ID) == match.HomeTeamID {
				lead := point.HomeScore - point.AwayScore
				if lead > maxLead {
					maxLead = lead
				}
				deficit := point.AwayScore - point.HomeScore
				if deficit > maxDeficit {
					maxDeficit = deficit
				}
			} else {
				lead := point.AwayScore - point.HomeScore
				if lead > maxLead {
					maxLead = lead
				}
				deficit := point.HomeScore - point.AwayScore
				if deficit > maxDeficit {
					maxDeficit = deficit
				}
			}
		}

		overAllPoints = overAllPoints + pointWon
		if overAllLead < maxLead {
			overAllLead = maxLead
		}
		if overAllStreak < maxStreak {
			overAllStreak = maxStreak
		}
		if overAllDeficit < maxDeficit {
			overAllDeficit = maxDeficit
		}

		set := map[string]interface{}{
			"set_number":       setNumber,
			"point_won":        pointWon,
			"is_team_won":      isTeamWon,
			"max_streak":       maxStreak,
			"max_lead":         maxLead,
			"biggest_comeback": maxDeficit,
		}
		sets = append(sets, set)
	}

	overall := map[string]interface{}{
		"total_points_won": overAllPoints,
		"max_streak":       overAllStreak,
		"max_lead":         overAllLead,
		"biggest_comeback": overAllDeficit,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"sets":    sets,
			"overall": overall,
		},
	})
}
