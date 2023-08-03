CREATE TABLE "threads" (
                           "id" bigserial PRIMARY KEY,
                           "username" varchar NOT NULL,
                           "communities_name" varchar NULL,
                           "title" varchar NOT NULL,
                           "content" text NOT NULL,
                           "media_type" varchar,
                           "media_url" text,
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

CREATE TABLE "login" (
                         "username" varchar NOT NULL,
                         "password" varchar NOT NULL
);
-- CREATE TABLE "signin" (
--     "mobile_number" string NOT NULL,
--     "otp" string NOT NULL
-- );

CREATE TABLE "signup" (
                          "mobile_number" string UNIQUE NOT NULL,
                          "otp" string NOT NULL
);

CREATE TABLE "communities" (
                               "id" bigserial PRIMARY KEY,
                               "owner" varchar NOT NULL,
                               "communities_name" varchar NOT NULL,
                               "description" varchar NOT NULL,
                               "community_type" varchar NOT NULL,
                               "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "friends" (
                           "id" bigserial PRIMARY KEY,
                           "owner" varchar NOT NULL,
                           "friend_username" varchar NOT NULL
);

CREATE TABLE "friends_request" (
                                "id" bigserial PRIMARY KEY,
                                "sender_username" varchar NOT NULL,
                                "reciever_username" varchar NOT NULL,
                                "status" varchar NOT NULL,
                                "created_at" timestamp   NOT NULL DEFAULT (now())
);

-- CREATE TABLE "community_member" (
--                                     "communities_name" varchar NOT NULL,
--                                     "username" varchar NOT NULL
-- );

-- CREATE TABLE "search_bar" (
--                               "full_name" varchar NOT NULL ,
--                               "username" varchar NOT NULL ,
--                               "communities" varchar NOT NULL
-- );

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "threads" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "threads" ADD FOREIGN KEY ("communities_name") REFERENCES "communities" ("communities_name");

ALTER TABLE "communities" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "friends" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "friends" ADD FOREIGN KEY ("friend_username") REFERENCES "users" ("username");

ALTER TABLE "friends_request" ADD FOREIGN KEY ("reciever_username") REFERENCES "users" ("username");

ALTER TABLE "friends_request" ADD FOREIGN KEY ("sender_username") REFERENCES "users" ("username");

-- ALTER TABLE "search_bar" ADD FOREIGN KEY ("full_name") REFERENCES "users" ("full_name");
--
-- ALTER TABLE "search_bar" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
--
-- ALTER TABLE "search_bar" ADD FOREIGN KEY ("communities") REFERENCES "communities" ("communities_name")