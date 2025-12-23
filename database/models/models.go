package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	FullName     string    `json:"full_name"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	HashPassword *string   `json:"hash_password"`
	IsVerified   bool      `json:"is_verified"`
	IsBanned     bool      `json:"is_banned"`
	GoogleID     *string   `json:"google_id"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BowlerScore struct {
	ID              int64     `json:"id"`
	PublicID        uuid.UUID `json:"public_id"`
	MatchID         int32     `json:"match_id"`
	TeamID          int32     `json:"team_id"`
	BowlerID        int32     `json:"bowler_id"`
	InningNumber    int       `json:"inning_number"`
	BallNumber      int       `json:"ball_number"`
	Runs            int32     `json:"runs"`
	Wickets         int32     `json:"wickets"`
	Wide            int32     `json:"wide"`
	NoBall          int32     `json:"no_ball"`
	BowlingStatus   bool      `json:"bowling_status"`
	IsCurrentBowler bool      `json:"is_current_bowler"`
}

type BatsmanScore struct {
	ID                 int64     `json:"id"`
	PublicID           uuid.UUID `json:"public_id"`
	MatchID            int32     `json:"match_id"`
	TeamID             int32     `json:"team_id"`
	BatsmanID          int32     `json:"batsman_id"`
	InningNumber       int       `json:"inning_number"`
	Position           string    `json:"position"`
	RunsScored         int       `json:"runs_scored"`
	BallsFaced         int       `json:"balls_faced"`
	Fours              int       `json:"fours"`
	Sixes              int       `json:"sixes"`
	BattingStatus      bool      `json:"batting_status"`
	IsStriker          bool      `json:"is_striker"`
	IsCurrentlyBatting bool      `json:"is_currently_batting"`
}

type CricketScore struct {
	ID                int64     `json:"id"`
	PublicID          uuid.UUID `json:"public_id"`
	MatchID           int32     `json:"match_id"`
	TeamID            int32     `json:"team_id"`
	InningNumber      int       `json:"inning_number"`
	Score             int       `json:"score"`
	Wickets           int       `json:"wickets"`
	Overs             int       `json:"overs"`
	RunRate           string    `json:"run_rate"`
	TargetRunRate     string    `json:"target_run_rate"`
	IsInningCompleted bool      `json:"is_inning_completed"`
	FollowOn          bool      `json:"follow_on"`
	Declared          bool      `json:"declared"`
	InningStatus      string    `json:"inning_status"`
}

type CricketToss struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	MatchID      int32     `json:"match_id"`
	TossDecision string    `json:"toss_decision"`
	TossWin      int32     `json:"toss_win"`
}

type Comment struct {
	ID              int64     `json:"id"`
	PublicID        uuid.UUID `json:"public_id"`
	ThreadID        int32     `json:"thread_id"`
	UserID          int32     `json:"user_id"`
	ParentCommentID *int32    `json:"parent_comment_id"`
	CommentText     string    `json:"comment_text"`
	LikeCount       int       `json:"like_count"`
	ReplyCount      int       `json:"reply_count"`
	IsDeleted       bool      `json:"is_deleted"`
	IsEdited        bool      `json:"is_edited"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Communities struct {
	ID            int64     `json:"id"`
	PublicID      uuid.UUID `json:"public_id"`
	UserID        int32     `json:"user_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CommunityType string    `json:"community_type"`
	IsActive      bool      `json:"is_active"`
	MemberCount   int       `json:"member_count"`
	AvatarUrl     string    `json:"avatar_url"`
	CoverImageUrl string    `json:"cover_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CommunityMessage struct {
	ID          int64     `json:"id"`
	PublicID    uuid.UUID `json:"public_id"`
	CommunityID int32     `json:"community_ud"`
	SenderID    int32     `json:"sender_id"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	MediaUrl    string    `json:"media_url"`
	MediaType   string    `json:"media_type"`
	SentAt      time.Time `json:"sent_at"`
}

type ContentAdmin struct {
	ID        int64  `json:"id"`
	ContentID int64  `json:"content_id"`
	Admin     string `json:"admin"`
}

type UsersConnections struct {
	UserID       int32 `json:"user_id"`
	TargetUserID int32 `json:"target_user_id"`
}

