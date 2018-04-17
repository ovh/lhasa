package deployment

import (
	"encoding/json"

	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

// Depend deploys an application version to the given environment and removes old deployments
type Depend func(src *v1.Deployment, target *v1.Deployment, t string) error

// Dependency declare an observable dependency on one application
func Dependency(depRepo *Repository) Depend {
	return func(src *v1.Deployment, target *v1.Deployment, dependencyType string) error {
		dependencies := make([]v1.DeploymentDependency, 0)
		// Check current value
		if src.Dependencies.RawMessage != nil {
			json.Unmarshal(src.Dependencies.RawMessage, &dependencies)
		}
		// Check if link already exist
		var alreadyStored = false
		for _, deps := range dependencies {
			if deps.TargetID == target.PublicID {
				alreadyStored = true
			}
		}
		// Check if already stored
		if alreadyStored {
			// No update is needed
			return nil
		}
		dependencies = append(dependencies, v1.DeploymentDependency{TargetID: target.PublicID, Type: dependencyType})
		var marshalErr error
		src.Dependencies.RawMessage, marshalErr = json.Marshal(dependencies)
		if marshalErr != nil {
			return &(hateoas.InternalError{Message: marshalErr.Error(), Detail: src.PublicID})
		}
		return depRepo.Save(src)
	}
}
