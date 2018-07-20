package deployment

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/v1"
	"github.com/sirupsen/logrus"
)

// Deployer deploys an application version to the given environment and removes old deployments
type Deployer func(application v1.Release, environment v1.Environment, deployment *v1.Deployment, log logrus.FieldLogger) (*v1.Deployment, bool, error)

// ApplicationDeployer deploys an application version to the given environment and removes old deployments
func ApplicationDeployer(tm db.TransactionManager, depFactory RepositoryFactory) Deployer {
	return func(app v1.Release, env v1.Environment, dep *v1.Deployment, log logrus.FieldLogger) (*v1.Deployment, bool, error) {
		c := false
		created := &c
		var d **v1.Deployment
		err := tm.Transaction(func(db *gorm.DB) error {
			depRepo := depFactory(db)
			j, err := getProperties(dep, app.Domain, app.Name, app.Version, env.Slug)
			dep.Properties.UnmarshalJSON(j)
			dep.EnvironmentID = env.ID
			dep.ApplicationID = app.ID

			// Looking for a previous deployment on the same release / environment
			prevs, err := depRepo.FindActivesByRelease(app.Domain, app.Name, app.Version, map[string]interface{}{"environment_id": env.ID})
			if err != nil {
				log.
					WithField("domain", app.Domain).
					WithField("name", app.Name).
					WithField("version", app.Version).
					WithField("env", env.Slug).
					WithError(err).
					Warnf("was not able to find previous active deployment for this release")
			}
			if len(prevs) > 0 {
				prev := prevs[0]
				prev.Properties.UnmarshalJSON(j)
				d = &prev
				return depRepo.Save(prev)
			}

			if err := depRepo.UndeployByApplicationEnv(app.Domain, app.Name, env.ID); err != nil {
				return err
			}

			c := true
			created = &c
			d = &dep
			return depRepo.Save(dep)
		}, log)
		return *d, *created, err
	}
}

func getProperties(dep *v1.Deployment, domain, name, version, envSlug string) ([]byte, error) {
	props := map[string]interface{}{}
	if len(dep.Properties.RawMessage) > 0 {
		if err := json.Unmarshal(dep.Properties.RawMessage, &props); err != nil {
			return nil, errors.BadRequestf("properties field should be a valid json object: %s", err.Error())
		}
	}
	props["_app_domain"] = domain
	props["_app_name"] = name
	props["_app_version"] = version
	props["_env_slug"] = envSlug
	return json.Marshal(props)
}
