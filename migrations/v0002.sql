-- +migrate Down

DROP TABLE IF EXISTS "deployments";
DROP TABLE IF EXISTS "environments";

-- +migrate Up

CREATE TABLE "environments" (
  "id"         BIGSERIAL,
  "slug"       VARCHAR(255) NOT NULL DEFAULT '',
  "name"       VARCHAR(255),
  "properties" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY ("id")
);

CREATE TABLE "deployments" (
  "id"             BIGSERIAL,
  "public_id"      VARCHAR(255),
  "environment_id" BIGINT NOT NULL  DEFAULT 0,
  "application_id" BIGINT NOT NULL  DEFAULT 0,
  "properties"     JSONB,
  "created_at"     TIMESTAMP WITH TIME ZONE,
  "undeployed_at"  TIMESTAMP WITH TIME ZONE,
  "updated_at"     TIMESTAMP WITH TIME ZONE,
  "deleted_at"     TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("environment_id") REFERENCES "environments" ON DELETE CASCADE,
  FOREIGN KEY ("application_id") REFERENCES "applications" ON DELETE CASCADE
);
