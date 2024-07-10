package cricket

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CricketPlayerScoreServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCricketPlayerServer(store *db.Store, logger *logger.Logger) *CricketPlayerScoreServer {
	return &CricketPlayerScoreServer{store: store, logger: logger}
}

type addCricketTeamPlayerScoreRequest struct {
	MatchID          int64  `json:"match_id"`
	TournamentID     int64  `json:"tournament_id"`
	TeamID           int64  `json:"team_id"`
	BattingOrBowling string `json:"batting_or_bowling, omit_empty"`
	Position         string `json:"position"`
	PlayerID         int64  `json:"player_id"`
	RunsScored       string `json:"runs_scored"`
	BallsFaced       string `json:"balls_faced"`
	Fours            string `json:"fours"`
	Sixes            string `json:"sixes"`
	WicketsTaken     string `json:"wickets_taken"`
	OversBowled      string `json:"overs_bowled"`
	RunsConceded     string `json:"runs_conceded"`
	WicketTakenBy    string `json:"wicket_taken_by"`
	WicketOf         string `json:"wicket_of"`
}

func (s *CricketPlayerScoreServer) AddCricketPlayerScoreFunc(ctx *gin.Context) {
	var req addCricketTeamPlayerScoreRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var position int64
	if req.Position == "" {
		position = 0
	} else {
		position, err = strconv.ParseInt(req.Position, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse position: %v", err)
			return
		}
	}

	var runsScored int64
	if req.RunsScored == "" {
		runsScored = 0
	} else {
		runsScored, err = strconv.ParseInt(req.RunsScored, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse runs scored: %v", err)
			return
		}
	}

	var ballsFaced int64
	if req.BallsFaced == "" {
		ballsFaced = 0
	} else {
		ballsFaced, err = strconv.ParseInt(req.BallsFaced, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse ball faced: %v", err)
			return
		}
	}

	var fours int64
	if req.Fours == "" {
		fours = 0
	} else {
		fours, err = strconv.ParseInt(req.Fours, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse fours: %v", err)
			return
		}
	}

	var sixes int64
	if req.Sixes == "" {
		runsScored = 0
	} else {
		sixes, err = strconv.ParseInt(req.Sixes, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse sixes: %v", err)
			return
		}
	}

	var wicketsTaken int64
	if req.WicketsTaken == "" {
		wicketsTaken = 0
	} else {
		wicketsTaken, err = strconv.ParseInt(req.WicketsTaken, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse wicket taken: %v", err)
			return
		}
	}

	var runsConceded int64
	if req.RunsConceded == "" {
		runsConceded = 0
	} else {
		runsConceded, err = strconv.ParseInt(req.RunsConceded, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse runs conceded: %v", err)
			return
		}
	}

	var wicketTakenBy int64
	if req.WicketTakenBy == "" {
		wicketTakenBy = 0
	} else {
		wicketTakenBy, err = strconv.ParseInt(req.WicketTakenBy, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse wicket taken by: %v", err)
			return
		}
	}

	var wicketOf int64
	if req.WicketOf == "" {
		runsScored = 0
	} else {
		wicketOf, err = strconv.ParseInt(req.WicketOf, 10, 64)
		if err != nil {
			s.logger.Error("Failed to parse wicket of: %v", err)
			return
		}
	}
	arg := db.AddCricketTeamPlayerScoreParams{
		MatchID:          req.MatchID,
		TournamentID:     req.TournamentID,
		TeamID:           req.TeamID,
		BattingOrBowling: req.BattingOrBowling,
		Position:         position,
		PlayerID:         req.PlayerID,
		RunsScored:       runsScored,
		BallsFaced:       ballsFaced,
		Fours:            fours,
		Sixes:            sixes,
		WicketsTaken:     wicketsTaken,
		OversBowled:      req.OversBowled,
		RunsConceded:     runsConceded,
		WicketTakenBy:    wicketTakenBy,
		WicketOf:         wicketOf,
	}

	response, err := s.store.AddCricketTeamPlayerScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to add the cricket player score: %v", gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}

func (s *CricketPlayerScoreServer) GetCricketPlayerScoreFunc(ctx *gin.Context) {

	matchIDStr := ctx.Query("match_id")
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse match id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentIDStr := ctx.Query("tournament_id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to parse tournament id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamIdStr := ctx.Query("team_id")
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		s.logger.Error("Failed to team id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.GetCricketTeamPlayerScoreParams{
		MatchID:      matchID,
		TournamentID: tournamentID,
		TeamID:       teamID,
	}

	response, err := s.store.GetCricketTeamPlayerScore(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get cricket team player score : %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
	return
}