type Roles struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserRole struct {
	ID        int64     `json:"id"`
	UserID    int32     `json:"user_id"`
	RoleID    int32     `json:"role_id"`
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

type FootballIncident struct {
	ID                    int64     `json:"id"`
	PublicID              uuid.UUID `json:"public_id"`
	TournamentID          int32     `json:"tournament_id"`
	MatchID               int32     `json:"match_id"`
	TeamID                *int32    `json:"team_id"`
	Periods               string    `json:"periods"`
	IncidentType          string    `json:"incident_type"`
	IncidentTime          int       `json:"incident_time"`
	Description           string    `json:"description"`
	PenaltyShootoutScored bool      `json:"penalty_shootout_scored"`
	CreatedAt             int64     `json:"created_at"`
}

type FootballIncidentPlayer struct {
	ID         int64 `json:"id"`
	IncidentID int32 `json:"incident_id"`
	PlayerID   int32 `json:"player_id"`
}

type FootballScore struct {
	ID              int64     `json:"id"`
	PublicID        uuid.UUID `json:"public_id"`
	MatchID         int32     `json:"match_id"`
	TeamID          int32     `json:"team_id"`
	FirstHalf       int32     `json:"first_half"`
	SecondHalf      int32     `json:"second_half"`
	Goals           int       `json:"goals"`
	PenaltyShootOut *int      `json:"penalty_shootout"`
}

type FootballStatistics struct {
	ID              int64     `json:"id"`
	PublicID        uuid.UUID `json:"public_id"`
	MatchID         int32     `json:"match_id"`
	TeamID          int32     `json:"team_id"`
	ShotsOnTarget   int32     `json:"shots_on_target"`
	TotalShots      int32     `json:"total_shots"`
	CornerKicks     int32     `json:"corner_kicks"`
	Fouls           int32     `json:"fouls"`
	GoalkeeperSaves int32     `json:"goalkeeper_saves"`
	FreeKicks       int32     `json:"free_kicks"`
	YellowCards     int32     `json:"yellow_cards"`
	RedCards        int32     `json:"red_cards"`
}

type FootballSubstitutionsPlayer struct {
	ID          int64 `json:"id"`
	IncidentID  int32 `json:"incident_id"`
	PlayerInID  int32 `json:"player_in_id"`
	PlayerOutID int32 `json:"player_out_id"`
}

type Game struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	MinPlayers int32  `json:"min_players"`
}

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type JoinCommunity struct {
	ID          int64     `json:"id"`
	PublicID    uuid.UUID `json:"public_id"`
	CommunityID int32     `json:"community_id"`
	UserID      int32     `json:"user_id"`
}

type UserLikeThread struct {
	ID       int64  `json:"id"`
	ThreadID int32  `json:"thread_id"`
	UserID   string `json:"user_id"`
}

type Match struct {
	ID              int64     `json:"id"`
	PublicID        uuid.UUID `json:"public_id"`
	TournamentID    int32     `json:"tournament_id"`
	AwayTeamID      int32     `json:"away_team_id"`
	HomeTeamID      int32     `json:"home_team_id"`
	StartTimestamp  int32     `json:"start_timestamp"`
	EndTimestamp    int32     `json:"end_timestamp"`
	Type            string    `json:"type"`
	StatusCode      string    `json:"status_code"`
	Result          *int32    `json:"match"`
	Stage           string    `json:"stage"`
	KnockoutLevelID *int32    `json:"knockout_level_id"`
	MatchFormat     *string   `json:"match_format"`
	DayNumber       *int      `json:"day_number"`
	SubStatus       *string   `json:"sub_status"`
	LocationID      *int32    `json:"location_id"`
}

type Message struct {
	ID          int64     `json:"id"`
	PublicID    uuid.UUID `json:"public_id"`
	SenderID    int32     `json:"sender_id"`
	ReceiverID  int32     `json:"receiver_id"`
	Content     string    `json:"content"`
	MediaUrl    string    `json:"media_url"`
	MediaType   string    `json:"media_type"`
	IsSeen      bool      `json:"is_seen"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	SentAt      time.Time `json:"sent_at"`
	IsDelivered bool      `json:"is_delivered"`
}

// type Messagemedium struct {
// 	MessageID int64 `json:"message_id"`
// 	MediaID   int64 `json:"media_id"`
// }

type OldPlayer struct {
	ID        int64     `json:"id"`
	PublicID  uuid.UUID `json:"public_id"`
	UserID    int32     `json:"user_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ShortName string    `json:"short_name"`
	MediaUrl  string    `json:"media_url"`
	Positions string    `json:"positions"`
	Country   string    `json:"country"`
	GameID    int64     `json:"game_id"`
}

