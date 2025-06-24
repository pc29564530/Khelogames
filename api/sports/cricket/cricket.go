package cricket

import (
	"context"
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addCricketScoreRequest struct {
	MatchID      int64 `json:"match_id"`
	TeamID       int64 `json:"team_id"`
	InningNumber int   `json:"inning_number"`
	FollowOn     bool  `json:"follow_on"`
}

func (s *CricketServer) AddCricketScoreFunc(ctx *gin.Context) {

	var req addCricketScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	arg := db.NewCricketScoreParams{
		MatchID:           req.MatchID,
		TeamID:            req.TeamID,
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

func (s *CricketServer) GetCricketScore(matches []db.GetMatchByIDRow, tournamentID int64) []map[string]interface{} {
	ctx := context.Background()

	tournament, err := s.store.GetTournament(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
	}

	var matchDetail []map[string]interface{}
	var groupMatches []map[string]interface{}
	var knockoutRounds []map[string]interface{}

	for _, match := range matches {
		matchScore, err := s.store.GetCricketScores(ctx, match.ID)
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
			"start_timestamp": match.StartTimestamp,
			"end_timestamp":   match.EndTimestamp,
			"status_code":     match.StatusCode,
			"result":          match.Result,
			"stage":           match.Stage,
			"teams": map[string]interface{}{
				"home_team": map[string]interface{}{
					"id":         match.HomeTeamID,
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

type updateInningRequest struct {
	InningNumber int   `json:"inning_number"`
	MatchID      int64 `json:"match_id"`
	TeamID       int64 `json:"team_id"`
}

func (s *CricketServer) UpdateCricketInningsFunc(ctx *gin.Context) {
	var req updateInningRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("unable to bind the json: ", err)
		return
	}

	//get match by match id
	match, err := s.store.GetMatchByID(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to get match by match id: ", err)
		return
	}

	arg := db.UpdateCricketInningsParams{
		InningNumber: req.InningNumber,
		MatchID:      req.MatchID,
		TeamID:       req.TeamID,
	}

	response, err := s.store.UpdateCricketInnings(ctx, arg)
	if err != nil {
		s.logger.Error("unable to update the innings: ", err)
		return
	}
	if response.IsInningCompleted {
		batTeamID := response.TeamID
		var bowlTeamID int64
		if match.AwayTeamID == batTeamID {
			bowlTeamID = match.HomeTeamID
		} else {
			bowlTeamID = match.AwayTeamID
		}

		//player batting score
		playerBatsScore, err := s.store.GetCricketBatsmanScoreByTeamID(ctx, batTeamID)
		if err != nil {
			s.logger.Error("Failed to get cricket batsman score: ", err)
			return
		}

		for _, item := range *playerBatsScore {
			playerBatsmanData, err := s.store.GetPlayerCricketStatsByMatchType(ctx, item.BatsmanID)
			if err != nil {
				s.logger.Error("Failed to get the player batting stats: ", err)
				return
			}

			if playerBatsmanData == nil {
				for _, item := range *playerBatsmanData {
					if item.MatchType == *match.MatchFormat {
						_, err := s.store.AddPlayerCricketStats(ctx,
							int32(item.PlayerID),
							*&item.MatchType,
							item.Matches,
							item.BattingInnings,
							item.Runs,
							item.Balls,
							item.Fours,
							item.Sixes,
							item.Fifties,
							item.Hundreds,
							item.BestScore,
							item.BowlingInnings,
							item.Wickets,
							item.RunsConceded,
							item.BallsBowled,
							item.FourWickets,
							item.FiveWickets)
						if err != nil {
							s.logger.Error("Failed to get the player batting stats: ", err)
							return
						}
					}
				}
			} else {
				//Update the player batting stats:
				for _, item := range *playerBatsmanData {
					if item.MatchType == *match.MatchFormat {
						_, err := s.store.UpdatePlayerBattingStats(ctx,
							int32(item.PlayerID),
							*&item.MatchType,
							item.Runs,
							item.Balls,
							item.Fours,
							item.Sixes,
							item.Fifties,
							item.Hundreds,
							item.BestScore)
						if err != nil {
							s.logger.Error("Failed to get the player batting stats: ", err)
							return
						}
					}
				}
			}
		}

		//player bowling stats:
		playerBallScore, err := s.store.GetCricketBowlerScoreByTeamID(ctx, bowlTeamID)
		if err != nil {
			s.logger.Error("Failed to get cricket bowler score: ", err)
			return
		}

		for _, item := range *playerBallScore {
			playerBowlerData, err := s.store.GetPlayerCricketStatsByMatchType(ctx, item.BowlerID)
			if err != nil {
				s.logger.Error("Failed to get the player batting stats: ", err)
				return
			}
			if playerBowlerData == nil {
				for _, item := range *playerBowlerData {
					if item.MatchType == *match.MatchFormat {
						_, err := s.store.AddPlayerCricketStats(ctx,
							int32(item.PlayerID),
							*&item.MatchType,
							item.Matches,
							item.BattingInnings,
							item.Runs,
							item.Balls,
							item.Fours,
							item.Sixes,
							item.Fifties,
							item.Hundreds,
							item.BestScore,
							item.BowlingInnings,
							item.Wickets,
							item.RunsConceded,
							item.BallsBowled,
							item.FourWickets,
							item.FiveWickets)
						if err != nil {
							s.logger.Error("Failed to get the player bowling stats: ", err)
							return
						}
					}
				}
			} else {
				//Update the player bowling stats:
				for _, item := range *playerBowlerData {
					if item.MatchType == *match.MatchFormat {
						_, err := s.store.UpdatePlayerBowlingStats(ctx,
							int32(item.PlayerID),
							*&item.MatchType,
							item.Wickets,
							item.RunsConceded,
							item.BallsBowled,
							item.FourWickets,
							item.FiveWickets)
						if err != nil {
							s.logger.Error("Failed to get the player bowling stats: ", err)
							return
						}
					}
				}
			}
		}

	}

	ctx.JSON(http.StatusAccepted, response)
}

type updateCricketEndInningRequest struct {
	MatchID      int64 `json:"match_id"`
	TeamID       int64 `json:"team_id"`
	InningNumber int   `json:"inning_number"`
}

func (s *CricketServer) UpdateCricketEndInningsFunc(ctx *gin.Context) {

	var req updateCricketEndInningRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("unable to bind the json: ", err)
		return
	}

	inningResponse, batsmanResponse, bowlerResponse, err := s.store.UpdateInningEndStatus(ctx, req.MatchID, req.TeamID, req.InningNumber)
	if err != nil {
		s.logger.Error("Failed to update inning end: ", err)
	}

	batsmanPlayer, err := s.store.GetPlayer(ctx, batsmanResponse.BatsmanID)
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
	}

	bowlerPlayer, err := s.store.GetPlayer(ctx, bowlerResponse.BowlerID)
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
	}

	batsman := map[string]interface{}{
		"player":               map[string]interface{}{"id": batsmanPlayer.ID, "name": batsmanPlayer.PlayerName, "slug": batsmanPlayer.Slug, "shortName": batsmanPlayer.ShortName, "position": batsmanPlayer.Positions, "username": batsmanPlayer.Username},
		"id":                   batsmanResponse.ID,
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
		"player":            map[string]interface{}{"id": bowlerPlayer.ID, "name": bowlerPlayer.PlayerName, "slug": bowlerPlayer.Slug, "shortName": bowlerPlayer.ShortName, "position": bowlerPlayer.Positions, "username": bowlerPlayer.Username},
		"id":                bowlerResponse.ID,
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
