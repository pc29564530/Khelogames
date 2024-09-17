
CREATE TABLE teams (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    shortName VARCHAR(255) NOT NULL,
    admin VARCHAR(255) REFERENCES users(username) NOT NULL,
    media_url VARCHAR(255) NOT NULL,
    gender CHAR(1) CHECK (gender IN ('M', 'F')) NOT NULL,
    national BOOLEAN NOT NULL,
    country VARCHAR(255) NOT NULL,
    type VARCHAR(255) CHECK (type IN ('team', 'individual', 'double')) NOT NULL,
    sports VARCHAR(255) NOT NULL,
    games_id BIGINT REFERENCES games (id) NOT NULL,
    player_count INT NOT NULL
);

CREATE TABLE team_players (
    team_id BIGSERIAL REFERENCES teams(id),
    player_id BIGSERIAL REFERENCES players(id),
    current_team VARCHAR(2) NOT NULL CHECK(current_team IN ('t', 'f')),
    PRIMARY KEY (team_id, player_id)

);

CREATE TABLE matches (
    id BIGSERIAL PRIMARY KEY,
    tournament_id BIGSERIAL REFERENCES Tournaments(id) NOT NULL,
    away_team_id BIGSERIAL NOT NULL,
    home_team_id BIGSERIAL NOT NULL,
    start_timestamp BIGINT NOT NULL,
    end_timestamp BIGINT NOT NULL,
    type VARCHAR(255) CHECK(type IN ('team', 'individual', 'double')) NOT NULL,
    status_code VARCHAR(255) NOT NULL
);
	
CREATE TABLE cricket_score (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL NOT NULL REFERENCES matches(id),
    team_id BIGSERIAL NOT NULL,
    inning VARCHAR CHECK (inning IN ('inning1', 'inning2')) NOT NULL,
    score INT CHECK (score >= 0) NOT NULL,
    wickets INT CHECK (wickets BETWEEN 0 AND 10) NOT NULL,
    overs INT CHECK (overs>=0) NOT NULL,
    run_rate DECIMAL(5,2) CHECK (run_rate >= 0) NOT NULL,
    target_run_rate DECIMAL(5,2) CHECK (target_run_rate >= 0) NOT NULL
);



CREATE TABLE cricket_toss (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL NOT NULL REFERENCES matches(id),
    toss_decision VARCHAR CHECK (toss_decision IN ('Batting', 'Bowling')) NOT NULL,
    toss_win BIGINT NOT NULL REFERENCES teams(id)
);

CREATE TABLE football_score (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL NOT NULL,
    team_id BIGSERIAL NOT NULL,
    first_half INT NOT NULL,
    second_half INT NOT NULL,
    goals BIGINT NOT NULL
);

CREATE TABLE goals (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL REFERENCES matches(id) NOT NULL,
    team_id BIGSERIAL REFERENCES teams(id) NOT NULL,
    player_id BIGSERIAL REFERENCES players(id) NOT NULL,
    goal_time BIGINT NOT NULL
);
	
CREATE TABLE bats (
    id BIGSERIAL PRIMARY KEY,
    batsman_id BIGSERIAL REFERENCES players(id) NOT NULL,
    team_id BIGSERIAL REFERENCES teams(id) NOT NULL,
    match_id BIGSERIAL NOT NULL REFERENCES matches(id),
    position INT NOT NULL,
    runs_scored INT NOT NULL,
    balls_faced INT NOT NULL,
    fours INT NOT NULL,
    sixes INT NOT NULL
);
	
CREATE TABLE balls (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGSERIAL REFERENCES teams(id) NOT NULL,
    match_id BIGSERIAL REFERENCES matches(id) NOT NULL,
    bowler_id BIGSERIAL REFERENCES players(id) NOT NULL,
    ball INT NOT NULL,
    runs INT NOT NULL,
    wickets INT NOT NULL,
    wide INT NOT NULL,
    no_ball INT NOT NULL
);

CREATE TABLE wickets (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL REFERENCES matches(id) NOT NULL,
    team_id BIGSERIAL REFERENCES teams(id) NOT NULL,
    batsman_id BIGSERIAL REFERENCES players(id) NOT NULL,
    bowler_id BIGSERIAL REFERENCES players(id) NOT NULL,
    wickets_number INT NOT NULL,
    wicket_type VARCHAR(50) NOT NULL,
    ball_number INT NOT NULL
);

CREATE TABLE players (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) REFERENCES users (username) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    short_name VARCHAR(255) NOT NULL,
    media_url VARCHAR(255) NOT NULL,
    positions VARCHAR(3) NOT NULL,
    sports VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    player_name VARCHAR(255) NOT NULL
);

CREATE TABLE football_statistics (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL REFERENCES matches (id) NOT NULL,
    team_id BIGSERIAL REFERENCES teams (id) NOT NULL,
    shots_on_target INT NOT NULL,
    total_shots INT NOT NULL,
    corner_kicks INT NOT NULL,
    fouls INT NOT NULL,
    goalkeeper_saves INT NOT NULL,
    free_kicks INT NOT NULL,
    yellow_cards INT NOT NULL,
    red_cards INT NOT NULL
);

CREATE TABLE football_incidents (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL REFERENCES matches (id) NOT NULL,
    team_id BIGSERIAL REFERENCES teams (id) NOT NULL,
    periods VARCHAR(50) CHECK (periods IN ('first_half', 'second_half', 'extra_first_half', 'extra_second_half')) NOT NULL,
    incident_type VARCHAR(50) NOT NULL,
    incident_time BIGINT NOT NULL,
    description VARCHAR NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) NOT NULL
);


CREATE TABLE football_substitutions_player (
    id BIGSERIAL PRIMARY KEY,
    incident_id BIGSERIAL REFERENCES football_incidents (id) NOT NULL,
    player_in_id BIGSERIAL REFERENCES players (id) NOT NULL,
    player_out_id BIGSERIAL REFERENCES players (id) NOT NULL
);

CREATE TABLE football_incident_player (
    id BIGSERIAL PRIMARY KEY,
    incident_id BIGSERIAL REFERENCES football_incidents (id) NOT NULL,
    player_id BIGSERIAL REFERENCES players (id) NOT NULL
);

-- this will we run by normal query 
CREATE TABLE games (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    min_players INT NOT NULL
);