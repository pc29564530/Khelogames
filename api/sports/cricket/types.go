package cricket

// // Request Types
// type AddCricketBatScoreRequest struct {
// 	MatchPublicID      string `json:"match_public_id" binding:"required"`
// 	TeamPublicID       string `json:"team_public_id" binding:"required"`
// 	BatsmanPublicID    string `json:"batsman_public_id" binding:"required"`
// 	Position           string `json:"position" binding:"required"`
// 	RunsScored         int32  `json:"runs_scored" binding:"min=0"`
// 	BallsFaced         int32  `json:"balls_faced" binding:"min=0"`
// 	Fours              int32  `json:"fours" binding:"min=0"`
// 	Sixes              int32  `json:"sixes" binding:"min=0"`
// 	BattingStatus      bool   `json:"batting_status"`
// 	IsStriker          bool   `json:"is_striker"`
// 	IsCurrentlyBatting bool   `json:"is_currently_batting"`
// 	InningNumber       int    `json:"inning_number" binding:"min=1"`
// }

// type AddCricketBallScoreRequest struct {
// 	MatchPublicID      string `json:"match_public_id" binding:"required"`
// 	TeamPublicID       string `json:"team_public_id" binding:"required"`
// 	BowlerPublicID     string `json:"bowler_public_id" binding:"required"`
// 	PrevBowlerPublicID string `json:"prev_bowler_public_id"`
// 	BallNumber         int32  `json:"ball_number" binding:"min=1,max=6"`
// 	Runs               int32  `json:"runs" binding:"min=0"`
// 	Wickets            int32  `json:"wickets" binding:"min=0,max=1"`
// 	Wide               int32  `json:"wide" binding:"min=0"`
// 	NoBall             int32  `json:"no_ball" binding:"min=0"`
// 	InningNumber       int    `json:"inning_number" binding:"min=1"`
// }

// type UpdateInningScoreRequest struct {
// 	MatchPublicID       string `json:"match_public_id" binding:"required"`
// 	BatsmanTeamPublicID string `json:"batsman_team_public_id" binding:"required"`
// 	BatsmanPublicID     string `json:"batsman_public_id" binding:"required"`
// 	BowlerPublicID      string `json:"bowler_public_id" binding:"required"`
// 	RunsScored          int    `json:"runs_scored" binding:"min=0,max=6"`
// 	InningNumber        int    `json:"inning_number" binding:"min=1"`
// }

// // Response Types
// type PlayerInfo struct {
// 	ID        int64     `json:"id"`
// 	PublicID  uuid.UUID `json:"public_id"`
// 	Name      string    `json:"name"`
// 	Slug      string    `json:"slug"`
// 	ShortName string    `json:"shortName"`
// 	Position  string    `json:"position"`
// }

// type BatsmanScore struct {
// 	Player             PlayerInfo `json:"player"`
// 	ID                 int64      `json:"id"`
// 	PublicID           uuid.UUID  `json:"public_id"`
// 	MatchID            int32      `json:"match_id"`
// 	TeamID             int32      `json:"team_id"`
// 	BatsmanID          int32      `json:"batsman_id"`
// 	RunsScored         int        `json:"runs_scored"`
// 	BallsFaced         int        `json:"balls_faced"`
// 	Fours              int        `json:"fours"`
// 	Sixes              int        `json:"sixes"`
// 	BattingStatus      bool       `json:"batting_status"`
// 	IsStriker          bool       `json:"is_striker"`
// 	IsCurrentlyBatting bool       `json:"is_currently_batting"`
// 	InningNumber       int        `json:"inning_number"`
// }

// type BowlerScore struct {
// 	Player          PlayerInfo `json:"player"`
// 	ID              int64      `json:"id"`
// 	PublicID        uuid.UUID  `json:"public_id"`
// 	MatchID         int32      `json:"match_id"`
// 	TeamID          int32      `json:"team_id"`
// 	BowlerID        int32      `json:"bowler_id"`
// 	BallNumber      int        `json:"ball_number"`
// 	Runs            int        `json:"runs"`
// 	Wide            int        `json:"wide"`
// 	NoBall          int        `json:"no_ball"`
// 	Wickets         int        `json:"wickets"`
// 	BowlingStatus   bool       `json:"bowling_status"`
// 	IsCurrentBowler bool       `json:"is_current_bowler"`
// 	InningNumber    int        `json:"inning_number"`
// }

// type InningScore struct {
// 	ID                int64     `json:"id"`
// 	PublicID          uuid.UUID `json:"public_id"`
// 	MatchID           int32     `json:"match_id"`
// 	TeamID            int32     `json:"team_id"`
// 	InningNumber      int       `json:"inning_number"`
// 	Score             int       `json:"score"`
// 	Wickets           int       `json:"wickets"`
// 	Overs             int       `json:"overs"`
// 	RunRate           string    `json:"run_rate"`
// 	TargetRunRate     string    `json:"target_run_rate"`
// 	FollowOn          bool      `json:"follow_on"`
// 	IsInningCompleted bool      `json:"is_inning_completed"`
// 	Declared          bool      `json:"declared"`
// }

// type CricketScoreUpdateResponse struct {
// 	Type    string      `json:"type"`
// 	Payload interface{} `json:"payload"`
// }

// type ErrorResponse struct {
// 	Error   string `json:"error"`
// 	Code    string `json:"code,omitempty"`
// 	Details string `json:"details,omitempty"`
// }

// // Constants
// const (
// 	EventTypeNormal = "normal"
// 	EventTypeWide   = "wide"
// 	EventTypeNoBall = "no_ball"
// 	EventTypeWicket = "wicket"
// 	EventTypeInning = "inning"

// 	ScoreUpdateType  = "UPDATE_SCORE"
// 	InningStatusType = "INNING_STATUS"
// 	AddBatsmanType   = "ADD_BATSMAN"
// 	AddBowlerType    = "ADD_BOWLER"

// 	InningStatusInProgress = "in_progress"
// 	InningStatusCompleted  = "completed"
// )

// // Validation helpers
// func (r *AddCricketBatScoreRequest) Validate() error {
// 	if r.RunsScored < 0 {
// 		return fmt.Errorf("runs scored cannot be negative")
// 	}
// 	if r.BallsFaced < 0 {
// 		return fmt.Errorf("balls faced cannot be negative")
// 	}
// 	if r.InningNumber < 1 {
// 		return fmt.Errorf("inning number must be at least 1")
// 	}
// 	return nil
// }

// func (r *AddCricketBallScoreRequest) Validate() error {
// 	if r.BallNumber < 1 || r.BallNumber > 6 {
// 		return fmt.Errorf("ball number must be between 1 and 6")
// 	}
// 	if r.Runs < 0 {
// 		return fmt.Errorf("runs cannot be negative")
// 	}
// 	if r.InningNumber < 1 {
// 		return fmt.Errorf("inning number must be at least 1")
// 	}
// 	return nil
// }
