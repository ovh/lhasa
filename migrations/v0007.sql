-- +migrate Down

DROP INDEX IF EXISTS "idx_applications_domain_name";

DROP TABLE IF EXISTS "applications";

ALTER TABLE IF EXISTS "releases"
  RENAME TO "applications";

ALTER INDEX IF EXISTS "idx_releases_domain_name_version"
RENAME TO "idx_applications_domain_name_version";

ALTER INDEX IF EXISTS "idx_releases_tags"
RENAME TO "idx_applications_tags";

-- +migrate Up

ALTER TABLE IF EXISTS "applications"
  RENAME TO "releases";

ALTER INDEX IF EXISTS "idx_applications_domain_name_version"
RENAME TO "idx_releases_domain_name_version";

ALTER INDEX IF EXISTS "idx_applications_tags"
RENAME TO "idx_releases_tags";

CREATE TABLE IF NOT EXISTS "applications" (
  "id"                BIGSERIAL,
  "domain"            VARCHAR(255) NOT NULL DEFAULT '',
  "name"              VARCHAR(255) NOT NULL DEFAULT '',
  "latest_release_id" BIGINT       NULL,
  "created_at"        TIMESTAMP WITH TIME ZONE,
  "updated_at"        TIMESTAMP WITH TIME ZONE,
  "deleted_at"        TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY ("id"),
  FOREIGN KEY (latest_release_id) REFERENCES "releases" ON DELETE SET NULL
);

-- IF NOT EXISTS is not compatible with postgres 9.4
CREATE UNIQUE INDEX "idx_applications_domain_name"
  ON "applications" ("domain", "name");

-- data migrations

INSERT INTO "applications" (
  "domain",
  "name",
  latest_release_id,
  "created_at",
  "updated_at"
)
  SELECT
    "v"."domain",
    "v"."name",
    max("v"."id"),
    now(),
    now()
  FROM "releases" as "v"
  WHERE "v"."deleted_at" IS NULL
        AND NOT EXISTS(SELECT 1
                       FROM "applications" as "a"
                       WHERE ("a"."domain", "a"."name") = ("v"."domain", "v"."name"))
  GROUP BY "domain", "name";
