package cricket

import (
	"context"
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addCricketScoreRequest struct {
	MatchPublicID uuid.UUID `json:"match_public_id"`
	TeamPublicID  uuid.UUID `json:"team_public_id"`
	InningNumber  int       `json:"inning_number"`
	FollowOn      bool      `json:"follow_on"`
}

func (s *CricketServer) AddCricketScoreFunc(ctx *gin.Context) {

	var req addCricketScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	match, err := s.store.GetMatchByID(ctx, req.MatchPublicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	team, err := s.store.GetTeamByPublicID(ctx, req.TeamPublicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.NewCricketScoreParams{
		MatchID:           int32(match.ID),
		TeamID:            int32(team.ID),
		InningNumber:      req.InningNumber,
		Score:             0,
		Wickets:           0,
		Overs:             0,
		RunRate:           "0.00",
		TargetRunRate:     "0.00",
		FollowOn:          req.FollowOn,
		IsInningCompleted: false,
		Declared:          false,
	}

	response, err := s.store.NewCricketScore(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusNoContent, err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return

}

func (s *CricketServer) GetCricketScore(matches []db.GetMatchByIDRow, tournamentPublicID uuid.UUID) []map[string]interface{} {
	ctx := context.Background()

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
	}

	var matchDetail []map[string]interface{}
	var groupMatches []map[string]interface{}
	var knockoutRounds []map[string]interface{}

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
					"public_id":  match.PublicID,
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
					"public_id":  match.PublicID,
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
			var roundName string
			switch *match.KnockoutLevelID {
			case 1:
				roundName = "final"
			case 2:
				roundName = "semifinal"
			case 3:
				roundName = "quaterfinal"
			case 4:
				roundName = "round_16"
			case 5:
				roundName = "round_32"
			case 6:
				roundName = "round_64"
			case 7:
				roundName = "round_128"
			}
			found := false
			for i, round := range knockoutRounds {
				if round["round"] == roundName {
					round["matches"] = append(round["matches"].([]map[string]interface{}), matchMap)
					knockoutRounds[i] = round
					found = true
					break
				}
			}
			if !found {
				knockoutRounds = append(knockoutRounds, map[string]interface{}{
					"round":   roundName,
					"matches": []map[string]interface{}{matchMap},
				})
			}
		}
	}

	matchDetail = append(matchDetail, map[string]interface{}{
		"tournament": map[string]interface{}{
			"id":              tournament.ID,
			"public_id":       tournament.PublicID,
			"name":            tournament.Name,
			"slug":            tournament.Slug,
			"country":         tournament.Country,
			"status_code":     tournament.StatusCode,
			"level":           tournament.Level,
			"start_timestamp": tournament.StartTimestamp,
			"game_id":         tournament.GameID,
			"group_count":     tournament.GroupCount,
			"max_group_team":  tournament.MaxGroupTeam,
		},
		"group_stage":    groupMatches,
		"knockout_stage": knockoutRounds,
	})
	return matchDetail
}

type updateCricketEndInningRequest struct {
	MatchPublicID uuid.UUID `json:"match_public_id"`
	TeamPublicID  uuid.UUID `json:"team_public_id"`
	InningNumber  int       `json:"inning_number"`
}

func (s *CricketServer) UpdateCricketEndInningsFunc(ctx *gin.Context) {

	var req updateCricketEndInningRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("unable to bind the json: ", err)
		return
	}

	inningResponse, batsmanResponse, bowlerResponse, err := s.store.UpdateInningEndStatusByPublicID(ctx, req.MatchPublicID, req.TeamPublicID, req.InningNumber)
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
		"player":               map[string]interface{}{"id": batsmanPlayer.ID, "name": batsmanPlayer.Name, "slug": batsmanPlayer.Slug, "shortName": batsmanPlayer.ShortName, "position": batsmanPlayer.Positions},
		"id":                   batsmanResponse.ID,
		"public_id":            batsmanResponse.PublicID,
		"match_id":             batsmanResponse.MatchID,
		"team_id":              batsmanResponse.TeamID,
		"batsman_id":           batsmanResponse.BatsmanID,
		"runs_scored":          batsmanResponse.RunsScored,
		"balls_faced":          batsmanResponse.BallsFaced,
		"fours":                batsmanResponse.Fours,
		"sixes":                batsmanResponse.Sixes,
		"batting_status":       batsmanResponse.BattingStatus,
		"is_striker":           batsmanResponse.IsStriker,
		"is_currently_batting": batsmanResponse.IsCurrentlyBatting,
	}

	bowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": bowlerPlayer.ID, "name": bowlerPlayer.Name, "slug": bowlerPlayer.Slug, "shortName": bowlerPlayer.ShortName, "position": bowlerPlayer.Positions},
		"id":                bowlerResponse.ID,
		"public_id":         bowlerResponse.PublicID,
		"match_id":          bowlerResponse.MatchID,
		"team_id":           bowlerResponse.TeamID,
		"bowler_id":         bowlerResponse.BowlerID,
		"ball":              bowlerResponse.Ball,
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
