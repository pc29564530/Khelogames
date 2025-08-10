package models

// import "github.com/google/uuid"

// type NewTournament struct {
// 	ID             int64     `json:"id"`
// 	PublicID       uuid.UUID `json:"public_id"`
// 	UserID         int32     `json:"user_id"`
// 	GameID         int32     `json:"game_id"`
// 	Name           string    `json:"name"`
// 	Slug           string    `json:"slug"`
// 	Shortname      string    `json:"shortname"`
// 	Country        string    `json:"country"`
// 	LogoUrl        string    `json:"logo_url"`
// 	Status         string    `json:"status"`
// 	Level          string    `json:"level"`
// 	Season         int       `json:"season"`
// 	TournamentType string    `json:"tournament_type"`
// 	StartTimestamp int64     `json:"start_timestamp"`
// 	EndTimestamp   *int64    `json:"end_timestamp"`
// 	CurrentStage   string    `json:"current_stage"`
// }

// type TournamentParticipants struct {
// 	ID           int64     `json:"id"`
// 	PublicID     uuid.UUID `json:"public_id"`
// 	TournamentID int32     `json:"tournament_id"`
// 	GroupID      int32     `json:"group_id"`
// 	EntityID     int32     `json:"entity_id"`
// 	EntityType   string    `json:"entity_type"`
// 	SeedNumber   int       `json:"seed_number"`
// 	Status       string    `json:"status"`
// }

// type Team struct {
// 	ID          int64     `json:"id"`
// 	PublicID    uuid.UUID `json:"public_id"`
// 	GameID      int32     `json:"game_id"`
// 	UserID      int32     `json:"user_id"`
// 	Name        string    `json:"name"`
// 	Slug        string    `json:"slug"`
// 	Shortname   string    `json:"shortname"`
// 	MediaUrl    string    `json:"media_url"`
// 	Gender      string    `json:"gender"`
// 	Country     string    `json:"country"`
// 	City        string    `json:"city"`
// 	TeamType    string    `json:"team_type"`
// 	IsActive    bool      `json:"is_active"`
// 	PlayerCount int32     `json:"player_count"`
// }

// type Player struct {
// 	ID          int64     `json:"id"`
// 	PublicID    uuid.UUID `json:"public_id"`
// 	UserID      int32     `json:"user_id"`
// 	GameID      int32     `json:"game_id"`
// 	Name        string    `json:"name"`
// 	Slug        string    `json:"slug"`
// 	Shortname   string    `json:"shortname"`
// 	Position    string    `json:"position"`
// 	MediaUrl    string    `json:"media_url"`
// 	Nationality string    `json:"nationality"`
// 	IsActive    bool      `json:"is_active"`
// }

// type Match struct {
// 	ID              int64     `json:"id"`
// 	PublicID        uuid.UUID `json:"public_id"`
// 	TournamentID    int32     `json:"tournament_id"`
// 	GameID          int32     `json:"game_id"`
// 	Season          int       `json:"season"`
// 	HomeTeamID      int64     `json:"home_team_id"`
// 	AwayTeamID      int64     `json:"away_team_id"`
// 	StartTimestamp  int64     `json:"start_timestamp"`
// 	EndTimestamp    int64     `json:"end_timestamp"`
// 	Venue           string    `json:"venue"`
// 	City            string    `json:"city"`
// 	Country         string    `json:"country"`
// 	Type            string    `json:"type"`
// 	Status          string    `json:"status"`
// 	MatchFormat     *string   `json:"MatchFormat"`
// 	Stage           string    `json:"stage"`
// 	KnockoutLevelID *int32    `json:"KnockoutLevelID"`
// 	Result          int32     `json:"result"`
// }

// // type NewTournamentStage struct {
// // 	ID int64 `json:"id"`
// //     TournamentID int32 `json:"tournament_id"`
// // 	StageName string `json:"stage_name"`
// //     StageType string `json:"stage_type"`
// //     stage_name VARCHAR(50) NOT NULL, -- 'group_stage', 'round_of_16', 'quarter_finals', 'semi_finals', 'final'
// //     stage_type VARCHAR(20) NOT NULL, -- 'group', 'knockout'
// //     stage_order INTEGER NOT NULL,
// //     status VARCHAR(20) DEFAULT 'upcoming', -- 'upcoming', 'active', 'completed'
// //     start_date DATE,
// //     end_date DATE,
// //     description TEXT,
// //     stage_config JSONB, -- Stage-specific settings
// //     created_at TIMESTAMP DEFAULT NOW(),
// //     updated_at TIMESTAMP DEFAULT NOW(),
// //     UNIQUE(tournament_id, stage_order)
// // }
