

ALTER TABLE IF EXISTS "account" DROP CONSTRAINT  IF EXISTS "owner_currency_acc";

ALTER TABLE IF EXISTS "account" DROP CONSTRAINT  IF EXISTS "account_owner_fkey";

DROP TABLE IF EXISTS users;

