CREATE TABLE "message" (
  "id" bigserial PRIMARY KEY,
  "content" text NOT NULL,
  "is_seen" boolean NOT NULL,
  "sender_username" varchar NOT NULL,
  "receiver_username" varchar NOT NULL,
  "sent_at" timestamp NOT NULL DEFAULT NOW(),
  "media_url" varchar NOT NULL,
  "media_type" varchar NOT NULL
);

ALTER TABLE "message" ADD FOREIGN KEY ("sender_username") REFERENCES "users" ("username");

ALTER TABLE "message" ADD FOREIGN KEY ("receiver_username") REFERENCES "users" ("username");



