package cricket

import (
	"context"
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addCricketScoreRequest struct {
	MatchPublicID string `json:"match_public_id"`
	TeamPublicID  string `json:"team_public_id"`
	InningNumber  int    `json:"inning_number"`
	FollowOn      bool   `json:"follow_on"`
}

func (s *CricketServer) AddCricketScoreFunc(ctx *gin.Context) {

	var req addCricketScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	_, err = s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
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
		ctx.JSON(http.StatusNoContent, err)
		return
	}
	fmt.Println("Inning Line no 74: ", responseScore)

	ctx.JSON(http.StatusAccepted, gin.H{
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
	})
	return

}

func (s *CricketServer) GetCricketScore(matches []db.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
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
		matchScore, err := s.store.GetCricketScores(ctx, int32(match.ID))
		if err != nil {
			s.logger.Error("Failed to get cricket scores: ", err)
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
			"teams": map[string]interface{}{
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
			},
			"scores": map[string]interface{}{
				"home_score": homeScore,
				"away_score": awayScore,
			},
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
	MatchPublicID string `json:"match_public_id"`
	TeamPublicID  string `json:"team_public_id"`
	InningNumber  int    `json:"inning_number"`
}

func (s *CricketServer) UpdateCricketEndInningsFunc(ctx *gin.Context) {

	var req updateCricketEndInningRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("unable to bind the json: ", err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	inningResponse, batsmanResponse, bowlerResponse, err := s.store.UpdateInningEndStatusByPublicID(ctx, matchPublicID, teamPublicID, req.InningNumber)
	if err != nil {
		s.logger.Error("Failed to update inning end: ", err)
	}

	batsmanPlayer, err := s.store.GetPlayerByID(ctx, int64(batsmanResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
	}

	bowlerPlayer, err := s.store.GetPlayerByID(ctx, int64(bowlerResponse.BowlerID))
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
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
		"inning":  inningResponse,
		"batsman": batsman,
		"bowler":  bowler,
	})
}

func (s *CricketServer) GetCricketCurrentInningFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse: ", err)
		return
	}

	cricketInning, err := s.store.GetCricketCurrentInning(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get current inning and status: ", err)
		return
	}

	batTeam, err := s.store.GetTeamByID(ctx, int64(cricketInning.TeamID))
	if err != nil {
		s.logger.Error("Failed to get team by id: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"inning":       cricketInning,
		"batting_team": batTeam,
	})
}
