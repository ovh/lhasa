package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ovh/lhasa/api/models"
)

// listResponseLinks is a subobject served in a LIST reponse to implement HATEOAS
type listResponseLinks struct {
	Base string `json:"base"` // base URL for the listing (e.g. /api/v0/applications)
	Next string `json:"next"` // URL to next page (e.g. /api/v0/applications/?start=30&size=10)
	Self string `json:"self"` // URL to the current sublist (e.g. /api/v0/applications/?start=20&size=10)
}

// listResponseLinksPagination is a subobject served in a LIST response to implement HATEOAS and pagination
type listResponseLinksPagination struct {
	Links listResponseLinks `json:"_links"`
	Start int               `json:"start"`
	Size  int               `json:"size"`
	Limit int               `json:"limit"`
}

// ListResponse is the response object served in response to a successful LIST request
type ListResponse struct {
	*listResponseLinksPagination
	Results []models.Application `json:"results"`
}

// listRequest gather requests parameters
type listRequest struct {
	Query string `form:"query"` // Search query to filter results
	Start int    `form:"start"` // pagination start (default = 0)
	Size  int    `form:"size"`  // page size  (default = 100)
}

// ListApplicationsHandler returns a filtered and paginated applications list
func ListApplicationsHandler(c *gin.Context) (*ListResponse, error) {
	var request listRequest
	if err := c.Bind(&request); err != nil {
		return nil, err
	}
	const LIMIT int = 2000
	if request.Size <= 0 {
		request.Size = 1000
	}
	if request.Size > LIMIT {
		request.Size = LIMIT
	}
	if request.Start < 0 {
		request.Start = 0
	}

	links := listResponseLinks{Base: "", Next: "", Self: ""}
	linkspagi := listResponseLinksPagination{Start: request.Start, Size: request.Size, Limit: LIMIT, Links: links}
	results, err := models.ListApplications(request.Query, request.Start, request.Size)
	if err != nil {
		return nil, err
	}
	response := ListResponse{listResponseLinksPagination: &linkspagi, Results: results}
	return &response, nil
}

// CreateApplicationHandler adds an application to the database
func CreateApplicationHandler(c *gin.Context) (*models.Application, error) {
	var app models.Application
	if err := c.BindJSON(&app); err != nil {
		return &app, err
	}
	if err := models.CreateApplication(&app); err != nil {
		return &app, err
	}
	return &app, nil
}

// DeleteApplicationHandler deletes an application from the database
func DeleteApplicationHandler(c *gin.Context) (string, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return "Could not parse ID", err
	}
	if _, err := models.DeleteApplication(int(id)); err != nil {
		return "An error occured while deleting from the database", err
	}
	return "application deleted", nil
}

// DetailApplicationHandler returns the details of an application
func DetailApplicationHandler(c *gin.Context) (*models.Application, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	intid := int(id)
	return models.DetailApplication(intid)
}
