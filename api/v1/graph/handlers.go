package graph

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/graphapi"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/deployment"
)

// HandlerGraph returns a dependency graph for a given deployment
func HandlerGraph(repo *Repository, depRepo *deployment.Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, deployment *v1.Deployment) (*graphapi.Graph, error) {
		entity, err := depRepo.FindOneBy(map[string]interface{}{"public_id": deployment.PublicID})
		if err != nil {
			return nil, err
		}
		deployment, ok := entity.(*v1.Deployment)
		if !ok {
			return nil, errors.New("internal type error")
		}

		return repo.FindByDeployment(deployment.PublicID)
	}, http.StatusOK)
}
