package cricket

import (
	"context"
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addCricketScoreRequest struct {
	MatchID int64  `json:"match_id"`
	TeamID  int64  `json:"team_id"`
	Inning  string `json:"inning"`
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
		Inning:            req.Inning,
		Score:             0,
		Wickets:           0,
		Overs:             0,
		RunRate:           "0.00",
		TargetRunRate:     "0.00",
		FollowOn:          false,
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
	knockoutMatches := map[string][]map[string]interface{}{
		"final":       {},
		"semifinal":   {},
		"quaterfinal": {},
		"round_16":    {},
		"round_32":    {},
		"round_64":    {},
		"round_128":   {},
	}
	var groupMatches []map[string]interface{}
	for _, match := range matches {

		homeTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.HomeTeamID}
		awayTeamArg := db.GetCricketScoreParams{MatchID: match.ID, TeamID: match.AwayTeamID}
		homeScore, err := s.store.GetCricketScore(ctx, homeTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for home team:", err)
		}
		awayScore, err := s.store.GetCricketScore(ctx, awayTeamArg)
		if err != nil {
			s.logger.Error("Failed to get cricket match score for away team:", err)
		}

		var awayScoreMap map[string]interface{}
		var homeScoreMap map[string]interface{}
		var emptyScore models.CricketScore
		if awayScore != emptyScore {
			awayScoreMap = map[string]interface{}{"id": awayScore.ID, "score": awayScore.Score, "wickets": awayScore.Wickets, "overs": awayScore.Overs, "inning": awayScore.Inning, "runRate": awayScore.RunRate, "targetRunRate": awayScore.TargetRunRate}
		}

		if homeScore != emptyScore {
			homeScoreMap = map[string]interface{}{"id": homeScore.ID, "score": homeScore.Score, "wickets": homeScore.Wickets, "overs": homeScore.Overs, "inning": homeScore.Inning, "runRate": homeScore.RunRate, "targetRunRate": homeScore.TargetRunRate}
		}

		game, err := s.store.GetGame(ctx, match.HomeGameID)
		if err != nil {
			s.logger.Error("Failed to get the game: ", err)
		}

		var inningsMap []map[string]interface{}
		matchToss, err := s.store.GetCricketToss(ctx, match.ID)
		if err != nil {
			s.logger.Error("Failed to get toss: ", err)
		}

		if matchToss.TossWin == match.HomeTeamID && matchToss.TossDecision == "Batting" {
			inningsMap = append(inningsMap, map[string]interface{}{
				"inning":              homeScore.Inning,
				"team_id":             match.HomeTeamID,
				"score":               homeScoreMap,
				"is_inning_completed": homeScoreMap["is_inning_completed"],
				"follow_on":           homeScoreMap["follow_on"],
			})
		}
		if matchToss.TossWin == match.AwayTeamID && matchToss.TossDecision == "Batting" {
			inningsMap = append(inningsMap, map[string]interface{}{
				"inning":              awayScore.Inning,
				"team_id":             match.AwayTeamID,
				"score":               awayScoreMap,
				"is_inning_completed": awayScoreMap["is_inning_completed"],
				"follow_on":           awayScoreMap["follow_on"],
			})
		}
		if matchToss.TossWin == match.HomeTeamID && matchToss.TossDecision == "Bowling" {
			inningsMap = append(inningsMap, map[string]interface{}{
				"inning":              homeScore.Inning,
				"team_id":             match.HomeTeamID,
				"score":               homeScoreMap,
				"is_inning_completed": homeScoreMap["is_inning_completed"],
				"follow_on":           homeScoreMap["follow_on"],
			})
		}
		if matchToss.TossWin == match.AwayTeamID && matchToss.TossDecision == "Bowling" {
			inningsMap = append(inningsMap, map[string]interface{}{
				"inning":              awayScore.Inning,
				"team_id":             match.AwayTeamID,
				"score":               awayScoreMap,
				"is_inning_completed": awayScoreMap["is_inning_completed"],
				"follow_on":           awayScoreMap["follow_on"],
			})
		}

		matchMap := map[string]interface{}{
			"matchId":         match.ID,
			"tournament":      map[string]interface{}{"id": tournament.ID, "name": tournament.Name, "slug": tournament.Slug, "country": tournament.Country, "sports": tournament.Sports},
			"homeTeam":        map[string]interface{}{"id": match.HomeTeamID, "name": match.HomeTeamName, "slug": match.HomeTeamSlug, "shortName": match.HomeTeamShortname, "gender": match.HomeTeamGender, "national": match.HomeTeamNational, "country": match.HomeTeamCountry, "type": match.HomeTeamType, "player_count": match.HomeTeamPlayerCount},
			"homeScore":       homeScoreMap,
			"awayTeam":        map[string]interface{}{"id": match.AwayTeamID, "name": match.AwayTeamName, "slug": match.AwayTeamSlug, "shortName": match.AwayTeamShortname, "gender": match.AwayTeamGender, "national": match.AwayTeamNational, "country": match.AwayTeamCountry, "type": match.AwayTeamType, "player_count": match.AwayTeamPlayerCount},
			"awayScore":       awayScoreMap,
			"startTimeStamp":  match.StartTimestamp,
			"endTimestamp":    match.EndTimestamp,
			"status_code":     match.StatusCode,
			"game":            game,
			"result":          match.Result,
			"stage":           match.Stage,
			"knockoutLevelId": match.KnockoutLevelID,
			"innings":         inningsMap,
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
		}
	}
	matchDetail = append(matchDetail, map[string]interface{}{
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
		"group_stage":    groupMatches,
		"knockout_stage": knockoutMatches,
	})
	return matchDetail
}

