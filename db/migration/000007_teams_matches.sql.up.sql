
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
    sports VARCHAR(255) NOT NULL
);

CREATE TABLE team_players (
    team_id BIGSERIAL REFERENCES teams(id),
    player_id BIGSERIAL REFERENCES players(id),
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
    status_code JSONB NOT NULL
);
	
CREATE TABLE cricket_score (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL NOT NULL REFERENCES matches(id),
    team_id BIGSERIAL NOT NULL,
    inning INT CHECK (inning BETWEEN 1 AND 4) NOT NULL,
    score INT CHECK (score >= 0) NOT NULL,
    wickets INT CHECK (wickets BETWEEN 0 AND 10) NOT NULL,
    overs DECIMAL(5,1) CHECK (overs >= 0) NOT NULL,
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
	
CREATE TABLE Bats (
    id BIGSERIAL PRIMARY KEY,
    batsman_id BIGSERIAL REFERENCES Players(id),
    match_id BIGSERIAL ,
    runs_scored BIGINT,
    balls_faced BIGINT,
    wickets_taken BIGINT
);
	
CREATE TABLE Balls (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGSERIAL,
    over_number BIGINT,
    ball_number BIGINT,
    runs BIGINT,
    wickets BIGINT,
    batsman_id BIGSERIAL REFERENCES Players(id),
    bowler_id BIGSERIAL REFERENCES Players(id),
    fielder_id BIGSERIAL REFERENCES Players(id)
);

-- CREATE TABLE Score (
-- 	id BIGSERIAL PRIMARY KEY,
-- 	team_id BIGSERIAL,
-- 	match_id BIGSERIAL,
-- 	score BIGINT,
-- 	period1 BIGINT,
-- 	period2 BIGINT
-- );

-- CREATE A INCIDENT TABLE
-- CREATE A STATUS TABLE AND STATUS TIME TABLE
-- this for football and other sport that played on the timer
-- CREATE TABLE Status_Time (
-- 	sport varchar(255), -- [football, hockey]
-- 	current_time bigint,
-- 	initial bigint,
-- 	max bigint,
-- 	extra bigint
-- );

-- create a table of the time

-- create a table for player
CREATE TABLE players (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) REFERENCES users (username) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    short_name VARCHAR(255) NOT NULL,
    media_url VARCHAR(255) NOT NULL,
    positions VARCHAR(3) NOT NULL,
    sports VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL
);