package tournaments

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TournamentMatchServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewTournamentMatchServer(store *db.Store, logger *logger.Logger) *TournamentMatchServer {
	return &TournamentMatchServer{store: store, logger: logger}
}

func (s *TournamentMatchServer) GetTournamentMatch(ctx *gin.Context) {

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusNotAcceptable, err)
		return
	}
	sports := ctx.Query("sports")
	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentID, sports))
	arg := db.GetTournamentMatchParams{
		TournamentID: tournamentID,
		Sports:       sports,
	}
	s.logger.Debug("Tournament match params: %v", arg)

	matches, err := s.store.GetTournamentMatch(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get tournament match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	tournament, err := s.store.GetTournament(ctx, tournamentID)
	if err != nil {
		s.logger.Error("Failed to get tournament: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	var matchDetails []map[string]interface{}

	for _, matchData := range matches {
		team1Name, err1 := s.store.GetClub(ctx, matchData.Team1ID)
		if err1 != nil {
			s.logger.Error("Failed to get club details for team1: %v", err1)
			continue
		}
		team2Name, err2 := s.store.GetClub(ctx, matchData.Team2ID)
		if err2 != nil {
			s.logger.Error("Failed to get club details for team2: %v", err2)
			continue
		}

		matchDetail := map[string]interface{}{
			"tournament_id":   matchData.TournamentID,
			"tournament_name": tournament.TournamentName,
			"match_id":        matchData.MatchID,
			"team1_id":        matchData.Team1ID,
			"team2_id":        matchData.Team2ID,
			"team1_name":      team1Name.ClubName,
			"team2_name":      team2Name.ClubName,
			"start_time":      matchData.StartTime,
			"end_time":        matchData.EndTime,
			"date_on":         matchData.DateOn,
			"sports":          matchData.Sports,
		}
		//s.logger.Debug("football match details: %v ", matchDetails)
		matchDetails = append(matchDetails, matchDetail)

	}
	checkSportServer := util.NewCheckSport(s.store, s.logger)
	tournamentMatches := checkSportServer.CheckSport(sports, matches, matchDetails)
	s.logger.Info("successfully  get the tournament match: %v", tournamentMatches)
	ctx.JSON(http.StatusAccepted, tournamentMatches)
	return
}

//add the instance from method

//add the sport matches:

type createTournamentMatchRequest struct {
	OrganizerID  int64     `json:"organizer_id"`
	TournamentID int64     `json:"tournament_id"`
	Team1ID      int64     `json:"team1_id"`
	Team2ID      int64     `json:"team2_id"`
	DateON       time.Time `json:"date_on"`
	StartTime    time.Time `json:"start_time"`
	Stage        string    `json:"stage"`
	Sports       string    `json:"sports"`
	EndTime      time.Time `json:"end_time"`
}

func (s *TournamentMatchServer) CreateTournamentMatch(ctx *gin.Context) {
	var req createTournamentMatchRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)
	arg := db.CreateMatchParams{
		OrganizerID:  req.OrganizerID,
		TournamentID: req.TournamentID,
		Team1ID:      req.Team1ID,
		Team2ID:      req.Team2ID,
		DateOn:       req.DateON,
		StartTime:    req.StartTime,
		Stage:        req.Stage,
		Sports:       req.Sports,
		EndTime:      req.EndTime,
	}

	s.logger.Debug("Create match params: %v", arg)

	response, err := s.store.CreateMatch(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create match: %v", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Debug("Successfully create match: %v", response)
	s.logger.Info("Successfully create match")

	ctx.JSON(http.StatusAccepted, response)
	return
}
