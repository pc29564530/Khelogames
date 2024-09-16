CREATE TABLE tournament_team (
    tournament_id bigserial REFERENCES tournaments(id),
    team_id bigserial REFERENCES teams(id),
    PRIMARY KEY (tournament_id, team_id)
);

CREATE TABLE tournament_standing (
    standing_id bigserial PRIMARY KEY,
    tournament_id bigserial NOT NULL REFERENCES tournaments (id),
    group_id bigserial NOT NULL REFERENCES groups (id),
    team_id bigint NOT NULL REFERENCES teams (id),
    wins bigint NOT NULL,
    loss bigint NOT NULL,
    draw bigint NOT NULL,
    goal_for bigint NOT NULL,
    goal_against bigint NOT NULL,
    goal_difference bigint NOT NULL,
    points bigint NOT NULL
);

CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    tournament_id BIGSERIAL NOT NULL REFERENCES tournaments (id),
    strength INT NOT NULL
);

CREATE TABLE teams_group (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGSERIAL NOT NULL REFERENCES groups (id),
    team_id BIGSERIAL NOT NULL REFERENCES teams (id),
    tournament_id BIGSERIAL NOT NULL REFERENCES tournaments (id)
);

CREATE TABLE tournaments (
    id BIGSERIAL PRIMARY KEY,
    tournament_name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    sports VARCHAR NOT NULL,
    country VARCHAR NOT NULL,
    status_code VARCHAR NOT NULL,
    level VARCHAR(255) CHECK (level IN ('international', 'country', 'local')) NOT NULL,
    start_timestamp BIGINT NOT NULL
);

CREATE TABLE content_admin (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGSERIAL REFERENCES tournaments (id) NOT NULL,
    admin VARCHAR NOT NULL REFERENCES users (username)
);