// Player - references User (many:1 with User)
type Player struct {
	ID        int64     `json:"id" db:"id"`
	PublicID  uuid.UUID `json:"public_id" db:"public_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	GameID    int64     `json:"game_id" db:"game_id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	ShortName string    `json:"short_name" db:"short_name"`
	MediaUrl  string    `json:"media_url" db:"media_url"`
	Positions string    `json:"positions" db:"positions"`
	Country   string    `json:"country" db:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserProfiles struct {
	ID         int64     `json:"id" db:"id"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Bio        string    `json:"bio" db:"bio"`
	AvatarUrl  string    `json:"avatar_url" db:"avatar_url"`
	Location   string    `json:"location" db:"location"`
	LocationID *int32    `json:"location_id" db:"location_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Session struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	UserID       int32     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// type Signup struct {
// 	MobileNumber string `json:"mobile_number"`
// 	Otp          string `json:"otp"`
// }

type Team struct {
	ID          int64     `json:"id"`
	PublicID    uuid.UUID `json:"public_id"`
	UserID      int32     `json:"user_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Shortname   string    `json:"shortname"`
	MediaUrl    string    `json:"media_url"`
	Gender      string    `json:"gender"`
	National    bool      `json:"national"`
	Country     string    `json:"country"`
	Type        string    `json:"type"`
	PlayerCount int32     `json:"player_count"`
	GameID      int64     `json:"game_id"`
}

type TournamentParticipants struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	TournamentID int32     `json:"tournament_id"`
	GroupID      *int32    `json:"group_id"`
	EntityID     int32     `json:"entity_id"`   //team or player
	EntityType   string    `json:"entity_type"` //team or player
	SeedNumber   *int      `json:"seed_number"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:created_at`
}

type TeamPlayer struct {
	TeamID    int32  `json:"team_id"`
	PlayerID  int32  `json:"player_id"`
	JoinDate  int32  `json:"join_date"`
	LeaveDate *int32 `json:"leave_date"`
}

type TeamsGroup struct {
	ID           int64 `json:"id"`
	GroupID      int64 `json:"group_id"`
	TeamID       int32 `json:"team_id"`
	TournamentID int32 `json:"tournament_id"`
}

type Thread struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	UserID       int32     `json:"user_id"`
	CommunityID  *int32    `json:"community_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	MediaUrl     string    `json:"media_url"`
	MediaType    string    `json:"media_type"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:comment_count"`
	IsDeleted    bool      `json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
}

type Tournament struct {
	ID             int64     `json:"id"`
	PublicID       uuid.UUID `json:"public_id"`
	UserID         int32     `json:"user_id"`
	GameID         int64     `json:"game_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Description    *string   `json:"description"`
	Country        string    `json:"country"`
	Status         string    `json:"status"`
	Season         *int      `json:"season"`
	Level          string    `json:"level"`
	StartTimestamp int64     `json:"start_timestamp"`
	GroupCount     *int      `json:"group_count"`
	MaxGroupTeam   *int      `json:"max_group_team"`
	Stage          string    `json:"stage"`
	HasKnockout    bool      `json:"has_knockout"`
	IsPublic       bool      `json:"is_public"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LocationID     *int32    `json:"location_id"`
}

type FootballStanding struct {
	ID             int64     `json:"id"`
	PublicID       uuid.UUID `json:"public_id"`
	TournamentID   int32     `json:"tournament_id"`
	GroupID        *int32    `json:"group_id"`
	TeamID         int32     `json:"team_id"`
	Matches        *int      `json:"matches"`
	Wins           *int      `json:"wins"`
	Loss           *int      `json:"loss"`
	Draw           *int      `json:"draw"`
	GoalFor        *int      `json:"goal_for"`
	GoalAgainst    *int      `json:"goal_against"`
	GoalDifference *int      `json:"goal_difference"`
	Points         *int      `json:"points"`
}

type CricketStanding struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	TournamentID int32     `json:"tournament_id"`
	GroupID      *int64    `json:"group_id"`
	TeamID       int32     `json:"team_id"`
	Matches      *int      `json:"matches"`
	Wins         *int      `json:"wins"`
	Loss         *int      `json:"loss"`
	Draw         *int      `json:"draw"`
	Points       *int      `json:"points"`
}

type TournamentTeam struct {
	TournamentID int32 `json:"tournament_id"`
	TeamID       int32 `json:"team_id"`
}

// type Uploadmedium struct {
// 	ID        int64     `json:"id"`
// 	MediaUrl  string    `json:"media_url"`
// 	MediaType string    `json:"media_type"`
// 	SentAt    time.Time `json:"sent_at"`
// }

type Wicket struct {
	ID            int64     `json:"id"`
	PublicID      uuid.UUID `json:"public_id"`
	MatchID       int32     `json:"match_id"`
	TeamID        int32     `json:"team_id"`
	BatsmanID     int32     `json:"batsman_id"`
	BowlerID      int32     `json:"bowler_id"`
	InningNumber  int       `json:"inning_number"`
	WicketsNumber int       `json:"wickets_number"`
	WicketType    string    `json:"wicket_type"`
	BallNumber    int       `json:"ball_number"`
	FielderID     *int32    `json:"fielder_id"`
	Score         *int      `json:"score"`
}

type GetPlayerByTeam struct {
	ID         int64     `json:"id"`
	PublicID   uuid.UUID `json:"public_id"`
	TeamID     int32     `json:"team_id"`
	PlayerID   int32     `json:"player_id"`
	JoinDate   *int32    `json:"join_date"`
	LeaveDate  *int32    `json:"leave_date"`
	Slug       string    `json:"slug"`
	ShortName  string    `json:"short_name"`
	MediaUrl   string    `json:"media_url"`
	Positions  string    `json:"positions"`
	Country    string    `json:"country"`
	PlayerName string    `json:"player_name"`
	GameID     int64     `json:"game_id"`
}

type GetTeamByPlayer struct {
	ID          int64     `json:"id"`
	PublicID    uuid.UUID `json:"public_id"`
	TeamID      int32     `json:"team_id"`
	PlayerID    int32     `json:"player_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Shortname   string    `json:"shortname"`
	MediaUrl    string    `json:"media_url"`
	Gender      string    `json:"gender"`
	National    bool      `json:"national"`
	Country     string    `json:"country"`
	Type        string    `json:"type"`
	PlayerCount int32     `json:"player_count"`
	GameID      int64     `json:"game_id"`
}

