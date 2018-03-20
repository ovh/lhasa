package service

import (
	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/v1/models"
	v1repo "github.com/ovh/lhasa/api/v1/repositories"
)

// Deployer deploys an application version to the given environment and removes old deployments
type Deployer func(models.Application, models.Environment, *models.Deployment) error

// ApplicationDeployer deploy an application version to the given environment and removes old deployments
func ApplicationDeployer(db *gorm.DB, depRepo *v1repo.DeploymentRepository) Deployer {
	return func(app models.Application, env models.Environment, dep *models.Deployment) error {
		if err := depRepo.UndeployByApplicationEnv(app.Domain, app.Name, env.ID); err != nil {
			return err
		}

		dep.EnvironmentID = env.ID
		dep.ApplicationID = app.ID
		return depRepo.Save(dep)
	}
}
