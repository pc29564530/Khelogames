package models

import (
	"time"

	"github.com/google/uuid"
)

type Ball struct {
	ID              int64 `json:"id"`
	TeamID          int64 `json:"team_id"`
	MatchID         int64 `json:"match_id"`
	BowlerID        int64 `json:"bowler_id"`
	Ball            int32 `json:"ball"`
	Runs            int32 `json:"runs"`
	Wickets         int32 `json:"wickets"`
	Wide            int32 `json:"wide"`
	NoBall          int32 `json:"no_ball"`
	BowlingStatus   bool  `json:"bowling_status"`
	IsCurrentBowler bool  `json:"is_current_bowler"`
	InningNumber    int   `json:"inning_number"`
}

type Bat struct {
	ID                 int64  `json:"id"`
	BatsmanID          int64  `json:"batsman_id"`
	TeamID             int64  `json:"team_id"`
	MatchID            int64  `json:"match_id"`
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
	ID                int64  `json:"id"`
	MatchID           int64  `json:"match_id"`
	TeamID            int64  `json:"team_id"`
	InningNumber      int    `json:"inning_number"`
	Score             int32  `json:"score"`
	Wickets           int32  `json:"wickets"`
	Overs             int32  `json:"overs"`
	RunRate           string `json:"run_rate"`
	TargetRunRate     string `json:"target_run_rate"`
	IsInningCompleted bool   `json:"is_inning_completed"`
	FollowOn          bool   `json:"follow_on"`
	Declared          bool   `json:"declared"`
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
	TournamentID          int32  `json:"tournament_id"`
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
	ID   int64  `json:"id"`
	Name string `json:"name"`
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
	ID              int64   `json:"id"`
	TournamentID    int64   `json:"tournament_id"`
	AwayTeamID      int64   `json:"away_team_id"`
	HomeTeamID      int64   `json:"home_team_id"`
	StartTimestamp  int64   `json:"start_timestamp"`
	EndTimestamp    int64   `json:"end_timestamp"`
	Type            string  `json:"type"`
	StatusCode      string  `json:"status_code"`
	Result          *int64  `json:"match"`
	Stage           string  `json:"stage"`
	KnockoutLevelID *int32  `json:"KnockoutLevelID"`
	MatchFormat     *string `json:"MatchFormat"`
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
	Country    string `json:"country"`
	PlayerName string `json:"player_name"`
	GameID     int64  `json:"game_id"`
	ProfileID  *int32 `json:"profile_id"`
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
	PlayerCount int32  `json:"player_count"`
	GameID      int64  `json:"game_id"`
}

type TeamPlayer struct {
	TeamID    int64  `json:"team_id"`
	PlayerID  int64  `json:"player_id"`
	JoinDate  int32  `json:"join_date"`
	LeaveDate *int32 `json:"leave_date"`
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
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Country        string `json:"country"`
	StatusCode     string `json:"status_code"`
	Level          string `json:"level"`
	StartTimestamp int64  `json:"start_timestamp"`
	GameID         *int64 `json:"game_id"`
	GroupCount     *int32 `json:"group_count"`
	MaxGroupTeam   *int32 `json:"max_group_team"`
	Stage          string `json:"stage"`
	HasKnockout    bool   `json:"has_knockout"`
}

type FootballStanding struct {
	ID             int64  `json:"id"`
	TournamentID   int64  `json:"tournament_id"`
	GroupID        *int64 `json:"group_id"`
	TeamID         int64  `json:"team_id"`
	Matches        *int64 `json:"matches"`
	Wins           *int64 `json:"wins"`
	Loss           *int64 `json:"loss"`
	Draw           *int64 `json:"draw"`
	GoalFor        *int64 `json:"goal_for"`
	GoalAgainst    *int64 `json:"goal_against"`
	GoalDifference *int64 `json:"goal_difference"`
	Points         *int64 `json:"points"`
}

type CricketStanding struct {
	ID           int64  `json:"id"`
	TournamentID int64  `json:"tournament_id"`
	GroupID      *int64 `json:"group_id"`
	TeamID       int64  `json:"team_id"`
	Matches      *int64 `json:"matches"`
	Wins         *int64 `json:"wins"`
	Loss         *int64 `json:"loss"`
	Draw         *int64 `json:"draw"`
	Points       *int64 `json:"points"`
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
	Username     string  `json:"username"`
	MobileNumber *string `json:"mobile_number"`
	Role         string  `json:"role"`
	Gmail        *string `json:"gmail"`
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
	FielderID     *int64 `json:"fielder_id"`
	Score         *int32 `json:"score"`
	InningNumber  int    `json:"inning_number"`
}

type GetPlayerByTeam struct {
	TeamID     int64  `json:"team_id"`
	PlayerID   int64  `json:"player_id"`
	JoinDate   *int32 `json:"join_date"`
	LeaveDate  *int32 `json:"leave_date"`
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Slug       string `json:"slug"`
	ShortName  string `json:"short_name"`
	MediaUrl   string `json:"media_url"`
	Positions  string `json:"positions"`
	Country    string `json:"country"`
	PlayerName string `json:"player_name"`
	GameID     int64  `json:"game_id"`
	ProfileID  *int64 `json:"profile_id"`
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
	PlayerCount int32  `json:"player_count"`
	GameID      int64  `json:"game_id"`
}

