package application_test

import (
	"net/http"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/gavv/httpexpect"
	"github.com/ovh/lhasa/api/tests"
)

func TestApplicationListEmpty(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	countReply := []map[string]interface{}{{"count(*)": 0}}
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "applications"  WHERE "applications"."deleted_at" IS NULL`).WithReply(countReply)

	e := httpexpect.New(t, server.URL)
	jsonObj := e.GET("/api/v1/applications/").
		Expect().
		Status(http.StatusPartialContent).
		JSON().Object()

	jsonObj.Keys().ContainsOnly("content", "pageMetadata", "_links")
	pageMetadata := jsonObj.Value("pageMetadata").Object()
	pageMetadata.Keys().ContainsOnly("totalElements", "totalPages", "size", "number")
	pageMetadata.ValueEqual("totalElements", 0).
		ValueEqual("totalPages", 0).
		// TODO: totalPages should be 1 for an empty list
		ValueEqual("size", 20).
		ValueEqual("number", 0)

	content := jsonObj.Value("content").Array()
	content.Empty()
}

func TestApplicationList(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	countReply := []map[string]interface{}{{"count(*)": 2}}
	listReply := []map[string]interface{}{{"name": "myapp1", "domain": "mydomain"}, {"name": "myapp2", "domain": "mydomain"}}

	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "applications"  WHERE`).WithReply(listReply)
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "applications"  WHERE "applications"."deleted_at" IS NULL`).WithReply(countReply)

	e := httpexpect.New(t, server.URL)
	jsonObj := e.GET("/api/v1/applications/").
		Expect().
		Status(http.StatusPartialContent).
		JSON().Object()

	jsonObj.Keys().ContainsOnly("content", "pageMetadata", "_links")
	pageMetadata := jsonObj.Value("pageMetadata").Object()
	pageMetadata.Keys().ContainsOnly("totalElements", "totalPages", "size", "number")
	pageMetadata.ValueEqual("totalElements", 2).
		ValueEqual("totalPages", 1).
		ValueEqual("size", 20).
		ValueEqual("number", 0)

	content := jsonObj.Value("content").Array()
	content.Length().Equal(2)
	content.Element(0).Object().ValueEqual("domain", "mydomain").ValueEqual("name", "myapp1")
	content.Element(1).Object().ValueEqual("domain", "mydomain").ValueEqual("name", "myapp2")

}

func TestApplicationAdd(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	countReply := []map[string]interface{}{{"count(*)": 0}}
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "applications"`).WithReply(countReply)

	e := httpexpect.New(t, server.URL)
	app := map[string]interface{}{
		"domain":   "mydomain",
		"name":     "myapp",
		"version":  "1.0.0",
		"manifest": "",
	}

	e.PUT("/api/v1/applications/mydomain/myapp/versions/1.0.0").
		WithHeader("Content-Type", "application/json").
		WithJSON(app).
		Expect().
		Status(http.StatusCreated)
}

func TestApplicationUpdate(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	countReply := []map[string]interface{}{{"count(*)": 1}}
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "applications"`).WithReply(countReply)

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "applications"  WHERE`).WithReply(listReply)

	e := httpexpect.New(t, server.URL)
	app := map[string]interface{}{
		"domain":   "mydomain",
		"name":     "myapp",
		"version":  "1.0.0",
		"manifest": "",
	}

	e.PUT("/api/v1/applications/mydomain/myapp/versions/1.0.0").
		WithHeader("Content-Type", "application/json").
		WithJSON(app).
		Expect().
		Status(http.StatusOK)
}

func TestApplicationDelete(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "applications"  WHERE`).WithReply(listReply)

	e := httpexpect.New(t, server.URL)

	e.DELETE("/api/v1/applications/mydomain/myapp/versions/1.0.0").
		Expect().
		Status(http.StatusNoContent)
}

func TestApplicationAssistantErrorApplicationDoesNotExist(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.POST("/api/v1/applications/mydomain/myapp/assistant").
		Expect().
		Status(http.StatusNotFound).JSON().Object().ValueEqual("error", "Application must exist not found")
	// TODO Fix the error message above
}

func TestApplicationAssistantErrorNoManifest(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "applications"  WHERE`).WithReply(listReply)

	app := map[string]interface{}{
		"domain":     "mydomain",
		"name":       "myname",
		"version":    "1.0.0",
		"repository": "bitbucket",
	}

	e := httpexpect.New(t, server.URL)
	e.POST("/api/v1/applications/mydomain/myapp/assistant").
		WithJSON(app).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("error", "Manifest cannot be null")
}

func TestApplicationAssistantErrorUserNotAuthenticated(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "applications"  WHERE`).WithReply(listReply)

	app := map[string]interface{}{
		"domain":     "mydomain",
		"name":       "myname",
		"version":    "1.0.0",
		"repository": "bitbucket",
		"manifest": map[string]interface{}{
			"name":    "myapp",
			"profile": "mydomain",
		},
	}

	e := httpexpect.New(t, server.URL)
	e.POST("/api/v1/applications/mydomain/myapp/assistant").
		WithJSON(app).
		Expect().
		Status(http.StatusUnauthorized).JSON().Object().
		ValueEqual("error", "User is not authorized")
}

func TestApplicationAssistant(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "applications"  WHERE`).WithReply(listReply)

	e := httpexpect.New(t, server.URL)

	app := map[string]interface{}{
		"domain":     "mydomain",
		"name":       "myname",
		"version":    "1.0.0",
		"repository": "bitbucket",
		"manifest": map[string]interface{}{
			"name":    "myapp",
			"profile": "mydomain",
		},
	}

	e.POST("/api/v1/applications/mydomain/myapp/assistant").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Remote-User", "john.doe").
		WithJSON(app).
		Expect().
		Status(http.StatusInternalServerError)
	//TODO: remove this test after the PR creation feature has been removed
}
