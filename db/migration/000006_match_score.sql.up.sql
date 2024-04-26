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
    score_id bigserial PRIMARY KEY,
    match_id bigserial NOT NULL REFERENCES tournament_match (match_id),
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
    team_id bigserial NOT NULL REFERENCES club (id),
    score bigint NOT NULL,
    wickets bigint NOT NULL,
    overs bigint NOT NULL,
    extras bigint NOT NULL,
    innings bigint NOT NULL
);