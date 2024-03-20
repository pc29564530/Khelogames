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

ALTER TABLE "club" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("club_member") REFERENCES "users" ("username");
ALTER TABLE "club_member" ADD FOREIGN KEY ("club_name") REFERENCES "club" ("club_name");