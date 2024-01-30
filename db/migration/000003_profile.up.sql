CREATE TABLE "profile" (
        "id" bigserial PRIMARY KEY,
        "owner" varchar UNIQUE NOT NULL,
        "full_name" varchar NOT NULL,
        "bio" text NOT NULL,
        "avatar_url" varchar NOT NULL,
        "cover_url" varchar NOT NULL,
        "avatar_type" sql.NULLString,
        "created_at" timestamptz NOT NULL DEFAULT NOW()
        
);

ALTER TABLE "profile" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");