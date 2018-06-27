package deployment

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/v1"
)

// Deployer deploys an application version to the given environment and removes old deployments
type Deployer func(application v1.Release, environment v1.Environment, deployment *v1.Deployment) error

// ApplicationDeployer deploys an application version to the given environment and removes old deployments
func ApplicationDeployer(tm db.TransactionManager, depFactory RepositoryFactory) Deployer {
	return func(app v1.Release, env v1.Environment, dep *v1.Deployment) error {
		return tm.Transaction(func(db *gorm.DB) error {
			depRepo := depFactory(db)

			if err := depRepo.UndeployByApplicationEnv(app.Domain, app.Name, env.ID); err != nil {
				return err
			}
			props := map[string]interface{}{}
			if len(dep.Properties.RawMessage) > 0 {
				if err := json.Unmarshal(dep.Properties.RawMessage, &props); err != nil {
					return errors.BadRequestf("properties field should be a valid json object: %s", err.Error())
				}
			}
			props["_app_domain"] = app.Domain
			props["_app_name"] = app.Name
			props["_app_version"] = app.Version
			props["_env_slug"] = env.Slug
			j, err := json.Marshal(props)
			if err != nil {
				return err
			}

			dep.Properties.UnmarshalJSON(j)
			dep.EnvironmentID = env.ID
			dep.ApplicationID = app.ID
			return depRepo.Save(dep)
		})
	}
}
