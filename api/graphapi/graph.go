package graphapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
)

// HandlerFindAll returns a resource list
func HandlerFindAll(repository Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (*Graph, error) {

		// params and query are user to filter resultset
		graphResult, err := repository.FindAll()
		if err != nil {
			return nil, err
		}

		return graphResult, nil
	}, http.StatusOK)
}

// HandlerFindAllActive returns a resource list
func HandlerFindAllActive(repository Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (*Graph, error) {

		// params and query are user to filter resultset
		graphResult, err := repository.FindAllActive()
		if err != nil {
			return nil, err
		}

		return graphResult, nil
	}, http.StatusOK)
}

// Convert to dependencies node behaviour
func Convert(repo Repository, entities []interface{}) []ImplementNode {
	var nodes = []ImplementNode{}
	for _, entity := range entities {
		mappers := map[string]ImplementNode{}
		// Resolve child dependencies
		repo.Resolve(entity, mappers)
		// Add this resolved node
		nodes = append(nodes, repo.Map(entity, mappers))
	}
	return nodes
}