type updateInningRequest struct {
	Inning  string `json:"inning"`
	MatchID int64  `json:"match_id"`
	TeamID  int64  `json:"team_id"`
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
		Inning:  req.Inning,
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
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
			playerBatsmanData, err := s.store.GetCricketPlayerBattingStatsByMatchType(ctx, item.BatsmanID, *match.MatchFormat)
			if err != nil {
				s.logger.Error("Failed to get the player batting stats: ", err)
				return
			}
			if playerBatsmanData == nil {
				playerStats, err := s.store.AddPlayerBattingStats(ctx, int32(item.BatsmanID), *match.MatchFormat, playerBatsmanData.TotalMatches, playerBatsmanData.TotalInnings, playerBatsmanData.Runs, playerBatsmanData.Balls, playerBatsmanData.Fours, playerBatsmanData.Sixes, playerBatsmanData.Fifties, playerBatsmanData.Hundreds, playerBatsmanData.BestScore, playerBatsmanData.Average, playerBatsmanData.StrikeRate)
				if err != nil {
					s.logger.Error("Failed to get the player batting stats: ", err)
					return
				}
				fmt.Println("Player Stats: ", playerStats)
			} else {
				//Update the player batting stats:
				playerStats, err := s.store.UpdatePlayerBattingStats(ctx, int32(item.BatsmanID), *match.MatchFormat, playerBatsmanData.TotalMatches, playerBatsmanData.TotalInnings, playerBatsmanData.Runs, playerBatsmanData.Balls, playerBatsmanData.Fours, playerBatsmanData.Sixes, playerBatsmanData.Fifties, playerBatsmanData.Hundreds, playerBatsmanData.BestScore, playerBatsmanData.Average, playerBatsmanData.StrikeRate)
				if err != nil {
					s.logger.Error("Failed to get the player batting stats: ", err)
					return
				}
				fmt.Println("Player Stats: ", playerStats)
			}
		}

		//player bowling stats:
		playerBallScore, err := s.store.GetCricketBowlerScoreByTeamID(ctx, bowlTeamID)
		if err != nil {
			s.logger.Error("Failed to get cricket bowler score: ", err)
			return
		}

		for _, item := range *playerBallScore {
			playerBowlerData, err := s.store.GetCricketPlayerBowlingStatsByMatchType(ctx, item.BowlerID, *match.MatchFormat)
			if err != nil {
				s.logger.Error("Failed to get the player batting stats: ", err)
				return
			}
			if playerBowlerData == nil {
				playerStats, err := s.store.AddPlayerBowlingStats(ctx, int32(item.BowlerID), *match.MatchFormat, playerBowlerData.Matches, playerBowlerData.Innings, playerBowlerData.Wickets, playerBowlerData.Runs, playerBowlerData.Balls, playerBowlerData.Average, playerBowlerData.StrikeRate, playerBowlerData.EconomyRate, playerBowlerData.FourWickets, playerBowlerData.FiveWickets)
				if err != nil {
					s.logger.Error("Failed to get the player batting stats: ", err)
					return
				}
				fmt.Println("Player Stats: ", playerStats)
			} else {
				//Update the player batting stats:
				playerStats, err := s.store.UpdatePlayerBowlingStats(ctx, int32(item.BowlerID), *match.MatchFormat, playerBowlerData.Matches, playerBowlerData.Innings, playerBowlerData.Wickets, playerBowlerData.Runs, playerBowlerData.Balls, playerBowlerData.Average, playerBowlerData.StrikeRate, playerBowlerData.EconomyRate, playerBowlerData.FourWickets, playerBowlerData.FiveWickets)
				if err != nil {
					s.logger.Error("Failed to get the player batting stats: ", err)
					return
				}
				fmt.Println("Player Stats: ", playerStats)
			}
		}

	}

	ctx.JSON(http.StatusAccepted, response)
}

type updateCricketEndInningRequest struct {
	MatchID int64  `json:"match_id"`
	TeamID  int64  `json:"team_id"`
	Inning  string `json:"inning"`
}

func (s *CricketServer) UpdateCricketEndInningsFunc(ctx *gin.Context) {

	var req updateCricketEndInningRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("unable to bind the json: ", err)
		return
	}

	inningResponse, batsmanResponse, bowlerResponse, err := s.store.UpdateInningEndStatus(ctx, req.MatchID, req.TeamID, req.Inning)
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
