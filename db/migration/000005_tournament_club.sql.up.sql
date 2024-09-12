-- organizer only work for the local or school tournament or matches
CREATE TABLE organizer (
    organizer_id bigserial PRIMARY KEY,
    organizer_name varchar NOT NULL,
    tournament_id bigserial  NOT NULL
);

-- we need the organizer for only the local event
CREATE TABLE tournament_organizer (
    organizer_id bigserial REFERENCES organizer(organizer_id),
    tournament_id bigserial REFERENCES tournaments(id),
    PRIMARY KEY (organizer_id, tournament_id)
);

CREATE TABLE tournament_team (
    tournament_id bigserial REFERENCES tournaments(id),
    team_id bigserial REFERENCES teams(id),
    PRIMARY KEY (tournament_id, team_id)
);

CREATE TABLE tournament_organization (
    id bigserial PRIMARY KEY,
    tournament_id bigint NOT NULL REFERENCES tournaments (id),
    tournament_start timestamp NOT NULL DEFAULT 'now()',
    player_count bigint NOT NULL,
    team_count bigint NOT NULL,
    group_count bigint NOT NULL,
    advanced_team bigint NOT NULL
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

CREATE TABLE group_team (
  group_team_id bigserial PRIMARY KEY,
  group_id bigserial NOT NULL,
  team_id bigserial NOT NULL,
  tournament_id bigserial NOT NULL
);

CREATE TABLE group_league (
  group_id bigserial PRIMARY KEY,
  group_name varchar NOT NULL,
  tournament_id bigserial NOT NULL,
  group_strength bigserial NOT NULL
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

ALTER TABLE "group_team" ADD FOREIGN KEY ("group_id") REFERENCES "group_league" ("group_id");
ALTER TABLE "group_team" ADD FOREIGN KEY ("team_id") REFERENCES "club" ("id");
ALTER TABLE "group_team" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");
ALTER TABLE "group_league" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");
ALTER TABLE "group_league" ADD FOREIGN KEY ("team_id") REFERENCES "club" ("id");
ALTER TABLE "organizer" ADD FOREIGN KEY ("organizer_name") REFERENCES "users" ("username");
ALTER TABLE "organizer" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");
