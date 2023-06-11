CREATE TABLE "blogs" (
                         "id" bigserial PRIMARY KEY,
                         "username" varchar NOT NULL,
                         "title" varchar NOT NULL,
                         "content" text NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "users" (
                         "username" varchar UNIQUE NOT NULL,
                         "mobile_number" string UNIQUE NOT NULL,
                         "hashed_password" varchar UNIQUE NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "sessions" (
                            "id" uuid PRIMARY KEY,
                            "username" varchar NOT NULL,
                            "refresh_token" varchar NOT NULL,
                            "user_agent" varchar NOT NULL,
                            "client_ip" varchar NOT NULL,
                            "expires_at" timestamptz NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

CREATE TABLE "login" (
                         "username" varchar NOT NULL,
                         "password" varchar NOT NULL
);

CREATE TABLE "signup" (
                          "mobile_number" string NOT NULL,
                          "otp" string NOT NULL
);

CREATE TABLE "communities" (
                               "id" bigserial PRIMARY KEY,
                               "communities_name" varchar NOT NULL,
                               "description" varchar NOT NULL,
                               "community_type" varchar NOT NULL,
                               "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "friends" (
                           "id" bigserial,
                           "friend_username" varchar,
                           "friend_name" varchar
);

CREATE TABLE "friends_request" (
                                "id" bigserial PRIMARY KEY,
                                "sender_username" varchar NOT NULL,
                                "reciever_username" varchar NOT NULL,
                                "status" varchar NOT NULL,
                                "created_at" timestamp   NOT NULL DEFAULT (now())
);

ALTER TABLE "blogs" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "blogs" ADD FOREIGN KEY ("created_at") REFERENCES "blogs" ("username");

ALTER TABLE "friends" ADD FOREIGN KEY ("friend_username") REFERENCES "users" ("username");

ALTER TABLE "friends_request" ADD FOREIGN KEY ("reciever_username") REFERENCES "users" ("username");

ALTER TABLE "friends_request" ADD FOREIGN KEY ("sender_username") REFERENCES "users" ("username");