package cricket

import (
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addCricketBatScore struct {
	BatsmanID          int64  `json:"batsman_id"`
	MatchID            int64  `json:"match_id"`
	TeamID             int64  `json:"team_id"`
	Position           string `json:"position"`
	RunsScored         int32  `json:"runs_scored"`
	BallsFaced         int32  `json:"balls_faced"`
	Fours              int32  `json:"fours"`
	Sixes              int32  `json:"sixes"`
	BattingStatus      bool   `json:"batting_status"`
	IsStriker          bool   `json:"is_striker"`
	IsCurrentlyBatting bool   `json:"is_currently_batting"`
}

func (s *CricketServer) AddCricketBatScoreFunc(ctx *gin.Context) {
	var req addCricketBatScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind player batting score: ", err)
	}

	strickerResponse, err := s.store.GetCricketStricker(ctx, req.MatchID, req.TeamID)
	if err != nil {
		s.logger.Error("failed to get stricker: ", err)
		return
	}

	arg := db.AddCricketBatsScoreParams{
		BatsmanID:          req.BatsmanID,
		MatchID:            req.MatchID,
		TeamID:             req.TeamID,
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

	ctx.JSON(http.StatusAccepted, response)
	return
}

type addCricketBallScore struct {
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
	BowlerID        int64 `json:"bowler_id"`
	Ball            int32 `json:"ball"`
	Runs            int32 `json:"runs"`
	Wickets         int32 `json:"wickets"`
	Wide            int32 `json:"wide"`
	NoBall          int32 `json:"no_ball"`
	BowlingStatus   bool  `json:"bowling_status"`
	IsCurrentBowler bool  `json:"is_current_bowler"`
}

func (s *CricketServer) AddCricketBallFunc(ctx *gin.Context) {
	var req addCricketBallScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := db.AddCricketBallParams{
		MatchID:         req.MatchID,
		TeamID:          req.TeamID,
		BowlerID:        req.BowlerID,
		Ball:            req.Ball,
		Runs:            req.Runs,
		Wickets:         req.Wickets,
		Wide:            req.Wide,
		NoBall:          req.NoBall,
		BowlingStatus:   req.BowlingStatus,
		IsCurrentBowler: req.IsCurrentBowler,
	}

	response, err := s.store.AddCricketBall(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket bowler data: ", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

type addCricketWicketScore struct {
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	BatsmanID     int64  `json:"batsman_id"`
	BowlerID      int64  `json:"bowler_id"`
	WicketsNumber int32  `json:"wickets_number"`
	WicketType    string `json:"wicket_type"`
	BallNumber    int32  `json:"ball_number"`
	FielderID     *int32 `json:"fielder_id"`
}

func (s *CricketServer) AddCricketWicketFunc(ctx *gin.Context) {
	var req addCricketWicketScore
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind add cricket wickets: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := db.AddCricketWicketsParams{
		MatchID:       req.MatchID,
		TeamID:        req.TeamID,
		BatsmanID:     req.BatsmanID,
		BowlerID:      req.BowlerID,
		WicketsNumber: req.WicketsNumber,
		WicketType:    req.WicketType,
		BallNumber:    req.BallNumber,
		FielderID:     *req.FielderID,
	}

	response, err := s.store.AddCricketWickets(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket wicket: ", gin.H{"error": err.Error()})
		return
	}

	var updageCricketWickets *models.Wicket

	if updageCricketWickets != nil {
		arg := db.UpdateCricketWicketsParams{
			MatchID: req.MatchID,
			TeamID:  req.TeamID,
		}

		_, err := s.store.UpdateCricketWickets(ctx, arg)
		if err != nil {
			s.logger.Error("Failed to update the cricket wicket: ", gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketBatRequest struct {
	BatsmanID  int64  `json:"batsman_id"`
	TeamID     int64  `json:"team_id"`
	MatchID    int64  `json:"match_id"`
	Position   string `json:"position"`
	RunsScored int32  `json:"runs_scored"`
	BallsFaced int32  `json:"balls_faced"`
	Fours      int32  `json:"fours"`
	Sixes      int32  `json:"sixes"`
}

func (s *CricketServer) UpdateCricketBatScoreFunc(ctx *gin.Context) {
	var req updateCricketBatRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind udpate cricket bat score: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Debug("successfully bind :", req)

	arg := db.UpdateCricketRunsScoredParams{
		RunsScored: req.RunsScored,
		BallsFaced: req.BallsFaced,
		Fours:      req.Fours,
		Sixes:      req.Sixes,
		MatchID:    req.MatchID,
		BatsmanID:  req.BatsmanID,
		TeamID:     req.TeamID,
	}

	response, err := s.store.UpdateCricketRunsScored(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the cricket player runs: ", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateCricketBallRequest struct {
	BowlerID int64 `json:"bowler_id"`
	TeamID   int64 `json:"team_id"`
	MatchID  int64 `json:"match_id"`
	Ball     int32 `json:"ball"`
	Runs     int32 `json:"runs"`
	Wickets  int32 `json:"wickets"`
	Wide     int32 `json:"wide"`
	NoBall   int32 `json:"no_ball"`
}

func (s *CricketServer) UpdateCricketBallFunc(ctx *gin.Context) {
	var req updateCricketBallRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Debug("successfully bind: ", req)

	arg := db.UpdateCricketBowlerParams{
		Ball:     req.Ball,
		Runs:     req.Runs,
		Wickets:  req.Wickets,
		Wide:     req.Wide,
		NoBall:   req.NoBall,
		MatchID:  req.MatchID,
		BowlerID: req.BowlerID,
		TeamID:   req.TeamID,
	}

	fmt.Println("Arg: ", arg)
	response, err := s.store.UpdateCricketBowler(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the cricket bowler: ", gin.H{"error": err.Error()})
		return
	}

	argUpdateOver := db.UpdateCricketOversParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	_, err = s.store.UpdateCricketOvers(ctx, argUpdateOver)
	if err != nil {
		s.logger.Error("Failed to update the cricket overs: ", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

// type getPlayerScoreRequest struct {
// 	MatchID int64 `json:"match_id" form:"match_id"`
// 	TeamID  int64 `json:"team_id" form:"team_id"`
// }

func (s *CricketServer) GetPlayerScoreFunc(ctx *gin.Context) {
	// var req getPlayerScoreRequest
	// err := ctx.ShouldBindQuery(&req)
	// if err != nil {
	// 	s.logger.Error("Failed to bind player score: ", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	matchIDString := ctx.Query("match_id")
	teamIDString := ctx.Query("team_id")
	matchID, err := strconv.ParseInt(matchIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match id ", err)
		return
	}

	teamID, err := strconv.ParseInt(teamIDString, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse team id ", err)
		return
	}

	arg := db.GetCricketPlayersScoreParams{
		TeamID:  teamID,
		MatchID: matchID,
	}

	fmt.Println("Player score; ", arg)

	teamPlayerScore, err := s.store.GetCricketPlayersScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get players score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Team Player Score: ", teamPlayerScore)
	match, err := s.store.GetMatchByMatchID(ctx, matchID)
	if err != nil {
		s.logger.Error("Failed to get match:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var battingTeamId int64
	var bowlingTeamId int64
	if teamID == match.HomeTeamID {
		battingTeamId = teamID
		bowlingTeamId = match.AwayTeamID
	} else {
		battingTeamId = teamID
		bowlingTeamId = match.HomeTeamID
	}

	battingTeam, err := s.store.GetTeam(ctx, battingTeamId)
	if err != nil {
		s.logger.Error("Failed to get batting team : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var battingDetails []map[string]interface{}
	for _, playerScore := range teamPlayerScore {
		playerData, err := s.store.GetPlayer(ctx, playerScore.BatsmanID)
		if err != nil {
			s.logger.Error("Failed to get players data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		battingDetails = append(battingDetails, map[string]interface{}{
			"player":               map[string]interface{}{"id": playerData.ID, "name": playerData.PlayerName, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions, "username": playerData.Username},
			"runsScored":           playerScore.RunsScored,
			"ballFaced":            playerScore.BallsFaced,
			"fours":                playerScore.Fours,
			"sixes":                playerScore.Sixes,
			"batting_status":       playerScore.BattingStatus,
			"is_striker":           playerScore.IsStriker,
			"is_currently_batting": playerScore.IsCurrentlyBatting,
		})
	}
	var scoreDetails map[string]interface{}
	var emptyDetails map[string]interface{}
	if len(battingDetails) >= 1 {
		scoreDetails = map[string]interface{}{
			"battingTeam": map[string]interface{}{"id": battingTeam.ID, "name": battingTeam.Name, "slug": battingTeam.Slug, "shortName": battingTeam.Shortname, "gender": battingTeam.Gender, "national": battingTeam.National, "country": battingTeam.Country, "type": battingTeam.Type},
			"innings":     battingDetails,
		}
	} else {
		scoreDetails = map[string]interface{}{
			"battingTeam": map[string]interface{}{"id": battingTeam.ID, "name": battingTeam.Name, "slug": battingTeam.Slug, "shortName": battingTeam.Shortname, "gender": battingTeam.Gender, "national": battingTeam.National, "country": battingTeam.Country, "type": battingTeam.Type},
			"innings":     emptyDetails,
		}
	}

	argCricketScore := db.UpdateCricketScoreParams{
		MatchID:       matchID,
		BattingTeamID: battingTeamId,
		BowlingTeamID: bowlingTeamId,
	}
	fmt.Println("arg: cricket score: ", argCricketScore)
	_, err = s.store.UpdateCricketScore(ctx, argCricketScore)
	if err != nil {
		s.logger.Error("Failed to update the cricket score: ", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, scoreDetails)
}

type getCricketBowlersRequest struct {
	MatchID int64 `json:"match_id" form:"match_id"`
	TeamID  int64 `json:"team_id" form:"team_id"`
}

func (s *CricketServer) GetCricketBowlerFunc(ctx *gin.Context) {
	var req getCricketBowlersRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketBallsParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}
	playerScore, err := s.store.GetCricketBalls(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get cricket bowler data : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := s.store.GetMatchByMatchID(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to get player :", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var battingTeamId int64
	var bowlingTeamId int64
	if req.TeamID == match.HomeTeamID {
		battingTeamId = match.AwayTeamID
		bowlingTeamId = req.TeamID
	} else {
		battingTeamId = match.HomeTeamID
		bowlingTeamId = req.TeamID
	}

	bowlingTeam, err := s.store.GetTeam(ctx, bowlingTeamId)
	if err != nil {
		s.logger.Error("Failed to get players score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bowlingDetails = make([]map[string]interface{}, len(playerScore))
	for i, playerScore := range playerScore {
		playerData, err := s.store.GetPlayer(ctx, playerScore.BowlerID)
		if err != nil {
			s.logger.Error("Failed to get players data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		bowlingDetails[i] = map[string]interface{}{
			"player":            map[string]interface{}{"id": playerData.ID, "name": playerData.PlayerName, "slug": playerData.Slug, "shortName": playerData.ShortName, "position": playerData.Positions, "username": playerData.Username},
			"runs":              playerScore.Runs,
			"ball":              playerScore.Ball,
			"wide":              playerScore.Wide,
			"noBall":            playerScore.NoBall,
			"wickets":           playerScore.Wickets,
			"bowling_status":    playerScore.BowlingStatus,
			"is_current_bowler": playerScore.IsCurrentBowler,
		}
	}

	scoreDetails := map[string]interface{}{
		"bowlingTeam": map[string]interface{}{
			"id":        bowlingTeam.ID,
			"name":      bowlingTeam.Name,
			"slug":      bowlingTeam.Slug,
			"shortName": bowlingTeam.Shortname,
			"gender":    bowlingTeam.Gender,
			"national":  bowlingTeam.National,
			"country":   bowlingTeam.Country,
			"type":      bowlingTeam.Type,
		},
		"innings": bowlingDetails,
	}

	arg1 := db.UpdateCricketOversParams{
		MatchID: req.MatchID,
		TeamID:  battingTeamId,
	}

	_, err = s.store.UpdateCricketOvers(ctx, arg1)
	if err != nil {
		s.logger.Error("Failed to update the cricket overs: ", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, scoreDetails)
}

type getCricketWicketsRequest struct {
	MatchID int64 `json:"match_id" form:"match_id"`
	TeamID  int64 `json:"team_id" form:"team_id"`
}

func (s *CricketServer) GetCricketWicketsFunc(ctx *gin.Context) {
	var req getCricketWicketsRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		s.logger.Error("Failed to bind  get cricket wickets : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketWicketsParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}
	s.logger.Debug("cricket wicket arg: ", arg)
	response, err := s.store.GetCricketWickets(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get cricket bowler score : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Successfully get the wickets: ", response)

	argCricketTeamWicket := db.UpdateCricketWicketsParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	_, err = s.store.UpdateCricketWickets(ctx, argCricketTeamWicket)
	if err != nil {
		s.logger.Error("Failed to upate cricket wicket : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wicketsData []map[string]interface{}

	// teamData, err := s.store.GetTeam()
	argMatchScore := db.GetCricketScoreParams{
		MatchID: req.MatchID,
		TeamID:  req.TeamID,
	}

	_, err = s.store.GetCricketScore(ctx, argMatchScore)
	if err != nil {
		s.logger.Error("Failed to get current score data : ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	for _, wicket := range response {
		batsmanData, err := s.store.GetPlayer(ctx, wicket.BatsmanID)
		if err != nil {
			s.logger.Error("Failed to get batsman data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		bowlerData, err := s.store.GetPlayer(ctx, wicket.BowlerID)
		if err != nil {
			s.logger.Error("Failed to get bowler data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		fielderData, err := s.store.GetPlayer(ctx, *wicket.FielderID)
		if err != nil {
			s.logger.Error("Failed to get fielder data : ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var emptyPlayer models.Player
		var fielder map[string]interface{}
		if fielderData == emptyPlayer {
			fielder = map[string]interface{}{"id": fielderData.ID, "name": fielderData.PlayerName, "slug": fielderData.Slug, "shortName": fielderData.ShortName, "position": fielderData.Positions, "username": fielderData.Username}
		}

		wicketData := map[string]interface{}{
			"batsman":      map[string]interface{}{"id": batsmanData.ID, "name": batsmanData.PlayerName, "slug": batsmanData.Slug, "shortName": batsmanData.ShortName, "position": batsmanData.Positions, "username": batsmanData.Username},
			"bowler":       map[string]interface{}{"id": bowlerData.ID, "name": bowlerData.PlayerName, "slug": bowlerData.Slug, "shortName": bowlerData.ShortName, "position": bowlerData.Positions, "username": bowlerData.Username},
			"wicketNumber": wicket.WicketsNumber,
			"wicketType":   wicket.WicketType,
			"Overs":        wicket.BallNumber,
			"fielder":      fielder,
		}
		wicketsData = append(wicketsData, wicketData)
	}
	fmt.Println("Wicket : ", wicketsData)
	s.logger.Debug("Successfully update the wickets: ", wicketsData)

	ctx.JSON(http.StatusAccepted, wicketsData)
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

func (s *CricketServer) UpdateRunsScoreFunc(ctx *gin.Context) {
	var req updateCricketPlayerStatsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind udpate cricket bat score: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	_, err = s.store.UpdateCricketBatsmanScore(ctx, req.RunsScored, req.MatchID, req.BatsmanID)
	if err != nil {
		s.logger.Error("Failed to update batsman scored: ", err)
		return
	}

	bowlerStats, err := s.store.UpdateBowlerStats(ctx, req.RunsScored, req.MatchID, req.BowlerID)
	if err != nil {
		s.logger.Error("Failed to update bowler scored: ", err)
		return
	}

	if bowlerStats.Ball%6 == 0 && req.RunsScored%2 == 0 {
		err = s.store.ToggleCricketStricker(ctx, req.MatchID)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
	} else if bowlerStats.Ball%6 != 0 && req.RunsScored%2 != 0 {
		err = s.store.ToggleCricketStricker(ctx, req.MatchID)
		if err != nil {
			s.logger.Error("Failed to update stricker: ", err)
		}
	}
}

type updateWideRunsRequest struct {
	MatchID       int64 `json:"match_id"`
	BowlerID      int64 `json:"bowler_id"`
	BattingTeamID int64 `json:"batting_team_id"`
}

func (s *CricketServer) UpdateWideBallFunc(ctx *gin.Context) {
	var req updateWideRunsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	err = s.store.UpdateWideRuns(ctx, req.MatchID, req.BowlerID, req.BattingTeamID)
	if err != nil {
		s.logger.Error("Failed to update wide: ", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}
}

type updateNoBallRuns struct {
	RunsScored    int32 `json:"runs_scored"`
	MatchID       int64 `json:"match_id"`
	BowlerID      int64 `json:"bowler_id"`
	BattingTeamID int32 `json:"team_id"`
}

func (s *CricketServer) UpdateNoBallsRunsFunc(ctx *gin.Context) {
	var req updateNoBallRuns
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	err = s.store.UpdateNoBallsRuns(ctx, req.RunsScored, req.MatchID, req.BowlerID, int64(req.BattingTeamID))
	if err != nil {
		s.logger.Error("Failed to update wide: ", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}
}

type addCricketWicketReq struct {
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
	BatsmanID    int64  `json:"batsman_id"`
	BowlerID     int64  `json:"bowler_id"`
	WicketNumber int    `json:"wicket_number"`
	WicketType   string `json:"wicket_type"`
	BallNumber   int    `json:"ball_number"`
	FielderID    int64  `json:"fielder_id"`
	RunsScored   int32  `json:"runs_scored"`
}

func (s *CricketServer) AddCricketWicketsFunc(ctx *gin.Context) {
	var req addCricketWicketReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("failed to bind: ", err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	err = s.store.AddCricketWicket(ctx, req.MatchID, req.TeamID, req.BatsmanID, req.BowlerID, req.WicketNumber, req.WicketType, req.BallNumber, req.FielderID, req.RunsScored)
	if err != nil {
		s.logger.Error("failed to add cricket wicket: ", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("failed to commit transcation: ", err)
		return
	}
}
