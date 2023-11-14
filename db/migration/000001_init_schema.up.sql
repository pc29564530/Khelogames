CREATE TABLE "threads" (
                           "id" bigserial PRIMARY KEY,
                           "username" varchar NOT NULL,
                           "communities_name" varchar NOT NULL ,
                           "title" varchar NOT NULL,
                           "content" text NOT NULL,
                           "media_type" varchar NOT NULL ,
                           "media_url" text NOT NULL ,
                           "like_count" bigint NOT NULL,
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



-- CREATE TABLE "community_member" (
--                                     "communities_name" varchar NOT NULL,
--                                     "username" varchar NOT NULL
-- );

-- CREATE TABLE "search_bar" (
--                               "full_name" varchar NOT NULL ,
--                               "username" varchar NOT NULL ,
--                               "communities" varchar NOT NULL
-- );

-- CREATE TABLE "likes" (
--                         "thread_id" bigserial PRIMARY KEY,
--                         "count" bigInt NOT NULL
-- );
--

CREATE TABLE "follow" (
                          "id" bigserial PRIMARY KEY,
                          "follower_owner" varchar NOT NULL,
                          "following_owner" varchar NOT NULL,
                          "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "comment" (
                           "id" bigserial PRIMARY KEY,
                           "thread_id" bigserial NOT NULL,
                           "owner" varchar NOT NULL,
                           "comment_text" text NOT NULL,
                           "created_at" timestamp NOT NULL DEFAULT 'now()'
);

ALTER TABLE "like_thread" ADD FOREIGN KEY ("thread_id") REFERENCES "threads" ("id");

ALTER TABLE "like_thread" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "comment" ADD FOREIGN KEY ("thread_id") REFERENCES "threads" ("id");

ALTER TABLE "comment" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "follow" ADD FOREIGN KEY ("follower_owner") REFERENCES "users" ("username");

ALTER TABLE "follow" ADD FOREIGN KEY ("following_owner") REFERENCES "users" ("username");

--
-- ALTER TABLE "likes" ADD FOREIGN KEY ("thread_id") REFERENCES "threads" ("id");
--
-- ALTER TABLE "comment" ADD FOREIGN KEY ("thread_id") REFERENCES "threads" ("id");
--
-- ALTER TABLE "comment" ADD FOREIGN KEY ("user_username") REFERENCES "users" ("username");
ALTER TABLE "comment" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "comment" ADD FOREIGN KEY ("thread_id") REFERENCES "threads" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "threads" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "threads" ADD FOREIGN KEY ("communities_name") REFERENCES "communities" ("communities_name");

ALTER TABLE "communities" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

-- ALTER TABLE "search_bar" ADD FOREIGN KEY ("full_name") REFERENCES "users" ("full_name");
--
-- ALTER TABLE "search_bar" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
--
-- ALTER TABLE "search_bar" ADD FOREIGN KEY ("communities") REFERENCES "communities" ("communities_name")