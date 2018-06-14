-- +migrate Down

ALTER TABLE releases DROP COLUMN IF EXISTS "properties";

-- +migrate Up

ALTER TABLE releases ADD COLUMN "properties" JSONB;
