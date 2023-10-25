CREATE TABLE "profile" (
                   "id" bigInt PRIMARY KEY,
                   "owner" varchar NOT NULL UNIQUE,
                   "full_name" varchar NOT NULL,
                   "bio" text NOT NULL,
                   "following_owner" bigInt NOT NULL,
                   "follower_owner" bigInt NOT NULL,
                   "avatar_url" varchar NOT NULL,
                   "created_at" timestamp NOT NULL DEFAULT 'now()'
);

ALTER TABLE "profile" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");