type SearchTeam struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FootballSquad struct {
	ID           int64     `json:"id"`
	MatchID      *int64    `json:"match_id"`
	TeamID       int64     `json:"team_id"`
	PlayerID     int64     `json:"player_id"`
	Position     *string   `json:"position"`
	IsSubstitute bool      `json:"is_substitute"`
	Role         *string   `json:"role"`
	CreatedAT    time.Time `json:"created_at"`
}

type CricketSquad struct {
	ID        int64     `json:"id"`
	MatchID   *int64    `json:"match_id"`
	TeamID    int64     `json:"team_id"`
	PlayerID  int64     `json:"player_id"`
	Role      *string   `json:"role"`
	OnBench   bool      `json:"on_bench"`
	IsCaptain bool      `json:"is_captain"`
	CreatedAT time.Time `json:"created_at"`
}

type Roles struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserRole struct {
	ID        int64     `json:"id"`
	ProfileID int64     `json:"profile_id"`
	RoleID    int64     `json:"role_id"`
	CreatedAT time.Time `json:"created_at"`
}

type Organizations struct {
	ID               int64      `json:"id"`
	ProfileID        int64      `json:"profile_id"`
	OrganizationName string     `json:"organization_name"`
	Email            string     `json:"email"`
	PhoneNumber      string     `json:"phone_number"`
	IsVerified       bool       `json:"is_verified"`
	VerifiedAT       *time.Time `json:"vereified_at"`
	CreatedAT        time.Time  `json:"created_at"`
}

type Document struct {
	ID           int64     `json:"id"`
	OrganizerID  int32     `json:"organizer_id"`
	DocumentType string    `json:"document_type"`
	FilePath     string    `json:"file_path"`
	SubmittedAT  time.Time `json:"submitted_at"`
	Status       string    `json:"status"`
}

type PlayerBattingStats struct {
	ID           int64     `json:"id"`
	PlayerID     int32     `json:"player_id"`
	MatchType    string    `json:"match_type"`
	TotalMatches int       `json:"total_matches"`
	TotalInnings int       `json:"total_innings"`
	Runs         int       `json:"runs"`
	Balls        int       `json:"balls"`
	Sixes        int       `json:"sixes"`
	Fours        int       `json:"fours"`
	Fifties      int       `json:"fifties"`
	Hundreds     int       `json:"hundreds"`
	BestScore    int       `json:"best_score"`
	Average      string    `json:"average"`
	StrikeRate   string    `json:"strike_rate"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

type PlayerBowlingStats struct {
	ID          int64     `json:"id"`
	PlayerID    int32     `json:"player_id"`
	MatchType   string    `json:"match_type"`
	Matches     int       `json:"total_matches"`
	Innings     int       `json:"total_innings"`
	Wickets     int       `json:"wickets"`
	Runs        int       `json:"runs"`
	Balls       int       `json:"balls"`
	Average     string    `json:"average"`
	StrikeRate  string    `json:"strike_rate"`
	EconomyRate string    `json:"economy_rate"`
	FourWickets int       `json:"four_wickets"`
	FiveWickets int       `json:"five_wickets"`
	CreatedAT   time.Time `json:"created_at"`
	UpdatedAT   time.Time `json:"updated_at"`
}

type FootballPlayerStats struct {
	ID             int64     `json:"id"`
	PlayerID       int32     `json:"player_id"`
	PlayerPosition string    `json:"player_position"`
	Matches        int       `json"matches"`
	MinutesPlayed  int       `json:"minutes_played"`
	GoalsScored    int       `json:"goals_scored"`
	GoalsConceded  int       `json:"goals_conceded"`
	CleanSheet     int       `json:"CleanSheet"`
	Assists        int       `json:"assists"`
	YellowCards    int       `json:"yellow_cards"`
	RedCards       int       `json:"red_cards"`
	Average        string    `json:"average"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PlayerCricketStats struct {
	ID             int64     `json:"id"`
	PlayerID       int32     `json:"player_id"`
	MatchType      string    `json:"match_type"`
	Matches        int32     `json:"matches"`
	BattingInnings int32     `json:"batting_innings"`
	Runs           int       `json:"runs"`
	Balls          int       `json:"balls"`
	Sixes          int       `json:"sixes"`
	Fours          int       `json:"fours"`
	Fifties        int       `json:"fifties"`
	Hundreds       int       `json:"hundreds"`
	BestScore      int       `json:"best_score"`
	BowlingInnings int32     `json:"bowling_innings"`
	Wickets        int       `json:"wickets"`
	RunsConceded   int       `json:"runs_conceded"`
	BallsBowled    int       `json:"balls_bowled"`
	FourWickets    int       `json:"four_wickets"`
	FiveWickets    int       `json:"five_wickets"`
	CreatedAT      time.Time `json:"created_at"`
	UpdatedAT      time.Time `json:"updated_at"`
}
