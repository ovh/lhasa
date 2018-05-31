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
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "releases"  WHERE "releases"."deleted_at" IS NULL`).WithReply(countReply)

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

	mocket.Catcher.NewMock().WithQuery(`SELECT "av".* FROM "releases" as "av"
JOIN "applications" ON "applications"."latest_release_id" = "av"."id"
WHERE`).WithReply(listReply)
	mocket.Catcher.NewMock().WithQuery(`SELECT count(1) as nbApp FROM "releases" as "av"
JOIN "applications" ON "applications"."latest_release_id" = "av"."id"
WHERE "av"."deleted_at" IS NULL`).WithReply(countReply)

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
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "releases"`).WithReply(countReply)

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
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "releases"`).WithReply(countReply)

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "releases"  WHERE`).WithReply(listReply)

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
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "releases"  WHERE`).WithReply(listReply)

	e := httpexpect.New(t, server.URL)

	e.DELETE("/api/v1/applications/mydomain/myapp/versions/1.0.0").
		Expect().
		Status(http.StatusNoContent)
}
