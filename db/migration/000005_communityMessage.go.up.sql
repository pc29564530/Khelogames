CREATE TABLE "communitymessage" (
  "id" BIGSERIAL PRIMARY KEY,
  "community_name" varchar NOT NULL,
  "sender_username" varchar NOT NULL,
  "content" text NOT NULL,
  "sent_at" timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE "uploadmedia" (
    "id" BIGSERIAL PRIMARY KEY,
    "media_url" varchar NOT NULL,
    "media_type" varchar NOT NULL,
    "sent_at" timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE "messagemedia" (
    "message_id" bigint NOT NULL,
    "media_id" bigint NOT NULL,
    PRIMARY KEY (message_id, media_id)
);

ALTER TABLE "messagemedia" ADD FOREIGN KEY ("message_id") REFERENCES "communitymessage" ("id");

ALTER TABLE "messagemedia" ADD FOREIGN KEY ("media_id") REFERENCES "uploadmedia" ("id");

ALTER TABLE "communitymessage" ADD FOREIGN KEY ("community_name") REFERENCES "communities" ("communities_name");

ALTER TABLE "communitymessage" ADD FOREIGN KEY ("sender_username") REFERENCES "users" ("username");