type SearchTeam struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FootballSquad struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	MatchID      *int32    `json:"match_id"`
	TeamID       int32     `json:"team_id"`
	PlayerID     int32     `json:"player_id"`
	Position     *string   `json:"position"`
	IsSubstitute bool      `json:"is_substitute"`
	Role         *string   `json:"role"`
	CreatedAT    time.Time `json:"created_at"`
}

type CricketSquad struct {
	ID        int64     `json:"id"`
	PublicID  uuid.UUID `json:"public_id"`
	MatchID   *int32    `json:"match_id"`
	TeamID    int32     `json:"team_id"`
	PlayerID  int32     `json:"player_id"`
	Role      *string   `json:"role"`
	OnBench   bool      `json:"on_bench"`
	IsCaptain bool      `json:"is_captain"`
	CreatedAT time.Time `json:"created_at"`
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
	ID         int64     `json:"id"`
	PublicID   uuid.UUID `json:"public_id"`
	PlayerID   int32     `json:"player_id"`
	MatchType  string    `json:"match_type"`
	Matches    int       `json:"matches"`
	Innings    int       `json:"innings"`
	Runs       int       `json:"runs"`
	Balls      int       `json:"balls"`
	Sixes      int       `json:"sixes"`
	Fours      int       `json:"fours"`
	Fifties    int       `json:"fifties"`
	Hundreds   int       `json:"hundreds"`
	BestScore  int       `json:"best_score"`
	Average    string    `json:"average"`
	StrikeRate string    `json:"strike_rate"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}

type PlayerBowlingStats struct {
	ID          int64     `json:"id"`
	PublicID    uuid.UUID `json:"public_id"`
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
	PublicID       uuid.UUID `json:"public_id"`
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

type CricketPlayerStats struct {
	ID             int64     `json:"id"`
	PublicID       uuid.UUID `json:"public_id"`
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

type CricketMatchInningDetails struct {
	ID            int64     `json:"id"`
	MatchID       int32     `json:"match_id"`
	BatsmanTeamID int32     `json:"batsman_team_id"`
	BowlerTeamID  int32     `json:"bowler_team_id"`
	InningNumber  int       `json:"inning_number"`
	InningStatus  string    `json:"inning_status"`
	LastUpdated   time.Time `json:"last_updated"`
}

type MatchHighlights struct {
	ID           int64     `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	UserID       int32     `json:"user_id"`
	TournamentID int32     `json:"tournament_id"`
	MatchID      int32     `json:"match_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	MediaURL     string    `json:"media_url"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

type TournamentUserRoles struct {
	ID           int64  `json:"id"`
	TournamentID int32  `json:"tournament_id"`
	UserID       int32  `json:"user_id"`
	Role         string `json:"role"`
}

type Locations struct {
	ID        int64     `json:"id"`
	PublicID  uuid.UUID `json:"public_id"`
	City      *string   `json:"city"`
	State     *string   `json:"state"`
	Country   *string   `json:"country"`
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
