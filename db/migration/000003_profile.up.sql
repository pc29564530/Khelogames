-- divide the profile into two cases user and player

CREATE TABLE "profile" (
                   "id" bigserial PRIMARY KEY,
                   "owner" varchar UNIQUE NOT NULL,
                   "full_name" varchar NOT NULL,
                   "bio" text NOT NULL,
                   "avatar_url" varchar NOT NULL,
                   "cover_url" varchar NOT NULL,
                   "created_at" timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE player_profile (
    id bigserial PRIMARY KEY,
    player_name varchar NOT NULL,
    player_avatar_url varchar NOT NULL,
    player_bio varchar NOT NULL,
    player_sport varchar NOT NULL,
    player_playing_category varchar NOT NULL,
    nation varchar NOT NULL
);

ALTER TABLE "profile" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");