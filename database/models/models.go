package models

import (
	"time"

	"github.com/google/uuid"
)

type Ball struct {
	ID       int64 `json:"id"`
	TeamID   int64 `json:"team_id"`
	MatchID  int64 `json:"match_id"`
	BowlerID int64 `json:"bowler_id"`
	Ball     int32 `json:"ball"`
	Runs     int32 `json:"runs"`
	Wickets  int32 `json:"wickets"`
	Wide     int32 `json:"wide"`
	NoBall   int32 `json:"no_ball"`
}

type Bat struct {
	ID         int64 `json:"id"`
	BatsmanID  int64 `json:"batsman_id"`
	TeamID     int64 `json:"team_id"`
	MatchID    int64 `json:"match_id"`
	Position   int32 `json:"position"`
	RunsScored int32 `json:"runs_scored"`
	BallsFaced int32 `json:"balls_faced"`
	Fours      int32 `json:"fours"`
	Sixes      int32 `json:"sixes"`
}

type Comment struct {
	ID          int64     `json:"id"`
	ThreadID    int64     `json:"thread_id"`
	Owner       string    `json:"owner"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}

type Community struct {
	ID              int64     `json:"id"`
	Owner           string    `json:"owner"`
	CommunitiesName string    `json:"communities_name"`
	Description     string    `json:"description"`
	CommunityType   string    `json:"community_type"`
	CreatedAt       time.Time `json:"created_at"`
}

type Communitymessage struct {
	ID             int64     `json:"id"`
	CommunityName  string    `json:"community_name"`
	SenderUsername string    `json:"sender_username"`
	Content        string    `json:"content"`
	SentAt         time.Time `json:"sent_at"`
}

type ContentAdmin struct {
	ID        int64  `json:"id"`
	ContentID int64  `json:"content_id"`
	Admin     string `json:"admin"`
}

type CricketScore struct {
	ID            int64  `json:"id"`
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	Inning        string `json:"inning"`
	Score         int32  `json:"score"`
	Wickets       int32  `json:"wickets"`
	Overs         int32  `json:"overs"`
	RunRate       string `json:"run_rate"`
	TargetRunRate string `json:"target_run_rate"`
}

type CricketToss struct {
	ID           int64  `json:"id"`
	MatchID      int64  `json:"match_id"`
	TossDecision string `json:"toss_decision"`
	TossWin      int64  `json:"toss_win"`
}

type Follow struct {
	ID             int64     `json:"id"`
	FollowerOwner  string    `json:"follower_owner"`
	FollowingOwner string    `json:"following_owner"`
	CreatedAt      time.Time `json:"created_at"`
}

type FootballIncident struct {
	ID                    int64  `json:"id"`
	MatchID               int64  `json:"match_id"`
	TeamID                *int64 `json:"team_id"`
	Periods               string `json:"periods"`
	IncidentType          string `json:"incident_type"`
	IncidentTime          int64  `json:"incident_time"`
	Description           string `json:"description"`
	CreatedAt             int64  `json:"created_at"`
	PenaltyShootoutScored bool   `json:"penalty_shootout_scored"`
}

type FootballIncidentPlayer struct {
	ID         int64 `json:"id"`
	IncidentID int64 `json:"incident_id"`
	PlayerID   int64 `json:"player_id"`
}

type FootballLineup struct {
	ID       int64  `json:"id"`
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

type FootballScore struct {
	ID         int64 `json:"id"`
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
	FirstHalf  int32 `json:"first_half"`
	SecondHalf int32 `json:"second_half"`
	Goals      int64 `json:"goals"`
}

type FootballStatistic struct {
	ID              int64 `json:"id"`
	MatchID         int64 `json:"match_id"`
	TeamID          int64 `json:"team_id"`
	ShotsOnTarget   int32 `json:"shots_on_target"`
	TotalShots      int32 `json:"total_shots"`
	CornerKicks     int32 `json:"corner_kicks"`
	Fouls           int32 `json:"fouls"`
	GoalkeeperSaves int32 `json:"goalkeeper_saves"`
	FreeKicks       int32 `json:"free_kicks"`
	YellowCards     int32 `json:"yellow_cards"`
	RedCards        int32 `json:"red_cards"`
}

type FootballSubstitution struct {
	ID       int64  `json:"id"`
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

type FootballSubstitutionsPlayer struct {
	ID          int64 `json:"id"`
	IncidentID  int64 `json:"incident_id"`
	PlayerInID  int64 `json:"player_in_id"`
	PlayerOutID int64 `json:"player_out_id"`
}

type Game struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	MinPlayers int32  `json:"min_players"`
}

type Goal struct {
	ID       int64 `json:"id"`
	MatchID  int64 `json:"match_id"`
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
	GoalTime int64 `json:"goal_time"`
}

type Group struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	TournamentID int64  `json:"tournament_id"`
	Strength     int32  `json:"strength"`
}

type JoinCommunity struct {
	ID            int64  `json:"id"`
	CommunityName string `json:"community_name"`
	Username      string `json:"username"`
}

type LikeThread struct {
	ID       int64  `json:"id"`
	ThreadID int64  `json:"thread_id"`
	Username string `json:"username"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Match struct {
	ID             int64  `json:"id"`
	TournamentID   int64  `json:"tournament_id"`
	AwayTeamID     int64  `json:"away_team_id"`
	HomeTeamID     int64  `json:"home_team_id"`
	StartTimestamp int64  `json:"start_timestamp"`
	EndTimestamp   int64  `json:"end_timestamp"`
	Type           string `json:"type"`
	StatusCode     string `json:"status_code"`
}

type Message struct {
	ID               int64     `json:"id"`
	Content          string    `json:"content"`
	IsSeen           bool      `json:"is_seen"`
	SenderUsername   string    `json:"sender_username"`
	ReceiverUsername string    `json:"receiver_username"`
	SentAt           time.Time `json:"sent_at"`
	MediaUrl         string    `json:"media_url"`
	MediaType        string    `json:"media_type"`
	IsDeleted        bool      `json:"is_deleted"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type Messagemedium struct {
	MessageID int64 `json:"message_id"`
	MediaID   int64 `json:"media_id"`
}

type Player struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Slug       string `json:"slug"`
	ShortName  string `json:"short_name"`
	MediaUrl   string `json:"media_url"`
	Positions  string `json:"positions"`
	Sports     string `json:"sports"`
	Country    string `json:"country"`
	PlayerName string `json:"player_name"`
	GameID     int64  `json:"game_id"`
}

type Profile struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Signup struct {
	MobileNumber string `json:"mobile_number"`
	Otp          string `json:"otp"`
}

type Team struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Shortname   string `json:"shortname"`
	Admin       string `json:"admin"`
	MediaUrl    string `json:"media_url"`
	Gender      string `json:"gender"`
	National    bool   `json:"national"`
	Country     string `json:"country"`
	Type        string `json:"type"`
	Sports      string `json:"sports"`
	PlayerCount int32  `json:"player_count"`
	GameID      int64  `json:"game_id"`
}

type TeamPlayer struct {
	TeamID      int64  `json:"team_id"`
	PlayerID    int64  `json:"player_id"`
	CurrentTeam string `json:"current_team"`
}

type TeamsGroup struct {
	ID           int64 `json:"id"`
	GroupID      int64 `json:"group_id"`
	TeamID       int64 `json:"team_id"`
	TournamentID int64 `json:"tournament_id"`
}

type Thread struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	CommunitiesName string    `json:"communities_name"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	MediaType       string    `json:"media_type"`
	MediaUrl        string    `json:"media_url"`
	LikeCount       int64     `json:"like_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type Tournament struct {
	ID             int64  `json:"id"`
	TournamentName string `json:"tournament_name"`
	Slug           string `json:"slug"`
	Sports         string `json:"sports"`
	Country        string `json:"country"`
	StatusCode     string `json:"status_code"`
	Level          string `json:"level"`
	StartTimestamp int64  `json:"start_timestamp"`
	GameID         int64  `json:"game_id"`
}

type TournamentStanding struct {
	StandingID     int64 `json:"standing_id"`
	TournamentID   int64 `json:"tournament_id"`
	GroupID        int64 `json:"group_id"`
	TeamID         int64 `json:"team_id"`
	Wins           int64 `json:"wins"`
	Loss           int64 `json:"loss"`
	Draw           int64 `json:"draw"`
	GoalFor        int64 `json:"goal_for"`
	GoalAgainst    int64 `json:"goal_against"`
	GoalDifference int64 `json:"goal_difference"`
	Points         int64 `json:"points"`
}

type TournamentTeam struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

type Uploadmedium struct {
	ID        int64     `json:"id"`
	MediaUrl  string    `json:"media_url"`
	MediaType string    `json:"media_type"`
	SentAt    time.Time `json:"sent_at"`
}

type User struct {
	Username       string  `json:"username"`
	MobileNumber   *string `json:"mobile_number"`
	HashedPassword *string `json:"hashed_password"`
	Role           string  `json:"role"`
	Gmail          *string `json:"gmail"`
}

type Wicket struct {
	ID            int64  `json:"id"`
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	BatsmanID     int64  `json:"batsman_id"`
	BowlerID      int64  `json:"bowler_id"`
	WicketsNumber int32  `json:"wickets_number"`
	WicketType    string `json:"wicket_type"`
	BallNumber    int32  `json:"ball_number"`
}

type GetPlayerByTeam struct {
	TeamID      int64  `json:"team_id"`
	PlayerID    int64  `json:"player_id"`
	CurrentTeam string `json:"current_team"`
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Slug        string `json:"slug"`
	ShortName   string `json:"short_name"`
	MediaUrl    string `json:"media_url"`
	Positions   string `json:"positions"`
	Sports      string `json:"sports"`
	Country     string `json:"country"`
	PlayerName  string `json:"player_name"`
	GameID      int64  `json:"game_id"`
}

type GetTeamByPlayer struct {
	TeamID      int64  `json:"team_id"`
	PlayerID    int64  `json:"player_id"`
	CurrentTeam string `json:"current_team"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Shortname   string `json:"shortname"`
	Admin       string `json:"admin"`
	MediaUrl    string `json:"media_url"`
	Gender      string `json:"gender"`
	National    bool   `json:"national"`
	Country     string `json:"country"`
	Type        string `json:"type"`
	Sports      string `json:"sports"`
	PlayerCount int32  `json:"player_count"`
	GameID      int64  `json:"game_id"`
}

type SearchTeam struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}