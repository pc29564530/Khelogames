package handlers

import (
	"khelogames/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getAllMatchesReq struct {
	StartTimestamp string `json:"start_timestamp"`
}

func (s *HandlersServer) GetAllMatchesFunc(ctx *gin.Context) {

	// var req getAllMatchesReq
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	s.logger.Error("Failed to bind: ", err)
	// 	return
	// }

	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}
	startDateString := ctx.Query("start_timestamp")
	startDate, err := util.ConvertTimeStamp(startDateString)
	if err != nil {
		s.logger.Error("Failed to convert to second: ", err)
	}
	response, err := s.store.GetAllMatches(ctx, int32(startDate), game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by game: ", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *HandlersServer) GetMatchByMatchIDFunc(ctx *gin.Context) {
	sport := ctx.Param("sport")
	game, err := s.store.GetGamebyName(ctx, sport)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
		return
	}
	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match_id: ", err)
		return
	}

	match, err := s.store.GetMatchByMatchID(ctx, matchID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get matches by match id: ", err)
		return
	}

	var inningsMap []map[string]interface{}
	matchToss, err := s.store.GetCricketToss(ctx, matchID)
	if err != nil {
		s.logger.Error("Failed to get toss: ", err)
	}

	if matchToss.TossWin == int64(match["home_team_id"].(float64)) && matchToss.TossDecision == "Batting" {
		inningsMap = append(inningsMap, map[string]interface{}{
			"inning":              match["homeScore"].(map[string]interface{})["inning"],
			"team_id":             int64(match["home_team_id"].(float64)),
			"score":               match["homeScore"],
			"is_inning_completed": match["homeScore"].(map[string]interface{})["is_inning_completed"],
			"follow_on":           match["homeScore"].(map[string]interface{})["follow_on"],
		})
	}
	if matchToss.TossWin == int64(match["away_team_id"].(float64)) && matchToss.TossDecision == "Batting" {
		inningsMap = append(inningsMap, map[string]interface{}{
			"inning":              match["awayScore"].(map[string]interface{})["inning"],
			"team_id":             int64(match["away_team_id"].(float64)),
			"score":               match["awayScore"],
			"is_inning_completed": match["awayScore"].(map[string]bool)["is_inning_completed"],
			"follow_on":           match["awayScore"].(map[string]bool)["follow_on"],
		})
	}
	if matchToss.TossWin != int64(match["home_team_id"].(float64)) && matchToss.TossDecision == "Bowling" {
		inningsMap = append(inningsMap, map[string]interface{}{
			"inning":              match["homeScore"].(map[string]interface{})["inning"],
			"team_id":             int64(match["home_team_id"].(float64)),
			"score":               match["homeScore"],
			"is_inning_completed": match["homeScore"].(map[string]interface{})["is_inning_completed"],
			"follow_on":           match["homeScore"].(map[string]interface{})["follow_on"],
		})
	}
	if matchToss.TossWin != int64(match["away_team_id"].(float64)) && matchToss.TossDecision == "Bowling" {
		inningsMap = append(inningsMap, map[string]interface{}{
			"inning":              match["awayScore"].(map[string]interface{})["inning"],
			"team_id":             int64(match["away_team_id"].(float64)),
			"score":               match["awayScore"],
			"is_inning_completed": match["awayScore"].(map[string]interface{})["is_inning_completed"],
			"follow_on":           match["awayScore"].(map[string]interface{})["follow_on"],
		})
	}

	match["innings"] = inningsMap

	ctx.JSON(http.StatusAccepted, match)
	return
}
