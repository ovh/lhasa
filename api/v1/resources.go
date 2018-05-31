package v1

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ovh/lhasa/api/hateoas"
)

// GetID returns the public ID of the entity
func (dom *Domain) GetID() string {
	return dom.Name
}

// SetID sets up the new ID of the entity
func (dom *Domain) SetID(id string) error {
	dom.Name = id
	return nil
}

// ToResource implements Resourceable
func (dom *Domain) ToResource(baseURL string) {
	dom.Resource.Links = []hateoas.ResourceLink{
		{Rel: "self", Href: dom.GetSelfURL(baseURL)},
		{Rel: "applications", Href: fmt.Sprintf("%s%s/%s", baseURL, ApplicationBasePath, dom.Name)},
	}
}

// GetSelfURL implements Resourceable
func (dom *Domain) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, DomainBasePath, dom.Name)
}

// GetID returns the public ID of the entity
func (cont *Content) GetID() string {
	return string(cont.ID)
}

// SetID sets up the new ID of the entity
func (cont *Content) SetID(id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	cont.ID = uint(idInt)
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (cont *Content) GetDeletedAt() *time.Time {
	return cont.DeletedAt
}

// ToResource implements Resourceable
func (cont *Content) ToResource(baseURL string) {
	cont.Resource.Links = []hateoas.ResourceLink{
		{Rel: "self", Href: cont.GetSelfURL(baseURL)},
	}
}

// GetSelfURL implements Resourceable
func (cont *Content) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, ApplicationBasePath, cont.Name)
}

// GetID returns the public ID of the entity
func (app *Application) GetID() string {
	return string(app.ID)
}

// SetID sets up the new ID of the entity
func (app *Application) SetID(id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	app.ID = uint(idInt)
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (app *Application) GetDeletedAt() *time.Time {
	return app.DeletedAt
}

// ToResource implements Resourceable
func (app *Application) ToResource(baseURL string) {
	app.Resource.Links = []hateoas.ResourceLink{
		{Rel: "self", Href: app.GetSelfURL(baseURL)},
		{Rel: "latest", Href: app.GetSelfURL(baseURL) + "/latest"},
		{Rel: "versions", Href: app.GetSelfURL(baseURL) + "/versions"}}
}

// GetSelfURL implements Resourceable
func (app *Application) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s/%s", baseURL, ApplicationBasePath, app.Domain, app.Name)
}

// GetID returns the public ID of the entity
func (app *Release) GetID() string {
	return string(app.ID)
}

// SetID sets up the new ID of the entity
func (app *Release) SetID(id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	app.ID = uint(idInt)
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (app *Release) GetDeletedAt() *time.Time {
	return app.DeletedAt
}

// ToResource implements Resourceable
func (app *Release) ToResource(baseURL string) {
	app.Resource.Links = []hateoas.ResourceLink{
		{Rel: "self", Href: app.GetSelfURL(baseURL)},
		{Rel: "deployments", Href: app.GetSelfURL(baseURL) + "/deployments"},
		{Rel: "deploy", Href: app.GetSelfURL(baseURL) + "/deploy/:environment"}}
}

// GetSelfURL implements Resourceable
func (app *Release) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s/%s/versions/%s", baseURL, ApplicationBasePath, app.Domain, app.Name, app.Version)
}

// GetID returns the public ID of the entity
func (dep *Deployment) GetID() string {
	return dep.PublicID
}

// SetID sets up the new ID of the entity
func (dep *Deployment) SetID(id string) error {
	dep.PublicID = id
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (dep *Deployment) GetDeletedAt() *time.Time {
	return dep.DeletedAt
}

// ToResource implements Resourceable
func (dep *Deployment) ToResource(baseURL string) {
	dep.Resource.Links = []hateoas.ResourceLink{
		{Rel: "self", Href: dep.GetSelfURL(baseURL)},
		{Rel: "add_link", Href: dep.GetSelfURL(baseURL) + "/add_link/:target_deployment_id"},
	}
	if dep.Environment != nil {
		dep.Resource.Links = append(dep.Resource.Links, hateoas.ResourceLink{Rel: "environment", Href: dep.Environment.GetSelfURL(baseURL)})
	}
	if dep.Application != nil {
		dep.Resource.Links = append(dep.Resource.Links, hateoas.ResourceLink{Rel: "application", Href: dep.Application.GetSelfURL(baseURL)})
	}
}

// GetSelfURL implements Resourceable
func (dep *Deployment) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, DeploymentBasePath, dep.PublicID)
}

// GetID returns the public ID of the entity
func (env *Environment) GetID() string {
	return env.Slug
}

// SetID sets up the new ID of the entity
func (env *Environment) SetID(id string) error {
	env.Slug = id
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (env *Environment) GetDeletedAt() *time.Time {
	return env.DeletedAt
}

// ToResource implements Resourceable
func (env *Environment) ToResource(baseURL string) {
	env.Resource.Links = []hateoas.ResourceLink{{Rel: "self", Href: env.GetSelfURL(baseURL)}}
}

// GetSelfURL implements Resourceable
func (env *Environment) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, EnvironmentBasePath, env.Slug)
}

// BADGES

// GetID returns the public ID of the entity
func (badge *Badge) GetID() string {
	return badge.Slug
}

// SetID sets up the new ID of the entity
func (badge *Badge) SetID(id string) error {
	badge.Slug = id
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (badge *Badge) GetDeletedAt() *time.Time {
	return badge.DeletedAt
}

// ToResource implements Resourceable
func (badge *Badge) ToResource(baseURL string) {
	badge.Resource.Links = []hateoas.ResourceLink{{Rel: "self", Href: badge.GetSelfURL(baseURL)}}
}

// GetSelfURL implements Resourceable
func (badge *Badge) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, BadgeBasePath, badge.Slug)
}
