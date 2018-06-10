-- +migrate Down

-- Remove tables
DROP TABLE IF EXISTS "badges";

ALTER TABLE applications DROP COLUMN IF EXISTS "badge_ratings";

-- Remove indexes


-- +migrate Up

-- Tables
CREATE TABLE IF NOT EXISTS "badges" (
  "id"         BIGSERIAL,
  "slug"       VARCHAR(255) UNIQUE NOT NULL,
  "title"      VARCHAR(255),
  "type"       VARCHAR(20),
  "levels"     JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY ("id"),
  CONSTRAINT non_empty CHECK (slug <> '')
);

ALTER TABLE applications ADD COLUMN "badge_ratings" JSONB;

-- Indexes
