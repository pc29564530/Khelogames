package cricket

import (
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addCricketBatScore struct {
	MatchPublicID      string `json:"match_public_id"`
	TeamPublicID       string `json:"team_public_id"`
	BatsmanPublicID    string `json:"batsman_public_id"`
	Position           string `json:"position"`
	RunsScored         int32  `json:"runs_scored"`
	BallsFaced         int32  `json:"balls_faced"`
	Fours              int32  `json:"fours"`
	Sixes              int32  `json:"sixes"`
	BattingStatus      bool   `json:"batting_status"`
	IsStriker          bool   `json:"is_striker"`
	IsCurrentlyBatting bool   `json:"is_currently_batting"`
	InningNumber       int    `json:"inning_number"`
}

func (s *CricketServer) AddCricketBatScoreFunc(ctx *gin.Context) {
	var req addCricketBatScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind player batting score: ", err)
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

	batsmanPublicID, err := uuid.Parse(req.BatsmanPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	strickerResponse, err := s.store.GetCricketStricker(ctx, matchPublicID, teamPublicID, req.InningNumber)
	if err != nil {
		s.logger.Error("failed to get stricker: ", err)
		return
	}

	arg := db.AddCricketBatsScoreParams{
		MatchPublicID:      matchPublicID,
		TeamPublicID:       teamPublicID,
		BatsmanPublicID:    batsmanPublicID,
		InningNumber:       req.InningNumber,
		Position:           req.Position,
		RunsScored:         req.RunsScored,
		BallsFaced:         req.BallsFaced,
		Fours:              req.Fours,
		Sixes:              req.Sixes,
		BattingStatus:      req.BattingStatus,
		IsStriker:          req.IsStriker,
		IsCurrentlyBatting: req.IsCurrentlyBatting,
	}

	if strickerResponse != nil {
		arg.IsStriker = false
	} else {
		arg.IsStriker = true
	}

	response, err := s.store.AddCricketBatsScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket player score: ", gin.H{"error": err.Error()})
		return
	}

	playerData, err := s.store.GetPlayerByPublicID(ctx, batsmanPublicID)
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	batsman := map[string]interface{}{
		"id":                   response.ID,
		"public_id":            response.PublicID,
		"player":               map[string]interface{}{"id": playerData.ID, "public_id": playerData.PublicID, "name": playerData.Name, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions},
		"team_id":              response.TeamID,
		"match_id":             response.MatchID,
		"batsman_id":           response.BatsmanID,
		"runsScored":           response.RunsScored,
		"ballFaced":            response.BallsFaced,
		"fours":                response.Fours,
		"sixes":                response.Sixes,
		"batting_status":       response.BattingStatus,
		"is_striker":           response.IsStriker,
		"is_currently_batting": response.IsCurrentlyBatting,
		"inning_number":        response.InningNumber,
	}

	ctx.JSON(http.StatusAccepted, batsman)
	return
}

type addCricketBallScore struct {
	MatchPublicID      string `json:"match_public_id"`
	TeamPublicID       string `json:"team_public_id"`
	BowlerPublicID     string `json:"bowler_public_id"`
	PrevBowlerPublicID string `json:"prev_bowler_public_id"`
	BallNumber         int32  `json:"ball_number"`
	Runs               int32  `json:"runs"`
	Wickets            int32  `json:"wickets"`
	Wide               int32  `json:"wide"`
	NoBall             int32  `json:"no_ball"`
	InningNumber       int    `json:"inning_number"`
}

func (s *CricketServer) AddCricketBallFunc(ctx *gin.Context) {
	var req addCricketBallScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	bowlerPublicID, err := uuid.Parse(req.BowlerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var prevBowlerPublicID uuid.UUID
	if req.PrevBowlerPublicID != "" {
		prevBowlerPublicID, err = uuid.Parse(req.PrevBowlerPublicID)
		if err != nil {
			s.logger.Error("Invalid UUID format", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin Transcation")
	}

	defer tx.Rollback()

	var prevBowlerID string
	var currentBowlerResponse *models.BowlerScore
	var prevBowler map[string]interface{}

	if req.PrevBowlerPublicID != prevBowlerID {
		currentBowlerResponse, err = s.store.UpdateBowlingBowlerStatus(ctx, matchPublicID, bowlerPublicID, prevBowlerPublicID, req.InningNumber)
		if err != nil {
			s.logger.Error("Failed to update current bowler status: ", err)
			return
		}

		playerData, err := s.store.GetPlayerByPublicID(ctx, bowlerPublicID)
		if err != nil {
			s.logger.Error("Failed to get Player: ", err)
		}
		prevBowler = map[string]interface{}{
			"player":            map[string]interface{}{"id": playerData.ID, "public_id": playerData.PublicID, "name": playerData.Name, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions},
			"id":                currentBowlerResponse.ID,
			"public_id":         currentBowlerResponse.PublicID,
			"match_id":          currentBowlerResponse.MatchID,
			"team_id":           currentBowlerResponse.TeamID,
			"bowler_id":         currentBowlerResponse.BowlerID,
			"runs":              currentBowlerResponse.Runs,
			"inning_number":     currentBowlerResponse.InningNumber,
			"ball_number":       currentBowlerResponse.BallNumber,
			"wide":              currentBowlerResponse.Wide,
			"no_ball":           currentBowlerResponse.NoBall,
			"wickets":           currentBowlerResponse.Wickets,
			"bowling_status":    currentBowlerResponse.BowlingStatus,
			"is_current_bowler": currentBowlerResponse.IsCurrentBowler,
		}
	}

	arg := db.AddCricketBallParams{
		MatchPublicID:  matchPublicID,
		TeamPublicID:   teamPublicID,
		BowlerPublicID: bowlerPublicID,
		InningNumber:   req.InningNumber,
		BallNumber:     req.BallNumber,
		Runs:           req.Runs,
		Wickets:        req.Wickets,
		Wide:           req.Wide,
		NoBall:         req.NoBall,
	}

	response, err := s.store.AddCricketBall(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket bowler data: ", gin.H{"error": err.Error()})
		return
	}

	playerData, err := s.store.GetPlayerByPublicID(ctx, bowlerPublicID)
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
	}

	currentBowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": playerData.ID, "public_id": playerData.PublicID, "name": playerData.Name, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions},
		"id":                response.ID,
		"public_id":         response.PublicID,
		"match_id":          response.MatchID,
		"team_id":           response.TeamID,
		"bowler_id":         response.BowlerID,
		"runs":              response.Runs,
		"ball_number":       response.BallNumber,
		"wide":              response.Wide,
		"no_ball":           response.NoBall,
		"wickets":           response.Wickets,
		"bowling_status":    response.BowlingStatus,
		"is_current_bowler": response.IsCurrentBowler,
		"inning_number":     response.InningNumber,
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"current_bowler": prevBowler,
		"next_bowler":    currentBowler,
	})
}

func (s *CricketServer) GetPlayerScoreFunc(ctx *gin.Context) {

	matchPublicIDString := ctx.Query("match_public_id")
	teamPublicIDString := ctx.Query("team_public_id")
	gameName := ctx.Param("sport")
	matchPublicID, err := uuid.Parse(matchPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parse match id ", err)
		return
	}

	teamPublicID, err := uuid.Parse(teamPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parse team id ", err)
		return
	}

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get the game: ", gin.H{"error": err.Error()})
		return
	}

	teamPlayerScore, err := s.store.GetCricketPlayersScore(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get players score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(teamPlayerScore) == 0 {
		s.logger.Warn("No players score found for match and team")
		ctx.JSON(http.StatusOK, gin.H{
			"battingTeam": nil,
			"innings":     nil,
		})
		return
	}

	_, err = s.store.GetMatchByPublicId(ctx, matchPublicID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get match:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playerOut, err := s.store.GetCricketWickets(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get wicket: ", err)
		return
	}

	// var battingTeamId int64
	// fmt.Println("BattingTeamID: ", battingTeamId)
	// var bowlingTeamId int64
	// if int64(teamPlayerScore[0].TeamID) == int64(match["home_team_id"].(float64)) {
	// 	battingTeamId = int64(teamPlayerScore[0].TeamID)
	// 	bowlingTeamId = int64(match["away_team_id"].(float64))
	// } else {
	// 	battingTeamId = int64(teamPlayerScore[0].TeamID)
	// 	bowlingTeamId = int64(match["home_team_id"].(float64))
	// }

	battingTeam, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get batting team : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	inningData := make(map[string][]map[string]interface{})
	for _, playerScore := range teamPlayerScore {
		playerData, err := s.store.GetPlayerByID(ctx, int64(playerScore.BatsmanID))
		if err != nil {
			s.logger.Error("Failed to get players data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		playerDetails := map[string]interface{}{
			"player":               map[string]interface{}{"id": playerData.ID, "public_id": playerData.PublicID, "name": playerData.Name, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions},
			"id":                   playerScore.ID,
			"public_id":            playerData.PublicID,
			"match_id":             playerScore.MatchID,
			"team_id":              playerScore.TeamID,
			"batsman_id":           playerScore.BatsmanID,
			"runs_scored":          playerScore.RunsScored,
			"balls_faced":          playerScore.BallsFaced,
			"fours":                playerScore.Fours,
			"sixes":                playerScore.Sixes,
			"batting_status":       playerScore.BattingStatus,
			"is_striker":           playerScore.IsStriker,
			"is_currently_batting": playerScore.IsCurrentlyBatting,
			"inning_number":        playerScore.InningNumber,
		}

		for _, item := range playerOut {
			if int64(item.BatsmanID) == playerData.ID {
				bowlerData, err := s.store.GetPlayerByID(ctx, int64(item.BowlerID))
				if err != nil {
					s.logger.Error("Failed to get bowler data : ", err)
					ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				// Add wicket details to the same playerDetails map
				playerDetails["wicket_type"] = item.WicketType
				playerDetails["bowler_name"] = bowlerData.Name
				break
			}
		}
		playerScoreInningNumber := strconv.Itoa(playerScore.InningNumber)
		inningData[playerScoreInningNumber] = append(inningData[playerScoreInningNumber], playerDetails)

	}

	var scoreDetails map[string]interface{}
	var emptyDetails map[string]interface{}
	if len(inningData) >= 1 {
		scoreDetails = map[string]interface{}{
			"battingTeam": map[string]interface{}{"id": battingTeam.ID, "public_id": battingTeam.PublicID, "name": battingTeam.Name, "slug": battingTeam.Slug, "shortName": battingTeam.Shortname, "gender": battingTeam.Gender, "national": battingTeam.National, "country": battingTeam.Country, "type": battingTeam.Type},
			"innings":     inningData,
		}
	} else {
		scoreDetails = map[string]interface{}{
			"battingTeam": map[string]interface{}{"id": battingTeam.ID, "public_id": battingTeam.PublicID, "name": battingTeam.Name, "slug": battingTeam.Slug, "shortName": battingTeam.Shortname, "gender": battingTeam.Gender, "national": battingTeam.National, "country": battingTeam.Country, "type": battingTeam.Type},
			"innings":     emptyDetails,
		}
	}

	ctx.JSON(http.StatusAccepted, scoreDetails)
}

type getCricketBowlersRequest struct {
	MatchPublicID string `json:"match_public_id"`
	TeamPublicID  string `json:"team_public_id"`
}

func (s *CricketServer) GetCricketBowlerFunc(ctx *gin.Context) {
	matchPublicIDString := ctx.Query("match_public_id")
	teamPublicIDString := ctx.Query("team_public_id")

	matchPublicID, err := uuid.Parse(matchPublicIDString)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	teamPublicID, err := uuid.Parse(teamPublicIDString)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get the game: ", gin.H{"error": err.Error()})
		return
	}

	playerScore, err := s.store.GetCricketBalls(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get cricket bowler data : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := s.store.GetMatchByPublicId(ctx, matchPublicID, game.ID)
	if err != nil {
		s.logger.Error("Failed to get player :", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var battingTeamId int64
	var bowlingTeamId int64

	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)

	if team.ID == int64(match["home_team_id"].(float64)) {
		battingTeamId = int64(match["away_team_id"].(float64))
		bowlingTeamId = team.ID
	} else {
		battingTeamId = int64(match["home_team_id"].(float64))
		bowlingTeamId = team.ID
	}
	fmt.Println("Batting Team ID: ", battingTeamId)

	bowlingTeam, err := s.store.GetTeamByID(ctx, bowlingTeamId)
	if err != nil {
		s.logger.Error("Failed to get players score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inningData := make(map[string][]map[string]interface{})
	for _, playerScore := range playerScore {
		playerData, err := s.store.GetPlayerByID(ctx, int64(playerScore.BowlerID))
		if err != nil {
			s.logger.Error("Failed to get players data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		bowlingDetails := map[string]interface{}{
			"player":            map[string]interface{}{"id": playerData.ID, "public_id": playerData.PublicID, "name": playerData.Name, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions},
			"id":                playerScore.ID,
			"public_id":         playerData.PublicID,
			"match_id":          playerScore.MatchID,
			"team_id":           playerScore.TeamID,
			"bowler_id":         playerScore.BowlerID,
			"ball_number":       playerScore.BallNumber,
			"runs":              playerScore.Runs,
			"wide":              playerScore.Wide,
			"no_ball":           playerScore.NoBall,
			"wickets":           playerScore.Wickets,
			"bowling_status":    playerScore.BowlingStatus,
			"is_current_bowler": playerScore.IsCurrentBowler,
			"inning_number":     playerScore.InningNumber,
		}
		playerScoreInningNumber := strconv.Itoa(playerScore.InningNumber)
		inningData[playerScoreInningNumber] = append(inningData[playerScoreInningNumber], bowlingDetails)
	}

	scoreDetails := map[string]interface{}{
		"bowlingTeam": map[string]interface{}{
			"id":           bowlingTeam.ID,
			"public_id":    bowlingTeam.PublicID,
			"name":         bowlingTeam.Name,
			"slug":         bowlingTeam.Slug,
			"shortName":    bowlingTeam.Shortname,
			"gender":       bowlingTeam.Gender,
			"national":     bowlingTeam.National,
			"country":      bowlingTeam.Country,
			"type":         bowlingTeam.Type,
			"player_count": bowlingTeam.PlayerCount,
		},
		"innings": inningData,
	}

	ctx.JSON(http.StatusAccepted, scoreDetails)
}

func (s *CricketServer) GetCricketWicketsFunc(ctx *gin.Context) {

	matchPublicIDString := ctx.Query("match_public_id")
	teamPublicIDString := ctx.Query("team_public_id")

	matchPublicID, err := uuid.Parse(matchPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parse to int: ", err)
	}

	teamPublicID, err := uuid.Parse(teamPublicIDString)
	if err != nil {
		s.logger.Error("Failed to parse to int: ", err)
	}
	wicketsResponse, err := s.store.GetCricketWickets(ctx, matchPublicID, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get cricket bowler score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Successfully get the wickets: ", wicketsResponse)

	var wicketsData []map[string]interface{}

	match, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team, err := s.store.GetTeamByPublicID(ctx, teamPublicID)
	if err != nil {
		s.logger.Error("Failed to get team : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = s.store.GetCricketScore(ctx, int32(match.ID), int32(team.ID))
	if err != nil {
		s.logger.Error("Failed to get current score data : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	inningData := make(map[string][]map[string]interface{})

	for _, wicket := range wicketsResponse {
		batsmanPlayerData, err := s.store.GetPlayerByID(ctx, int64(wicket.BatsmanID))
		if err != nil {
			s.logger.Error("Failed to get batsman data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		bowlerPlayerData, err := s.store.GetPlayerByID(ctx, int64(wicket.BowlerID))
		if err != nil {
			s.logger.Error("Failed to get bowler data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var fielderData *models.Player
		if wicket.FielderID != nil {
			fielderData, err = s.store.GetPlayerByID(ctx, int64(*wicket.FielderID))
			if err != nil {
				s.logger.Error("Failed to get fielder data : ", err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}

		wicketData := map[string]interface{}{
			"batsman_player": map[string]interface{}{"id": batsmanPlayerData.ID, "public_id": batsmanPlayerData.PublicID, "name": batsmanPlayerData.Name, "slug": batsmanPlayerData.Slug, "shortName": batsmanPlayerData.ShortName, "position": batsmanPlayerData.Positions},
			"bowler_player":  map[string]interface{}{"id": bowlerPlayerData.ID, "public_id": bowlerPlayerData.PublicID, "name": bowlerPlayerData.Name, "slug": bowlerPlayerData.Slug, "shortName": bowlerPlayerData.ShortName, "position": bowlerPlayerData.Positions},
			"fielder_player": fielderData,
			"id":             wicket.ID,
			"public_id":      wicket.PublicID,
			"match_id":       wicket.MatchID,
			"team_id":        wicket.TeamID,
			"batsman_id":     wicket.BatsmanID,
			"bowler_id":      wicket.BowlerID,
			"inning_number":  wicket.InningNumber,
			"wicket_number":  wicket.WicketsNumber,
			"wicket_type":    wicket.WicketType,
			"ball_number":    wicket.BallNumber,
			"fielder_id":     wicket.FielderID,
			"score":          wicket.Score,
		}

		wicketsData = append(wicketsData, wicketData)
		wicketsInningNumber := strconv.Itoa(wicket.InningNumber)
		inningData[wicketsInningNumber] = append(inningData[wicketsInningNumber], wicketData)
	}

	s.logger.Debug("Successfully update the wickets: ", wicketsData)

	ctx.JSON(http.StatusAccepted, gin.H{
		"inning": inningData,
	})
}

type updateCricketBatsmanScoreRequest struct {
	BatsmanID  int64  `json:"batsman_id"`
	TeamID     int64  `json:"team_id"`
	MatchID    int64  `json:"match_id"`
	Position   string `json:"position"`
	RunsScored int32  `json:"runs_scored"`
	Fours      int32  `json:"fours"`
	Sixes      int32  `json:"sixes"`
}

type updateCricketPlayerStatsRequest struct {
	BatsmanID   int64  `json:"batsman_id"`
	BowlerID    int64  `json:"bowler_id"`
	MatchID     int64  `json:"match_id"`
	Position    string `json:"position"`
	RunsScored  int32  `json:"runs_scored"`
	BowlerBalls int32  `json:"bowler_balls"`
	Fours       int32  `json:"fours"`
	Sixes       int32  `json:"sixes"`
}

// type updateWideRunsRequest struct {
// 	MatchPublicID       string `json:"match_public_id"`
// 	BatsmanPublicID     string `json:"batsman_public_id"`
// 	BowlerPublicID      string `json:"bowler_public_id"`
// 	BattingTeamPublicID string `json:"batting_team_public_id"`
// 	RunsScored          int32  `json:"runs_scored"`
// 	InningNumber        int    `json:"inning_number"`
// }

func (s *CricketServer) UpdateWideBallWS(ctx *gin.Context, message map[string]interface{}) (data map[string]interface{}) {
	// var req updateWideRunsRequest
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	s.logger.Error("Failed to bind: ", err)
	// 	return
	// }

	matchPublicID, err := uuid.Parse(message["match_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	batsmanPublicID, err := uuid.Parse(message["batsman_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	bowlerPublicID, err := uuid.Parse(message["bowler_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	battingTeamPublicID, err := uuid.Parse(message["batting_team_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	runsScored := int32(message["runs_scored"].(float64))
	inningNumber := int(message["inning_number"].(float64))

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction: ", err)
		return
	}

	defer tx.Rollback()

	batsmanResponse, bowlerResponse, inningScore, err := s.store.UpdateWideRuns(ctx, matchPublicID, battingTeamPublicID, bowlerPublicID, runsScored, inningNumber)
	if err != nil {
		s.logger.Error("Failed to update wide: ", err)
		return
	}

	var currentBatsman []models.BatsmanScore
	var nonStrikerResponse models.BatsmanScore
	if bowlerResponse.BallNumber%6 == 0 && runsScored%2 == 0 {
		currentBatsman, err = s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
	} else if bowlerResponse.BallNumber%6 != 0 && runsScored%2 != 0 {
		currentBatsman, err = s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
	}

	verifyBatsman, err := s.store.GetPlayerByPublicID(ctx, batsmanPublicID)
	if err != nil {
		s.logger.Error("Failed to update stricker: ", err)
	}

	for _, curBatsman := range currentBatsman {
		if curBatsman.BatsmanID == int32(verifyBatsman.ID) && curBatsman.IsStriker {
			batsmanResponse.IsStriker = curBatsman.IsStriker
		} else if curBatsman.BatsmanID != int32(verifyBatsman.ID) && curBatsman.IsStriker {
			nonStrikerResponse = curBatsman
		} else if curBatsman.BatsmanID == int32(verifyBatsman.ID) && !curBatsman.IsStriker {
			batsmanResponse.IsStriker = curBatsman.IsStriker
		} else if curBatsman.BatsmanID != int32(verifyBatsman.ID) && !curBatsman.IsStriker {
			nonStrikerResponse = curBatsman
		}
	}

	batsmanPlayerData, err := s.store.GetPlayerByID(ctx, int64(batsmanResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	nonStrikerPlayerData, err := s.store.GetPlayerByID(ctx, int64(nonStrikerResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	bowlerPlayerData, err := s.store.GetPlayerByID(ctx, int64(bowlerResponse.BowlerID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	batsman := map[string]interface{}{
		"player":               map[string]interface{}{"id": batsmanPlayerData.ID, "public_id": batsmanPlayerData.PublicID, "name": batsmanPlayerData.Name, "slug": batsmanPlayerData.Slug, "shortName": batsmanPlayerData.ShortName, "position": batsmanPlayerData.Positions},
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

	var emptyBatsman models.BatsmanScore
	var nonStriker map[string]interface{}

	if nonStrikerResponse != emptyBatsman {
		nonStriker = map[string]interface{}{
			"player":               map[string]interface{}{"id": nonStrikerPlayerData.ID, "public_id": nonStrikerPlayerData.PublicID, "name": nonStrikerPlayerData.Name, "slug": nonStrikerPlayerData.Slug, "shortName": nonStrikerPlayerData.ShortName, "position": nonStrikerPlayerData.Positions},
			"id":                   nonStrikerResponse.ID,
			"public_id":            nonStrikerPlayerData.PublicID,
			"match_id":             nonStrikerResponse.MatchID,
			"team_id":              nonStrikerResponse.TeamID,
			"batsman_id":           nonStrikerResponse.BatsmanID,
			"runs_scored":          nonStrikerResponse.RunsScored,
			"balls_faced":          nonStrikerResponse.BallsFaced,
			"fours":                nonStrikerResponse.Fours,
			"sixes":                nonStrikerResponse.Sixes,
			"batting_status":       nonStrikerResponse.BattingStatus,
			"is_striker":           nonStrikerResponse.IsStriker,
			"is_currently_batting": nonStrikerResponse.IsCurrentlyBatting,
		}
	}

	bowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": bowlerPlayerData.ID, "public_id": bowlerPlayerData.PublicID, "name": bowlerPlayerData.Name, "slug": bowlerPlayerData.Slug, "shortName": bowlerPlayerData.ShortName, "position": bowlerPlayerData.Positions},
		"id":                bowlerResponse.ID,
		"public_id":         bowlerResponse.PublicID,
		"match_id":          bowlerResponse.MatchID,
		"team_id":           bowlerResponse.TeamID,
		"bowler_id":         bowlerResponse.BowlerID,
		"ball_number":       bowlerResponse.BallNumber,
		"runs":              bowlerResponse.Runs,
		"wide":              bowlerResponse.Wide,
		"no_ball":           bowlerResponse.NoBall,
		"wickets":           bowlerResponse.Wickets,
		"bowling_status":    bowlerResponse.BowlingStatus,
		"is_current_bowler": bowlerResponse.IsCurrentBowler,
	}

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"striker_batsman":     batsman,
	// 	"non_striker_batsman": nonStriker,
	// 	"bowler":              bowler,
	// 	"inning_score":        inningScore,
	// })

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
		return
	}

	data = map[string]interface{}{
		"type": "UPDATE_SCORE",
		"payload": map[string]interface{}{
			"striker_batsman":     batsman,
			"non_striker_batsman": nonStriker,
			"bowler":              bowler,
			"inning_score":        inningScore,
			"event_type":          "wide",
		},
	}
	return data
}

// type updateNoBallRuns struct {
// 	MatchPublicID       string `json:"match_public_id"`
// 	BatsmanPublicID     string `json:"batsman_public_id"`
// 	BowlerPublicID      string `json:"bowler_public_id"`
// 	BattingTeamPublicID string `json:"batting_team_public_id"`
// 	RunsScored          int32  `json:"runs_scored"`
// 	InningNumber        int    `json:"inning_number"`
// }

func (s *CricketServer) UpdateNoBallsRunsWS(ctx *gin.Context, message map[string]interface{}) (data map[string]interface{}) {
	// var req updateNoBallRuns
	// err := ctx.ShouldBindJSON(&req)
	// if err != nil {
	// 	s.logger.Error("Failed to bind: ", err)
	// 	return
	// }

	matchPublicID, err := uuid.Parse(message["match_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	batsmanPublicID, err := uuid.Parse(message["batsman_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	bowlerPublicID, err := uuid.Parse(message["bowler_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	battingTeamPublicID, err := uuid.Parse(message["batting_team_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	runsScored := int32(message["runs_scored"].(float64))
	inningNumber := int(message["inning_number"].(float64))

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction: ", err)
		return
	}

	defer tx.Rollback()

	batsmanResponse, bowlerResponse, inningScore, err := s.store.UpdateNoBallsRuns(ctx, matchPublicID, bowlerPublicID, battingTeamPublicID, runsScored, inningNumber)
	if err != nil {
		s.logger.Error("Failed to update no_ball: ", err)
		return
	}

	var currentBatsman []models.BatsmanScore
	var nonStrikerResponse models.BatsmanScore
	if bowlerResponse.BallNumber%6 == 0 && runsScored%2 == 0 {
		currentBatsman, err = s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
	} else if bowlerResponse.BallNumber%6 != 0 && runsScored%2 != 0 {
		currentBatsman, err = s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
	}

	verifyBatsman, err := s.store.GetPlayerByPublicID(ctx, batsmanPublicID)
	if err != nil {
		s.logger.Error("Failed to update stricker: ", err)
	}

	for _, curBatsman := range currentBatsman {
		if curBatsman.BatsmanID == int32(verifyBatsman.ID) && curBatsman.IsStriker {
			batsmanResponse.IsStriker = curBatsman.IsStriker
		} else if curBatsman.BatsmanID != int32(verifyBatsman.ID) && curBatsman.IsStriker {
			nonStrikerResponse = curBatsman
		} else if curBatsman.BatsmanID == int32(verifyBatsman.ID) && !curBatsman.IsStriker {
			batsmanResponse.IsStriker = curBatsman.IsStriker
		} else if curBatsman.BatsmanID != int32(verifyBatsman.ID) && !curBatsman.IsStriker {
			nonStrikerResponse = curBatsman
		}
	}

	batsmanPlayerData, err := s.store.GetPlayerByID(ctx, int64(batsmanResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	nonStrikerPlayerData, err := s.store.GetPlayerByID(ctx, int64(nonStrikerResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	bowlerPlayerData, err := s.store.GetPlayerByID(ctx, int64(bowlerResponse.BowlerID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	batsman := map[string]interface{}{
		"player":               map[string]interface{}{"id": batsmanPlayerData.ID, "public_id": batsmanPlayerData.PublicID, "name": batsmanPlayerData.Name, "slug": batsmanPlayerData.Slug, "shortName": batsmanPlayerData.ShortName, "position": batsmanPlayerData.Positions},
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
		"inning_number":        batsmanResponse.InningNumber,
	}

	var emptyBatsman models.BatsmanScore
	var nonStriker map[string]interface{}

	if nonStrikerResponse != emptyBatsman {
		nonStriker = map[string]interface{}{
			"player":               map[string]interface{}{"id": nonStrikerPlayerData.ID, "public_id": nonStrikerPlayerData.PublicID, "name": nonStrikerPlayerData.Name, "slug": nonStrikerPlayerData.Slug, "shortName": nonStrikerPlayerData.ShortName, "position": nonStrikerPlayerData.Positions},
			"id":                   nonStrikerResponse.ID,
			"public_id":            nonStrikerPlayerData.PublicID,
			"match_id":             nonStrikerResponse.MatchID,
			"team_id":              nonStrikerResponse.TeamID,
			"batsman_id":           nonStrikerResponse.BatsmanID,
			"runs_scored":          nonStrikerResponse.RunsScored,
			"balls_faced":          nonStrikerResponse.BallsFaced,
			"fours":                nonStrikerResponse.Fours,
			"sixes":                nonStrikerResponse.Sixes,
			"batting_status":       nonStrikerResponse.BattingStatus,
			"is_striker":           nonStrikerResponse.IsStriker,
			"is_currently_batting": nonStrikerResponse.IsCurrentlyBatting,
			"inning_number":        nonStrikerResponse.InningNumber,
		}
	}

	bowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": bowlerPlayerData.ID, "public_id": bowlerPlayerData.PublicID, "name": bowlerPlayerData.Name, "slug": bowlerPlayerData.Slug, "shortName": bowlerPlayerData.ShortName, "position": bowlerPlayerData.Positions},
		"id":                bowlerResponse.ID,
		"public_id":         bowlerResponse.PublicID,
		"match_id":          bowlerResponse.MatchID,
		"team_id":           bowlerResponse.TeamID,
		"bowler_id":         bowlerResponse.BowlerID,
		"ball_number":       bowlerResponse.BallNumber,
		"runs":              bowlerResponse.Runs,
		"wide":              bowlerResponse.Wide,
		"no_ball":           bowlerResponse.NoBall,
		"wickets":           bowlerResponse.Wickets,
		"bowling_status":    bowlerResponse.BowlingStatus,
		"is_current_bowler": bowlerResponse.IsCurrentBowler,
		"inning_number":     bowlerResponse.InningNumber,
	}

	data = map[string]interface{}{
		"type": "UPDATE_SCORE",
		"payload": map[string]interface{}{
			"striker_batsman":     batsman,
			"non_striker_batsman": nonStriker,
			"bowler":              bowler,
			"inning_score":        inningScore,
			"event_type":          "no_ball",
		},
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
		return
	}
	return data
}

func (s *CricketServer) AddCricketWicketsWS(ctx *gin.Context, message map[string]interface{}) (data map[string]interface{}) {

	matchPublicID, err := uuid.Parse(message["match_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	battingTeamID, err := uuid.Parse(message["batting_team_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	batsmanPublicID, err := uuid.Parse(message["batsman_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	bowlerPublicID, err := uuid.Parse(message["bowler_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	fielderPublicID, err := uuid.Parse(message["fielder_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	runsScored := int(message["runs_scored"].(float64))
	inningNumber := int(message["inning_number"].(float64))
	wicketNumber := int(message["wicket_number"].(float64))
	ballNumber := int(message["ball_number"].(float64))
	wicketType := message["wicket_type"].(string)
	bowlType := message["bowl_type"].(*string)
	toggleStriker := message["toggle_striker"].(bool)

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction: ", err)
		return
	}

	defer tx.Rollback()

	cricketScore, err := s.store.GetCricketScoreByInning(ctx, matchPublicID, battingTeamID, inningNumber)
	if err != nil {
		s.logger.Error("Failed to get cricket score: ", err)
	}

	var outBatsmanResponse *models.BatsmanScore
	var notOutBatsmanResponse *models.BatsmanScore
	var bowlerResponse *models.BowlerScore
	var inningScoreResponse *models.CricketScore
	var wicketResponse *models.Wicket
	if bowlType != nil {
		outBatsmanResponse, notOutBatsmanResponse, bowlerResponse, inningScoreResponse, wicketResponse, err = s.store.AddCricketWicketWithBowlType(ctx, matchPublicID, battingTeamID, batsmanPublicID, bowlerPublicID, wicketNumber, wicketType, ballNumber, &fielderPublicID, cricketScore.Score, *bowlType, inningNumber)
		if err != nil {
			s.logger.Error("failed to add cricket wicket with bowl type: ", err)
		}
	} else {
		outBatsmanResponse, notOutBatsmanResponse, bowlerResponse, inningScoreResponse, wicketResponse, err = s.store.AddCricketWicket(ctx, matchPublicID, battingTeamID, batsmanPublicID, bowlerPublicID, int(cricketScore.Wickets), wicketType, int(cricketScore.Overs), fielderPublicID, cricketScore.Score, int32(runsScored), inningNumber)
		if err != nil {
			s.logger.Error("failed to add cricket wicket: ", err)
			return
		}
	}

	matchData, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("failed to get match: ", err)
		return
	}

	if inningScoreResponse.Wickets == 10 {
		inningScoreResponse, notOutBatsmanResponse, bowlerResponse, err = s.store.UpdateInningEndStatus(ctx, int32(matchData.ID), notOutBatsmanResponse.TeamID, inningNumber)
		if err != nil {
			s.logger.Error("failed to update inning_numberscore: ", err)
			return
		}
	} else if *matchData.MatchFormat == "T20" && inningScoreResponse.Overs/6 == 20 {
		inningScoreResponse, notOutBatsmanResponse, bowlerResponse, err = s.store.UpdateInningEndStatus(ctx, int32(matchData.ID), notOutBatsmanResponse.TeamID, inningNumber)
		if err != nil {
			s.logger.Error("failed to update inning_numberscore: ", err)
			return
		}
	} else if *matchData.MatchFormat == "ODI" && inningScoreResponse.Overs/6 == 50 {
		inningScoreResponse, notOutBatsmanResponse, bowlerResponse, err = s.store.UpdateInningEndStatus(ctx, int32(matchData.ID), notOutBatsmanResponse.TeamID, inningNumber)
		if err != nil {
			s.logger.Error("failed to update inning_numberscore: ", err)
			return
		}
	}

	err = s.UpdateMatchStatusAndResult(ctx, inningScoreResponse, matchData, matchData.ID)
	if err != nil {
		s.logger.Error("Failed to update match status and result: ", err)
		return
	}

	if toggleStriker {
		notOut, err := s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("failed to toggle batsman: ", err)
			return
		}
		notOutBatsmanResponse = &notOut[0]
	}

	var currentBatsman *models.BatsmanScore
	currentBatsman = notOutBatsmanResponse
	if bowlerResponse.BallNumber%6 == 0 {
		currentBatsmanResponse, err := s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
		currentBatsman = &currentBatsmanResponse[0]
	}

	outBatsmanPlayerData, err := s.store.GetPlayerByID(ctx, int64(outBatsmanResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
	}

	notOutBatsmanPlayerData, err := s.store.GetPlayerByID(ctx, int64(currentBatsman.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
	}

	bowlerPlayerData, err := s.store.GetPlayerByID(ctx, int64(bowlerResponse.BowlerID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
	}

	var fielderPlayerData *models.Player

	if wicketResponse.FielderID != nil {
		fielderPlayerData, err = s.store.GetPlayerByID(ctx, int64(*wicketResponse.FielderID))
		if err != nil {
			s.logger.Error("Failed to get Player: ", err)
		}
	}

	outBatsmanScore := map[string]interface{}{
		"player":               map[string]interface{}{"id": outBatsmanPlayerData.ID, "public_id": outBatsmanPlayerData.PublicID, "name": outBatsmanPlayerData.Name, "slug": outBatsmanPlayerData.Slug, "shortName": outBatsmanPlayerData.ShortName, "position": outBatsmanPlayerData.Positions},
		"id":                   outBatsmanResponse.ID,
		"public_id":            outBatsmanResponse.PublicID,
		"match_id":             outBatsmanResponse.MatchID,
		"team_id":              outBatsmanResponse.TeamID,
		"batsman_id":           outBatsmanResponse.BatsmanID,
		"runs_scored":          outBatsmanResponse.RunsScored,
		"balls_faced":          outBatsmanResponse.BallsFaced,
		"fours":                outBatsmanResponse.Fours,
		"sixes":                outBatsmanResponse.Sixes,
		"batting_status":       outBatsmanResponse.BattingStatus,
		"is_striker":           outBatsmanResponse.IsStriker,
		"is_currently_batting": outBatsmanResponse.IsCurrentlyBatting,
		"inning_number":        outBatsmanResponse.InningNumber,
	}

	notOutBatsmanScore := map[string]interface{}{
		"player":               map[string]interface{}{"id": notOutBatsmanPlayerData.ID, "public_id": notOutBatsmanPlayerData.PublicID, "name": notOutBatsmanPlayerData.Name, "slug": notOutBatsmanPlayerData.Slug, "shortName": notOutBatsmanPlayerData.ShortName, "position": notOutBatsmanPlayerData.Positions},
		"id":                   notOutBatsmanResponse.ID,
		"public_id":            notOutBatsmanResponse.PublicID,
		"match_id":             notOutBatsmanResponse.MatchID,
		"team_id":              notOutBatsmanResponse.TeamID,
		"batsman_id":           notOutBatsmanResponse.BatsmanID,
		"runs_scored":          notOutBatsmanResponse.RunsScored,
		"balls_faced":          notOutBatsmanResponse.BallsFaced,
		"fours":                notOutBatsmanResponse.Fours,
		"sixes":                notOutBatsmanResponse.Sixes,
		"batting_status":       notOutBatsmanResponse.BattingStatus,
		"is_striker":           notOutBatsmanResponse.IsStriker,
		"is_currently_batting": notOutBatsmanResponse.IsCurrentlyBatting,
		"inning_number":        notOutBatsmanResponse.InningNumber,
	}

	bowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": bowlerPlayerData.ID, "public_id": bowlerPlayerData.PublicID, "name": bowlerPlayerData.Name, "slug": bowlerPlayerData.Slug, "shortName": bowlerPlayerData.ShortName, "position": bowlerPlayerData.Positions},
		"id":                bowlerResponse.ID,
		"public_id":         bowlerResponse.PublicID,
		"match_id":          bowlerResponse.MatchID,
		"team_id":           bowlerResponse.TeamID,
		"bowler_id":         bowlerResponse.BowlerID,
		"ball_number":       bowlerResponse.BallNumber,
		"runs":              bowlerResponse.Runs,
		"wide":              bowlerResponse.Wide,
		"no_ball":           bowlerResponse.NoBall,
		"wickets":           bowlerResponse.Wickets,
		"bowling_status":    bowlerResponse.BowlingStatus,
		"is_current_bowler": bowlerResponse.IsCurrentBowler,
	}

	wickets := map[string]interface{}{
		"batsman_player": map[string]interface{}{"id": outBatsmanPlayerData.ID, "public_id": outBatsmanPlayerData.PublicID, "name": outBatsmanPlayerData.Name, "slug": outBatsmanPlayerData.Slug, "shortName": outBatsmanPlayerData.ShortName, "position": outBatsmanPlayerData.Positions},
		"bowler_player":  map[string]interface{}{"id": bowlerPlayerData.ID, "public_id": bowlerPlayerData.PublicID, "name": bowlerPlayerData.Name, "slug": bowlerPlayerData.Slug, "shortName": bowlerPlayerData.ShortName, "position": bowlerPlayerData.Positions},
		"fielder_player": fielderPlayerData,
		"id":             wicketResponse.ID,
		"public_id":      wicketResponse.PublicID,
		"match_id":       wicketResponse.MatchID,
		"team_id":        wicketResponse.TeamID,
		"batsman_id":     wicketResponse.BatsmanID,
		"bowler_id":      wicketResponse.BowlerID,
		"wicket_number":  wicketResponse.WicketsNumber,
		"wicket_type":    wicketResponse.WicketType,
		"ball_number":    wicketResponse.BallNumber,
		"fielder_id":     wicketResponse.FielderID,
		"score":          wicketResponse.Score,
		"inning_number":  wicketResponse.InningNumber,
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("failed to commit transaction: ", err)
		return
	}

	data = map[string]interface{}{
		"type": "UPDATE_SCORE",
		"payload": map[string]interface{}{
			"out_batsman":     outBatsmanScore,
			"not_out_batsman": notOutBatsmanScore,
			"bowler":          bowler,
			"inning_score":    inningScoreResponse,
			"wickets":         wickets,
			"match":           matchData,
			"event_type":      "wicket",
		},
	}
	return data
}

func (s *CricketServer) UpdateInningScoreWS(ctx *gin.Context, message map[string]interface{}) (data map[string]interface{}) {
	fmt.Println("Line no 1222: ", message)

	fmt.Println("Message Line no 1237: ", message)

	matchPublicID, err := uuid.Parse(message["match_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	batsmanTeamPublicID, err := uuid.Parse(message["batsman_team_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	batsmanPublicID, err := uuid.Parse(message["batsman_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	fmt.Println("Batsman: ", batsmanPublicID)

	bowlerPublicID, err := uuid.Parse(message["bowler_public_id"].(string))
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	runsScored := int(message["runs_scored"].(float64))
	inningNumber := int(message["inning_number"].(float64))
	fmt.Println("Bowler: ", bowlerPublicID)

	batsmanResponse, bowlerResponse, inningScore, err := s.store.UpdateInningScore(ctx, matchPublicID, batsmanTeamPublicID, batsmanPublicID, bowlerPublicID, int32(runsScored), inningNumber)
	if err != nil {
		s.logger.Error("Failed to update innings: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var currentBatsman []models.BatsmanScore
	var strikerResponse models.BatsmanScore
	var nonStrikerResponse models.BatsmanScore
	if bowlerResponse.BallNumber%6 == 0 && runsScored%2 == 0 {
		currentBatsman, err = s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
		fmt.Println("Even Run and Last Ball: ", currentBatsman)
	} else if bowlerResponse.BallNumber%6 != 0 && runsScored%2 != 0 {
		currentBatsman, err = s.store.ToggleCricketStricker(ctx, matchPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
		fmt.Println("odd Run and no Last Ball: ", currentBatsman)
	} else {
		currentBatsman, err = s.store.GetCurrentBattingBatsman(ctx, matchPublicID, batsmanTeamPublicID, inningNumber)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
		fmt.Println("other type: ", currentBatsman)
	}
	fmt.Println("Before toggle:", currentBatsman)

	for _, curBatsman := range currentBatsman {
		if curBatsman.IsStriker {
			strikerResponse = curBatsman
		} else {
			nonStrikerResponse = curBatsman
		}
	}

	fmt.Println("After toggle:", currentBatsman)

	matchData, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("failed to get match by match id: ", err)
		return
	}

	if inningScore.Wickets == 10 {
		inningScore, batsmanResponse, bowlerResponse, err = s.store.UpdateInningEndStatus(ctx, int32(matchData.ID), batsmanResponse.TeamID, inningNumber)
		if err != nil {
			s.logger.Error("failed to update inning_numberscore: ", err)
			return
		}
	} else if *matchData.MatchFormat == "T20" && inningScore.Overs/6 == 20 {
		inningScore, batsmanResponse, bowlerResponse, err = s.store.UpdateInningEndStatus(ctx, int32(matchData.ID), batsmanResponse.TeamID, inningNumber)
		if err != nil {
			s.logger.Error("failed to update inning_numberscore: ", err)
			return
		}
	} else if *matchData.MatchFormat == "ODI" && inningScore.Overs/6 == 50 {
		inningScore, batsmanResponse, bowlerResponse, err = s.store.UpdateInningEndStatus(ctx, int32(matchData.ID), batsmanResponse.TeamID, inningNumber)
		if err != nil {
			s.logger.Error("failed to update inning_numberscore: ", err)
			return
		}
	}

	err = s.UpdateMatchStatusAndResult(ctx, inningScore, matchData, matchData.ID)
	if err != nil {
		s.logger.Error("Failed to update match status and result: ", err)
		return
	}

	fmt.Println("STriker Response; ", strikerResponse)
	if strikerResponse.BatsmanID == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No striker found"})
		return
	}

	strikerPlayerData, err := s.store.GetPlayerByID(ctx, int64(strikerResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}
	fmt.Println("Striker Player Data; ", strikerPlayerData)

	nonStrikerPlayerData, err := s.store.GetPlayerByID(ctx, int64(nonStrikerResponse.BatsmanID))
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	bowlerPlayerData, err := s.store.GetPlayerByPublicID(ctx, bowlerPublicID)
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	striker := map[string]interface{}{
		"player":               map[string]interface{}{"id": strikerPlayerData.ID, "public_id": strikerPlayerData.PublicID, "name": strikerPlayerData.Name, "slug": strikerPlayerData.Slug, "shortName": strikerPlayerData.ShortName, "position": strikerPlayerData.Positions},
		"id":                   strikerResponse.ID,
		"public_id":            strikerResponse.PublicID,
		"match_id":             strikerResponse.MatchID,
		"team_id":              strikerResponse.TeamID,
		"batsman_id":           strikerResponse.BatsmanID,
		"runs_scored":          strikerResponse.RunsScored,
		"balls_faced":          strikerResponse.BallsFaced,
		"fours":                strikerResponse.Fours,
		"sixes":                strikerResponse.Sixes,
		"batting_status":       strikerResponse.BattingStatus,
		"is_striker":           strikerResponse.IsStriker,
		"is_currently_batting": strikerResponse.IsCurrentlyBatting,
		"inning_number":        strikerResponse.InningNumber,
	}

	nonStriker := map[string]interface{}{
		"player":               map[string]interface{}{"id": nonStrikerPlayerData.ID, "public_id": nonStrikerPlayerData.PublicID, "name": nonStrikerPlayerData.Name, "slug": nonStrikerPlayerData.Slug, "shortName": nonStrikerPlayerData.ShortName, "position": nonStrikerPlayerData.Positions},
		"id":                   nonStrikerResponse.ID,
		"public_id":            nonStrikerResponse.PublicID,
		"match_id":             nonStrikerResponse.MatchID,
		"team_id":              nonStrikerResponse.TeamID,
		"batsman_id":           nonStrikerResponse.BatsmanID,
		"runs_scored":          nonStrikerResponse.RunsScored,
		"balls_faced":          nonStrikerResponse.BallsFaced,
		"fours":                nonStrikerResponse.Fours,
		"sixes":                nonStrikerResponse.Sixes,
		"batting_status":       nonStrikerResponse.BattingStatus,
		"is_striker":           nonStrikerResponse.IsStriker,
		"is_currently_batting": nonStrikerResponse.IsCurrentlyBatting,
		"inning_number":        nonStrikerResponse.InningNumber,
	}

	bowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": bowlerPlayerData.ID, "public_id": bowlerPlayerData.PublicID, "name": bowlerPlayerData.Name, "slug": bowlerPlayerData.Slug, "shortName": bowlerPlayerData.ShortName, "position": bowlerPlayerData.Positions},
		"id":                bowlerResponse.ID,
		"public_id":         bowlerResponse.PublicID,
		"match_id":          bowlerResponse.MatchID,
		"team_id":           bowlerResponse.TeamID,
		"bowler_id":         bowlerResponse.BowlerID,
		"ball_number":       bowlerResponse.BallNumber,
		"runs":              bowlerResponse.Runs,
		"wide":              bowlerResponse.Wide,
		"no_ball":           bowlerResponse.NoBall,
		"wickets":           bowlerResponse.Wickets,
		"bowling_status":    bowlerResponse.BowlingStatus,
		"is_current_bowler": bowlerResponse.IsCurrentBowler,
		"inning_number":     bowlerResponse.InningNumber,
	}

	data = map[string]interface{}{
		"type": "UPDATE_SCORE",
		"payload": map[string]interface{}{
			"striker_batsman":     striker,
			"non_striker_batsman": nonStriker,
			"bowler":              bowler,
			"inning_score":        inningScore,
			"match":               matchData,
			"event_type":          "normal",
		},
	}

	return data
}

func (s *CricketServer) UpdateBowlingBowlerFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID         string `json:"match_public_id"`
		TeamPublicID          string `json:"team_public_id"`
		CurrentBowlerPublicID string `json:"current_bowler_public_id"`
		NextBowlerPublicID    string `json:"next_bowler_public_id"`
		InningNumber          int    `json:"inning_number"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
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

	currentBowlerPublicID, err := uuid.Parse(req.CurrentBowlerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	nextBowlerPublicID, err := uuid.Parse(req.NextBowlerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin Transcation")
	}

	defer tx.Rollback()

	currentBowlerResponse, err := s.store.UpdateBowlingBowlerStatus(ctx, matchPublicID, teamPublicID, currentBowlerPublicID, req.InningNumber)
	if err != nil {
		s.logger.Error("Failed to update current bowler status: ", err)
		return
	}

	nextBowlerResponse, err := s.store.UpdateBowlingBowlerStatus(ctx, matchPublicID, teamPublicID, nextBowlerPublicID, req.InningNumber)
	if err != nil {
		s.logger.Error("Failed to update next bowler status: ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
	}

	nextPlayerData, err := s.store.GetPlayerByPublicID(ctx, nextBowlerPublicID)
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	currentPlayerData, err := s.store.GetPlayerByPublicID(ctx, currentBowlerPublicID)
	if err != nil {
		s.logger.Error("Failed to get Player: ", err)
		return
	}

	nextBowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": nextPlayerData.ID, "public_id": nextPlayerData.PublicID, "name": nextPlayerData.Name, "slug": nextPlayerData.Slug, "shortName": nextPlayerData.ShortName, "position": nextPlayerData.Positions},
		"id":                nextBowlerResponse.ID,
		"public_id":         nextBowlerResponse.PublicID,
		"match_id":          nextBowlerResponse.MatchID,
		"team_id":           nextBowlerResponse.TeamID,
		"bowler_id":         nextBowlerResponse.BowlerID,
		"ball_number":       nextBowlerResponse.BallNumber,
		"runs":              nextBowlerResponse.Runs,
		"wide":              nextBowlerResponse.Wide,
		"no_ball":           nextBowlerResponse.NoBall,
		"wickets":           nextBowlerResponse.Wickets,
		"bowling_status":    nextBowlerResponse.BowlingStatus,
		"is_current_bowler": nextBowlerResponse.IsCurrentBowler,
		"inning_number":     nextBowlerResponse.InningNumber,
	}

	currentBowler := map[string]interface{}{
		"player":            map[string]interface{}{"id": currentPlayerData.ID, "public_id": currentPlayerData.PublicID, "name": currentPlayerData.Name, "slug": currentPlayerData.Slug, "shortName": currentPlayerData.ShortName, "position": currentPlayerData.Positions},
		"id":                currentBowlerResponse.ID,
		"public_id":         currentBowlerResponse.PublicID,
		"match_id":          currentBowlerResponse.MatchID,
		"team_id":           currentBowlerResponse.TeamID,
		"bowler_id":         currentBowlerResponse.BowlerID,
		"ball_number":       currentBowlerResponse.BallNumber,
		"runs":              currentBowlerResponse.Runs,
		"wide":              currentBowlerResponse.Wide,
		"no_ball":           currentBowlerResponse.NoBall,
		"wickets":           currentBowlerResponse.Wickets,
		"bowling_status":    currentBowlerResponse.BowlingStatus,
		"is_current_bowler": currentBowlerResponse.IsCurrentBowler,
		"inning_number":     currentBowlerResponse.InningNumber,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"next_bowler":    nextBowler,
		"current_bowler": currentBowler,
	})

}
