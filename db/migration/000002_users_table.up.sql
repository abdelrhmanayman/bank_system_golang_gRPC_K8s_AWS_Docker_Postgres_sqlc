CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_pwd" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz DEFAULT ('0001-01-01 00:00:00Z'),
  "created_at" timestamptz DEFAULT (now())
);

ALTER TABLE "account" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "account" ADD CONSTRAINT "owner_currency_acc" UNIQUE ("owner", "currency");
