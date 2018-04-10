package deployment

import (
	"github.com/ovh/lhasa/api/v1"
)

// Deployer deploys an application version to the given environment and removes old deployments
type Deployer func(application v1.Application, environment v1.Environment, deployment *v1.Deployment) error

// ApplicationDeployer deploys an application version to the given environment and removes old deployments
func ApplicationDeployer(depRepo *Repository) Deployer {
	return func(app v1.Application, env v1.Environment, dep *v1.Deployment) error {
		if err := depRepo.UndeployByApplicationEnv(app.Domain, app.Name, env.ID); err != nil {
			return err
		}

		dep.EnvironmentID = env.ID
		dep.ApplicationID = app.ID
		return depRepo.Save(dep)
	}
}
