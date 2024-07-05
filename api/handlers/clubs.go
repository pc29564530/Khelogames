package handlers

import (
	"encoding/base64"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClubServer struct {
	store  *db.Store
	logger *logger.Logger
}
type createClubRequest struct {
	ClubName  string `json:"club_name"`
	AvatarURL string `json:"avatar_url"`
	Sport     string `json:"sport"`
}

func NewClubServer(store *db.Store, logger *logger.Logger) *ClubServer {
	return &ClubServer{store: store, logger: logger}
}

func (s *ClubServer) CreateClubFunc(ctx *gin.Context) {
	var req createClubRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind create club request: %v", err)
		return
	}

	var path string
	if req.AvatarURL != "" {
		b64Data := req.AvatarURL[strings.IndexByte(req.AvatarURL, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			s.logger.Error("Failed to decode string: %v", err)
			return
		}

		path, err = util.SaveImageToFile(data, "image")
		if err != nil {
			s.logger.Error("Failed to create file: %v", err)
			return
		}
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateClubParams{
		ClubName:  req.ClubName,
		AvatarUrl: path,
		Sport:     req.Sport,
		Owner:     authPayload.Username,
	}

	response, err := s.store.CreateClub(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetClubsFunc(ctx *gin.Context) {

	response, err := s.store.GetClubs(ctx)
	if err != nil {
		s.logger.Error("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type getClubRequest struct {
	ID int64 `uri:"id"`
}

func (s *ClubServer) GetClubFunc(ctx *gin.Context) {
	var req getClubRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	response, err := s.store.GetClub(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type updateAvatarUrlRequest struct {
	AvatarUrl string `json:"avatar_url"`
	ClubName  string `json:"club_name"`
}

func (s *ClubServer) UpdateClubAvatarUrlFunc(ctx *gin.Context) {
	var req updateAvatarUrlRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateAvatarUrlParams{
		AvatarUrl: req.AvatarUrl,
		ClubName:  req.ClubName,
	}

	response, err := s.store.UpdateAvatarUrl(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update avatar url: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

// type updateClubNameRequest struct {
// 	AvatarUrl string `json:"avatar_url"`
// 	ClubName  string `json:"club_name"`
// }

// func (server *Server) updateClubName(ctx *gin.Context) {
// 	var req updateClubNameRequest
// 	err := ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, (err))
// 		return
// 	}
// 	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
// 	arg := db.UpdateClubNameParams{
// 		ClubName: req.ClubName,
// 	}

// 	response, err :=s.store.UpdateClubName(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, (err))
// 		return
// 	}

// 	ctx.JSON(http.StatusAccepted, response)
// 	return
// }

type updateClubSport struct {
	ClubName string `json:"club_name"`
	Sport    string `json:"sport"`
}

func (s *ClubServer) UpdateClubSportFunc(ctx *gin.Context) {
	var req updateClubSport
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.UpdateClubSportParams{
		Sport:    req.Sport,
		ClubName: req.ClubName,
	}

	response, err := s.store.UpdateClubSport(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update club sport: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

type searchTeamRequest struct {
	ClubName string `json:"club_name"`
}

func (s *ClubServer) SearchTeamFunc(ctx *gin.Context) {
	var req searchTeamRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	searchQuery := "%" + req.ClubName + "%"

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	response, err := s.store.SearchTeam(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetClubsBySportFunc(ctx *gin.Context) {

	sports := ctx.Param("sport")

	response, err := s.store.GetClubsBySport(ctx, sports)
	if err != nil {
		s.logger.Error("Failed to get club by sport: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetTournamentsByClubFunc(ctx *gin.Context) {
	clubName := ctx.Query("club_name")
	response, err := s.store.GetTournamentsByClub(ctx, clubName)
	if err != nil {
		s.logger.Error("Failed to get tournament by club: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *ClubServer) GetMatchByClubNameFunc(ctx *gin.Context) {
	clubIdStr := ctx.Query("id")
	clubID, err := strconv.ParseInt(clubIdStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse club id: %v", err)
		return
	}
	matches, err := s.store.GetMatchByClubName(ctx, clubID)
	if err != nil {
		s.logger.Error("Failed to get match by clubname: %v", err)
		ctx.JSON(http.StatusNoContent, (err))
		return
	}

	var clubMatchDetails []map[string]interface{}

	for _, match := range matches {
		matchScoreData := s.getMatchScore(ctx, match)
		s.logger.Debug("Match Score Data: ", matchScoreData)
		s.logger.Debug("matches: ", match)
		clubMatchDetails = append(clubMatchDetails, matchScoreData)
	}

	s.logger.Info("Club match details: ", clubMatchDetails)

	ctx.JSON(http.StatusAccepted, clubMatchDetails)
	return
}

func (s *ClubServer) getMatchScore(ctx *gin.Context, match db.GetMatchByClubNameRow) map[string]interface{} {
	clubMatchDetail := map[string]interface{}{
		"tournament_id":   match.TournamentID,
		"tournament_name": match.TournamentName,
		"match_id":        match.MatchID,
		"team1_id":        match.Team1ID,
		"team2_id":        match.Team2ID,
		"team1_name":      match.Team1Name,
		"team2_name":      match.Team2Name,
		"date_on":         match.DateOn,
		"start_time":      match.StartTime,
		"end_time":        match.EndTime,
	}

	switch match.Sports {
	case "Cricket":
		return s.getCricketMatchScore(ctx, match, clubMatchDetail)
	case "Football":
		return s.getFootballMatchScore(ctx, match, clubMatchDetail)
	default:
		s.logger.Error("Unsupported sport type:", match.Sports)
		return nil
	}
}

func (s *ClubServer) getCricketMatchScore(ctx *gin.Context, match db.GetMatchByClubNameRow, clubMatchDetail map[string]interface{}) map[string]interface{} {
	arg1 := db.GetCricketMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team1ID}
	arg2 := db.GetCricketMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team2ID}

	matchScoreData1, err := s.store.GetCricketMatchScore(ctx, arg1)
	if err != nil {
		s.logger.Error("Failed to get cricket match score for team 1:", err)
		return nil
	}
	matchScoreData2, err := s.store.GetCricketMatchScore(ctx, arg2)
	if err != nil {
		s.logger.Error("Failed to get cricket match score for team 2:", err)
		return nil
	}

	clubMatchDetail["team1_score"] = matchScoreData1.Score
	clubMatchDetail["team1_wickets"] = matchScoreData1.Wickets
	clubMatchDetail["team1_extras"] = matchScoreData1.Extras
	clubMatchDetail["team1_overs"] = matchScoreData1.Overs
	clubMatchDetail["team1_innings"] = matchScoreData1.Innings
	clubMatchDetail["team2_score"] = matchScoreData2.Score
	clubMatchDetail["team2_wickets"] = matchScoreData2.Wickets
	clubMatchDetail["team2_extras"] = matchScoreData2.Extras
	clubMatchDetail["team2_overs"] = matchScoreData2.Overs
	clubMatchDetail["team2_innings"] = matchScoreData2.Innings

	return clubMatchDetail
}

func (s *ClubServer) getFootballMatchScore(ctx *gin.Context, match db.GetMatchByClubNameRow, clubMatchDetail map[string]interface{}) map[string]interface{} {
	arg1 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team1ID}
	arg2 := db.GetFootballMatchScoreParams{MatchID: match.MatchID, TeamID: match.Team2ID}

	matchScoreData1, err := s.store.GetFootballMatchScore(ctx, arg1)
	if err != nil {
		s.logger.Error("Failed to get football match score for team 1:", err)
		return nil
	}
	matchScoreData2, err := s.store.GetFootballMatchScore(ctx, arg2)
	if err != nil {
		s.logger.Error("Failed to get football match score for team 2:", err)
		return nil
	}

	clubMatchDetail["team1_score"] = matchScoreData1.GoalFor
	clubMatchDetail["team2_score"] = matchScoreData2.GoalFor

	return clubMatchDetail
}
