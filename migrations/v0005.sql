-- +migrate Down

-- Remove indexes
DROP INDEX IF EXISTS "idx_contents_name";

-- Remove tables
DROP TABLE IF EXISTS "contents";

-- +migrate Up

-- Tables
CREATE TABLE "contents" (
  "id"             BIGSERIAL,
  "name"           VARCHAR(255) NOT NULL  DEFAULT '',
  "content_type"   VARCHAR(128) NOT NULL  DEFAULT 'text/plain',
  "locale"         VARCHAR(255) NOT NULL  DEFAULT 'en-GB',
  "body"           BYTEA NOT NULL  DEFAULT '',
  "created_at"     TIMESTAMP WITH TIME ZONE,
  "updated_at"     TIMESTAMP WITH TIME ZONE,
  "deleted_at"     TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_contents_name
  ON "contents" ("name","locale");
