CREATE TABLE football_matches_score (
    id bigserial PRIMARY KEY,
    match_id bigserial NOT NULL REFERENCES tournament_match (match_id),
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
    team_id bigserial NOT NULL REFERENCES club (id),
    goal_for bigint NOT NULL,
    goal_against bigserial NOT NULL,
    goal_score_time timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE cricket_match_score (
    id bigserial PRIMARY KEY,
    match_id bigserial NOT NULL REFERENCES tournament_match (match_id),
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
    team_id bigserial NOT NULL REFERENCES club (id),
    score bigint NOT NULL,
    wickets bigint NOT NULL,
    overs bigint NOT NULL,
    extras bigint NOT NULL,
    innings bigint NOT NULL
);

CREATE TABLE cricket_team_player_score (
    id bigserial PRIMARY KEY,
    match_id bigserial NOT NULL REFERENCES tournament_match(match_id),
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
    team_id bigserial NOT NULL REFERENCES club(id),
    batting_or_bowling VARCHAR(10) NOT NULL,
    position bigint NOT NULL,
    player_id bigserial NOT NULL REFERENCES player_profile(id),
    runs_scored bigint NOT NULL,
    balls_faced bigint NOT NULL,
    fours bigint NOT NULL,
    sixes bigint NOT NULL,
    wickets_taken bigint NOT NULL,
    overs_bowled DECIMAL(4,1) NOT NULL,
    runs_conceded bigint NOT NULL,
    wicket_taken_by bigint NOT NULL,
    wicket_of bigint NOT NULL
);

CREATE TABLE football_team_player_score (
    id bigserial PRIMARY KEY,
    match_id bigserial NOT NULL REFERENCES tournament_match(match_id),
    team_id bigserial NOT NULL REFERENCES club(id),
    player_id bigserial NOT NULL REFERENCES player_profile (id),
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
    goal_score_time timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE cricket_match_team_toss (
    id bigserial PRIMARY KEY,
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
    match_id bigserial NOT NULL REFERENCES tournament_match (match_id),
    toss_won bigint NOT NULL REFERENCES club (id),
    bat_or_bowl varchar NOT NULL
);