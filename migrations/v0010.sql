-- +migrate Down

DROP INDEX IF EXISTS "idx_deployments_public_id";

-- +migrate Up

CREATE UNIQUE INDEX "idx_deployments_public_id"
  ON "deployments" ("public_id");
