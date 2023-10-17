CREATE TABLE "join_community" (
     "id" bigserial PRIMARY KEY,
     "community_name" varchar NOT NULL,
     "username" varchar NOT NULL
);

ALTER TABLE "join_community" ADD FOREIGN KEY ("community_name") REFERENCES "communities" ("communities_name");

ALTER TABLE "join_community" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
