CREATE TABLE football_matches_score (
    id bigserial PRIMARY KEY,
    match_id bigserial NOT NULL REFERENCES tournament_match (match_id),
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
    team_id bigserial NOT NULL REFERENCES club (id),
    goal_score bigInt NOT NULL,
    goal_score_time timestamp NOT NULL DEFAULT 'now()'
);