-- +migrate Down

-- Remove indexes
DROP INDEX IF EXISTS "idx_applications_domain_name_version";
DROP INDEX IF EXISTS "idx_applications_tags";

-- Remove tables
DROP TABLE IF EXISTS "dependencies";
DROP TABLE IF EXISTS "applications";

-- +migrate Up

-- Tables
CREATE TABLE "applications" (
  "id"         BIGSERIAL,
  "domain"     VARCHAR(255) NOT NULL  DEFAULT '',
  "name"       VARCHAR(255) NOT NULL  DEFAULT '',
  "version"    VARCHAR(255) NOT NULL  DEFAULT '',
  "tags"       VARCHAR(255) [],
  "manifest"   JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_applications_domain_name_version
  ON "applications" ("domain", "name", "version");

CREATE INDEX idx_applications_tags
  ON "applications" USING GIN ("tags");

CREATE TABLE "dependencies" (
  "id"        BIGSERIAL,
  "owner_id"  BIGSERIAL,
  "target_id" BIGSERIAL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("owner_id") REFERENCES "applications" ON DELETE CASCADE,
  FOREIGN KEY ("target_id") REFERENCES "applications" ON DELETE CASCADE
);
