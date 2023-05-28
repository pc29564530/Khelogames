CREATE TABLE "signup" (
                          "mobile_number" string NOT NULL,
                          "otp" string NOT NULL,
                          "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "users" (
                         "username" varchar UNIQUE PRIMARY KEY,
                         "mobile_number" string UNIQUE NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "blogs" (
                         "id" bigserial PRIMARY KEY,
                         "username" varchar NOT NULL,
                         "title" varchar NOT NULL,
                         "content" text NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "communities" (
                               "id" bigserial PRIMARY KEY,
                               "communities_name" varchar NOT NULL,
                               "description" varchar NOT NULL,
                               "community_type" varchar NOT NULL,
                               "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE UNIQUE INDEX ON "blogs" ("username");

ALTER TABLE "blogs" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");