-- +migrate Down

DROP INDEX IF EXISTS idx_deployments_application_id;
DROP INDEX IF EXISTS idx_deployments_environment_id;

-- +migrate Up

CREATE INDEX idx_deployments_application_id ON deployments (application_id);
CREATE INDEX idx_deployments_environment_id ON deployments (environment_id);

