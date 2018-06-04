-- +migrate Down

ALTER TABLE "deployments" DROP COLUMN IF EXISTS "dependencies";

-- +migrate Up

ALTER TABLE "deployments" ADD "dependencies" JSONB;

CREATE INDEX idx_deployments_dependencies ON "deployments" USING GIN ("dependencies");

