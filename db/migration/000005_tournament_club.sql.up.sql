-- some chage in club added country and gender
CREATE TABLE "club" (
    "id" bigserial PRIMARY KEY,
    "club_name" varchar NOT NULL,
    "avatar_url" varchar NOT NULL,
    "sport" varchar NOT NULL,
    "owner" varchar NOT NULL, -- remove the owner club cannot have the owner it only applicable for the local or school club
    "created_at" timestamp NOT NULL DEFAULT 'now()',
    "country" varchar NOT NULL,
    "gender" varchar NOT NULL
);
-- squad of the team
CREATE TABLE "club_member" (
    "id" bigserial PRIMARY KEY,
    "club_id" bigserial NOT NULL,
    "player_id" bigserial NOT NULL
);

CREATE TABLE club_played (
    played_id bigserial PRIMARY KEY,
    tournament_id bigserial NOT NULL REFERENCES tournaments (id),
    club_id bigserial NOT NULL REFERENCES club (id)
);

-- CREATE TABLE "tournament" (
--     tournament_id bigserial PRIMARY KEY,
--     tournament_name varchar(255) NOT NULL,
--     sport_type varchar(100) NOT NULL,
--     format varchar(100) NOT NULL,
--     teams_joined bigInt NOT NULL,
--     start_on timestamp NOT NULL,
--     end_on timestamp NOT NULL,
--     category varchar NOT NULL DEFAULT 'Global',
--     CONSTRAINT format_check CHECK (format IN ('group', 'league', 'custom'))
-- );

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

-- remove the end_time and date_on
CREATE TABLE tournament_match (
    match_id bigserial PRIMARY KEY,
    organizer_id bigint NOT NULL REFERENCES organizer (organizer_id),
    tournament_id bigint NOT NULL REFERENCES tournaments (id),
    team1_id bigint NOT NULL REFERENCES club (id),
    team2_id bigint NOT NULL REFERENCES club (id),
    start_time timestamp NOT NULL,
    stage varchar NOT NULL,
    sports varchar NOT NULL
);

CREATE TABLE tournament_standing (
    standing_id bigserial PRIMARY KEY,
    tournament_id bigserial NOT NULL REFERENCES tournaments (id),
    group_id bigserial NOT NULL REFERENCES tournament_group (group_id),
    team_id bigint NOT NULL REFERENCES club (id),
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
-- create a new cricket matches
CREATE TABLE cricket_matches (
    id bigserial PRIMARY KEY,
    tournament_id bigserial REFERENCES tournaments(id),
    home_team_id bigserial REFERENCES teams(id),
    away_team_id bigserial REFERENCES teams(id),
    venue VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    start_timestamp TIMESTAMP NOT NULL,
    end_timestamp TIMESTAMP,
    toss_winner VARCHAR,
    toss_decision VARCHAR,
    status_code INT,
    status_description VARCHAR,
    current_period VARCHAR,
    tv_umpire_name VARCHAR,
    umpire1_name VARCHAR,
    umpire2_name VARCHAR
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

-- create a new football matches

ALTER TABLE "group_team" ADD FOREIGN KEY ("group_id") REFERENCES "group_league" ("group_id");
ALTER TABLE "group_team" ADD FOREIGN KEY ("team_id") REFERENCES "club" ("id");
ALTER TABLE "group_team" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");
ALTER TABLE "group_league" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");
ALTER TABLE "group_league" ADD FOREIGN KEY ("team_id") REFERENCES "club" ("id");
ALTER TABLE "club" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("club_member") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("club_name") REFERENCES "club" ("club_name");
ALTER TABLE "organizer" ADD FOREIGN KEY ("organizer_name") REFERENCES "users" ("username");
ALTER TABLE "organizer" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournaments" ("id");
