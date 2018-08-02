package graph

import (
	"reflect"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/ovh/lhasa/api/graphapi"
	"github.com/ovh/lhasa/api/v1"
)

const (
	queryOutbound = `SELECT d2.* FROM "deployments" AS "d1"
CROSS JOIN jsonb_array_elements("d1"."dependencies") as "items"
JOIN "deployments" AS "d2" ON items->>'target' = "d2"."public_id"
  AND "d2"."deleted_at" IS NULL
  AND "d2"."undeployed_at" IS NULL
WHERE "d1"."deleted_at" IS NULL
  AND (("d1"."undeployed_at" IS NULL AND "d1"."public_id" = ?))`
)

// Repository is a repository manager
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetType returns the entity type managed by this repository
func (repo *Repository) GetType() reflect.Type {
	return reflect.TypeOf(graphapi.Graph{})
}

// FindAll fetch a collection of nodes matching each criteria
// criteria are not used now
func (repo *Repository) FindAll() (*graphapi.Graph, error) {
	var entities []*v1.Deployment
	err := repo.db.Preload("Application").Preload("Environment").Find(&entities).Error
	if err == nil {
		interfaces := []interface{}{}
		for _, entity := range entities {
			interfaces = append(interfaces, entity)
		}
		graph, err := repo.transformToGraph(graphapi.Convert(repo, interfaces))
		return graph, err
	}
	return nil, err
}

// FindAllActive fetch a collection of nodes matching each criteria
func (repo *Repository) FindAllActive(criteria map[string]interface{}) (*graphapi.Graph, error) {
	var entities []*v1.Deployment
	err := repo.db.Where(criteria).Where("\"undeployed_at\" is null").Preload("Application").Preload("Environment").Find(&entities).Error
	if err != nil {
		return nil, err
	}
	var interfaces []interface{}
	for _, entity := range entities {
		interfaces = append(interfaces, entity)
	}
	return repo.transformToGraph(graphapi.Convert(repo, interfaces))
}

// FindByDeployment find all inbound and outbound dependencies for a deployment
func (repo *Repository) FindByDeployment(publicID string) (*graphapi.Graph, error) {
	var entities []*v1.Deployment
	var interfaces []interface{}
	// load the deployment itself
	err := repo.db.Where("\"undeployed_at\" IS NULL AND \"public_id\" = ?", publicID).Preload("Application").Preload("Environment").Find(&entities).Error
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		interfaces = append(interfaces, entity)
	}
	// load the inbound dependencies
	err = repo.db.Where("\"undeployed_at\" IS NULL AND \"dependencies\" @> '[{\"target\": \"?\"}]'", gorm.Expr(publicID)).Preload("Application").Preload("Environment").Find(&entities).Error
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		interfaces = append(interfaces, entity)
	}
	// load the outbound dependencies
	rows, err := repo.db.Preload("Application").Preload("Environment").Raw(queryOutbound, publicID).Rows()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		entity := &v1.Deployment{}
		if err := repo.db.ScanRows(rows, entity); err != nil {
			return nil, err
		}
		app := &v1.Release{}
		if err := repo.db.Find(app, entity.ApplicationID).Error; err != nil {
			return nil, err
		}
		env := &v1.Environment{}
		if err := repo.db.Find(env, entity.EnvironmentID).Error; err != nil {
			return nil, err
		}
		entity.Application = app
		entity.Environment = env
		interfaces = append(interfaces, entity)
	}

	return repo.transformToGraph(graphapi.Convert(repo, interfaces))
}

// Resolve resolve dependencies
func (repo *Repository) Resolve(entity interface{}, mappers map[string]graphapi.ImplementNode) {
	repo.Map(entity.(*v1.Deployment).Application, mappers)
	repo.Map(entity.(*v1.Deployment).Environment, mappers)
}

// Map api entity
func (repo *Repository) internalMap(entity interface{}, mappers map[string]graphapi.ImplementNode) (string, graphapi.ImplementNode) {
	deployment, isDeployment := entity.(*v1.Deployment)
	if isDeployment {
		graphEntity := v1.DeploymentGraph{
			PublicID:     deployment.PublicID,
			Application:  mappers["application"],
			Environment:  mappers["environment"],
			Dependencies: deployment.Dependencies.RawMessage,
			Properties:   deployment.Properties.RawMessage,
			UndeployedAt: deployment.UndeployedAt,
			CreatedAt:    deployment.CreatedAt,
			UpdatedAt:    deployment.UpdatedAt,
		}
		return "deployemnt", &graphEntity
	}
	application, isApplication := entity.(*v1.Release)
	if isApplication {
		if application == nil {
			graphEntity := v1.ApplicationVersionGraph{}
			return "application", &graphEntity
		}
		uid, _ := uuid.NewV4()
		graphEntity := v1.ApplicationVersionGraph{
			PublicID:  uid.String(),
			Domain:    application.Domain,
			Name:      application.Name,
			Version:   application.Version,
			Tags:      application.Tags,
			CreatedAt: application.CreatedAt,
			UpdatedAt: application.UpdatedAt,
		}
		return "application", &graphEntity
	}
	environment, isenvironment := entity.(*v1.Environment)
	if isenvironment {
		if environment == nil {
			graphEntity := v1.EnvironmentGraph{}
			return "environment", &graphEntity
		}
		uid, _ := uuid.NewV4()
		graphEntity := v1.EnvironmentGraph{
			PublicID:   uid.String(),
			Slug:       environment.Slug,
			Name:       environment.Name,
			Properties: environment.Properties.RawMessage,
			CreatedAt:  environment.CreatedAt,
			UpdatedAt:  environment.UpdatedAt,
		}
		return "environment", &graphEntity
	}
	return "", nil
}

// Map api entity
func (repo *Repository) Map(entity interface{}, mappers map[string]graphapi.ImplementNode) graphapi.ImplementNode {
	name, mapped := repo.internalMap(entity, mappers)
	mappers[name] = mapped
	return mapped
}

// transformToGraph transform to standard graph
func (repo *Repository) transformToGraph(entities []graphapi.ImplementNode) (*graphapi.Graph, error) {
	// Convert to graph representation
	mappedGraph := graphapi.Graph{
		Nodes: make([]graphapi.Node, 0),
		Edges: make([]graphapi.Edge, 0),
	}

	// Map on all single edge in graph
	edges := make(map[string]graphapi.Edge)
	types := make(map[string]bool)

	// Iterate on entities
	for _, entity := range entities {
		mappedGraph.Nodes = append(mappedGraph.Nodes, graphapi.Node{
			ID:         entity.GetID(),
			Name:       entity.GetName(),
			Type:       entity.GetType(),
			Properties: entity,
		})
		// Type
		types[entity.GetType()] = true
		// Iterate on all dependencies
		for _, edge := range entity.GetEdges() {
			var key = edge.GetFrom() + "->" + edge.GetTo()
			_, found := edges[key]
			if !found {
				edges[key] = graphapi.Edge{
					ID:         edge.GetID(),
					From:       edge.GetFrom(),
					To:         edge.GetTo(),
					Type:       edge.GetType(),
					Properties: edge.GetProperties(),
				}
			}
		}
	}

	// Iterate on index to dump edge
	for _, edge := range edges {
		mappedGraph.Edges = append(mappedGraph.Edges, edge)
	}

	return &mappedGraph, nil
}
