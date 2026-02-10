package cricket

import (
	"context"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addCricketScoreRequest struct {
	MatchPublicID string `json:"match_public_id" binding:"required"`
	TeamPublicID  string `json:"team_public_id" binding:"required"`
	InningNumber  int    `json:"inning_number" binding:"required,min=1"`
	FollowOn      bool   `json:"follow_on"`
}

func (s *CricketServer) AddCricketScoreFunc(ctx *gin.Context) {
	var req addCricketScoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	_, err = s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match details: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team by public id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match details: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user tournament role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check user tournament role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	if !isExists {
		s.logger.Error("User does not own this match")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You do not own this match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	arg := db.NewCricketScoreParams{
		MatchPublicID:     matchPublicID,
		TeamPublicID:      teamPublicID,
		InningNumber:      req.InningNumber,
		Score:             0,
		Wickets:           0,
		Overs:             0,
		RunRate:           "0.00",
		TargetRunRate:     "0.00",
		FollowOn:          req.FollowOn,
		IsInningCompleted: false,
		Declared:          false,
		InningStatus:      "not_started",
	}

	responseScore, err := s.store.NewCricketScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add cricket score: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add cricket score",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": gin.H{
			"inning": gin.H{
				"id":                  responseScore.ID,
				"public_id":           responseScore.PublicID,
				"match_id":            responseScore.MatchID,
				"team_id":             responseScore.TeamID,
				"inning_number":       responseScore.InningNumber,
				"score":               responseScore.Score,
				"wickets":             responseScore.Wickets,
				"overs":               responseScore.Overs,
				"run_rate":            responseScore.RunRate,
				"target_run_rate":     responseScore.TargetRunRate,
				"follow_on":           responseScore.FollowOn,
				"is_inning_completed": responseScore.IsInningCompleted,
				"declared":            responseScore.Declared,
				"inning_status":       responseScore.InningStatus,
			},
			"team": team,
		},
	})
}

func (s *CricketServer) GetCricketScore(matches []db.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
	ctx := context.Background()

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		return nil
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
		matchScore, err := s.store.GetCricketScores(ctx, int32(match.ID))
		if err != nil {
			s.logger.Error("Failed to get cricket scores: ", err)
			continue
		}
		var homeScore []models.CricketScore
		var awayScore []models.CricketScore
		for _, score := range matchScore {
			if match.HomeTeamID == score.TeamID {
				homeScore = append(homeScore, score)
			} else {
				awayScore = append(awayScore, score)
			}
		}

		matchMap := map[string]interface{}{
			"id":              match.ID,
			"public_id":       match.PublicID,
			"start_timestamp": match.StartTimestamp,
			"end_timestamp":   match.EndTimestamp,
			"status_code":     match.StatusCode,
			"result":          match.Result,
			"stage":           match.Stage,
			"home_team": map[string]interface{}{
				"id":         match.HomeTeamID,
				"public_id":  match.HomeTeamPublicID,
				"name":       match.HomeTeamName,
				"slug":       match.HomeTeamSlug,
				"short_name": match.HomeTeamShortname,
				"gender":     match.HomeTeamGender,
				"national":   match.HomeTeamNational,
				"country":    match.HomeTeamCountry,
				"type":       match.HomeTeamType,
			},
			"away_team": map[string]interface{}{
				"id":         match.AwayTeamID,
				"public_id":  match.AwayTeamPublicID,
				"name":       match.AwayTeamName,
				"slug":       match.AwayTeamSlug,
				"short_name": match.AwayTeamShortname,
				"gender":     match.AwayTeamGender,
				"national":   match.AwayTeamNational,
				"country":    match.AwayTeamCountry,
				"type":       match.AwayTeamType,
			},
			"home_score":        homeScore,
			"away_score":        awayScore,
			"knockout_level_id": match.KnockoutLevelID,
		}

		if *match.Stage == "Group" {
			groupMatches = append(groupMatches, matchMap)
		} else if match.Stage != nil && *match.Stage == "Knockout" {
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
		} else if *match.Stage == "League" {
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
			"status":          tournament.Status,
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

type updateCricketEndInningRequest struct {
	MatchPublicID string `json:"match_public_id" binding:"required"`
	TeamPublicID  string `json:"team_public_id" binding:"required"`
	InningNumber  int    `json:"inning_number" binding:"required,min=1"`
}

func (s *CricketServer) UpdateCricketEndInningsFunc(ctx *gin.Context) {
	var req updateCricketEndInningRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match details: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user tournament role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check user tournament role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	if !isExists {
		s.logger.Error("User does not own this match")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You do not own this match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	inningResponse, batsmanResponse, bowlerResponse, err := s.store.UpdateInningEndStatusByPublicID(ctx, matchPublicID, teamPublicID, req.InningNumber)
	if err != nil {
		s.logger.Error("Failed to update inning end: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update inning end",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	batsmanPlayer, err := s.store.GetPlayerByID(ctx, int64(batsmanResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	bowlerPlayer, err := s.store.GetPlayerByID(ctx, int64(bowlerResponse.BowlerID))
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get player details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	batsman := map[string]interface{}{
		"player":               map[string]interface{}{"id": batsmanPlayer.ID, "public_id": batsmanPlayer.PublicID, "name": batsmanPlayer.Name, "slug": batsmanPlayer.Slug, "shortName": batsmanPlayer.ShortName, "position": batsmanPlayer.Positions},
		"id":                   batsmanResponse.ID,
		"public_id":            batsmanResponse.PublicID,
		"match_id":             batsmanResponse.MatchID,
		"team_id":              batsmanResponse.TeamID,
		"batsman_id":           batsmanResponse.BatsmanID,
		"inning_number":        bowlerResponse.InningNumber,
		"runs_scored":          batsmanResponse.RunsScored,
		"balls_faced":          batsmanResponse.BallsFaced,
		"fours":                batsmanResponse.Fours,
		"sixes":                batsmanResponse.Sixes,
		"batting_status":       batsmanResponse.BattingStatus,
		"is_striker":           batsmanResponse.IsStriker,
		"is_currently_batting": batsmanResponse.IsCurrentlyBatting,
	}

	bowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": bowlerPlayer.ID, "public_id": bowlerPlayer.PublicID, "name": bowlerPlayer.Name, "slug": bowlerPlayer.Slug, "shortName": bowlerPlayer.ShortName, "position": bowlerPlayer.Positions},
		"id":                bowlerResponse.ID,
		"public_id":         bowlerResponse.PublicID,
		"match_id":          bowlerResponse.MatchID,
		"team_id":           bowlerResponse.TeamID,
		"bowler_id":         bowlerResponse.BowlerID,
		"inning_number":     bowlerResponse.InningNumber,
		"ball_number":       bowlerResponse.BallNumber,
		"runs":              bowlerResponse.Runs,
		"wide":              bowlerResponse.Wide,
		"no_ball":           bowlerResponse.NoBall,
		"wickets":           bowlerResponse.Wickets,
		"bowling_status":    bowlerResponse.BowlingStatus,
		"is_current_bowler": bowlerResponse.IsCurrentBowler,
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": gin.H{
			"inning":  inningResponse,
			"batsman": batsman,
			"bowler":  bowler,
		},
	})
}

func (s *CricketServer) GetCricketCurrentInningFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	cricketInning, err := s.store.GetCricketCurrentInning(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get current inning and status: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get current inning and status",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	batTeam, err := s.store.GetTeamByID(ctx, int64(cricketInning.TeamID))
	if err != nil {
		s.logger.Error("Failed to get team by id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": gin.H{
			"inning":       cricketInning,
			"batting_team": batTeam,
		},
	})
}
