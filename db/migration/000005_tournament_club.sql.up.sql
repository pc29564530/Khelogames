CREATE TABLE "club" (
    "id" bigserial PRIMARY KEY,
    "club_name" varchar NOT NULL,
    "avatar_url" varchar NOT NULL,
    "sport" varchar NOT NULL,
    "owner" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "club_member" (
    "id" bigserial PRIMARY KEY,
    "club_name" varchar NOT NULL,
    "club_member" varchar NOT NULL,
    "joined_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "tournament" (
    tournament_id bigserial PRIMARY KEY,
    tournament_name varchar(255) NOT NULL,
    sport_type varchar(100) NOT NULL,
    format varchar(100) NOT NULL,
    teams_joined bigInt NOT NULL,
    CONSTRAINT format_check CHECK (format IN ('knockout', 'league', 'leagure+knockout', 'gourps+knockout', 'custom'))
);

CREATE TABLE organizer (
    organizer_id bigserial PRIMARY KEY,
    organizer_name varchar NOT NULL,
    tournament_id bigserial  NOT NULL
);

CREATE TABLE tournament_organizer (
    organizer_id bigserial REFERENCES organizer(organizer_id),
    tournament_id bigserial REFERENCES tournament(tournament_id),
    PRIMARY KEY (organizer_id, tournament_id)
);

CREATE TABLE tournament_team (
    tournament_id bigserial REFERENCES tournament(tournament_id),
    team_id bigserial REFERENCES club(id),
    PRIMARY KEY (tournament_id, team_id)
);

CREATE TABLE tournament_organization (
    id bigserial PRIMARY KEY,
    tournament_id bigint NOT NULL REFERENCES tournament (tournament_id),
    tournament_start timestamp NOT NULL DEFAULT 'now()',
    player_count bigint NOT NULL,
    team_count bigint NOT NULL,
    group_count bigint NOT NULL,
    advanced_team bigint NOT NULL
);

CREATE TABLE tournament_match (
    match_id bigserial PRIMARY KEY,
    organizer_id bigint NOT NULL REFERENCES "organizer" ("organizer_id"),
    tournament_id bigint NOT NULL REFERENCES "tournament" ("tournament_id"),
    team1_id bigint NOT NULL REFERENCES "club" ("id"),
    team2_id bigint NOT NULL REFERENCES "club" ("id"),
    date_on timestamp NOT NULL,
    start_at timestamp NOT NULL,
    stage varchar NOT NULL,
    created_at timestamp DEFAULT 'now()',
    sports varchar NOT NULL
);

CREATE TABLE tournament_standing (
    standing_id bigserial PRIMARY KEY,
    tournament_id bigserial NOT NULL REFERENCES tournament (tournament_id),
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

CREATE TABLE tournament_group (
    group_id bigserial PRIMARY KEY,
    tournament_id bigint NOT NULL REFERENCES tournament (tournament_id),
    team_id bigint NOT NULL REFERENCES club (id)
);

ALTER TABLE "club" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("club_member") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("club_name") REFERENCES "club" ("club_name");
ALTER TABLE "organizer" ADD FOREIGN KEY ("organizer_name") REFERENCES "users" ("username");
ALTER TABLE "organizer" ADD FOREIGN KEY ("tournament_id") REFERENCES "tournament" ("tournament_id");