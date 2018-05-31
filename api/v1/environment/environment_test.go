package environment_test

import (
	"net/http"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/gavv/httpexpect"
	"github.com/ovh/lhasa/api/tests"
)

func TestEnvironmentAdd(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	emptyListReply := []map[string]interface{}{}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "environments"  WHERE`).WithReply(emptyListReply)

	e := httpexpect.New(t, server.URL)
	env := map[string]interface{}{}

	e.PUT("/api/v1/environments/prod").
		WithHeader("Content-Type", "application/json").
		WithJSON(env).
		Expect().
		Status(http.StatusCreated)
}

func TestEnvironmentUpdate(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	emptyListReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "environments"  WHERE`).WithReply(emptyListReply)

	e := httpexpect.New(t, server.URL)
	env := map[string]interface{}{}

	e.PUT("/api/v1/environments/prod").
		WithHeader("Content-Type", "application/json").
		WithJSON(env).
		Expect().
		Status(http.StatusOK)
}

func TestEnvironmentDelete(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "environments"  WHERE`).WithReply(listReply)

	e := httpexpect.New(t, server.URL)

	e.DELETE("/api/v1/environments/prod").
		Expect().
		Status(http.StatusNoContent)
}

func TestEnvironmentList(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	countReply := []map[string]interface{}{{"count(*)": 2}}
	listReply := []map[string]interface{}{{"name": "prod"}, {"name": "staging"}}

	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "environments"  WHERE`).WithReply(listReply)
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "environments"  WHERE "environments"."deleted_at" IS NULL`).WithReply(countReply)

	e := httpexpect.New(t, server.URL)
	jsonObj := e.GET("/api/v1/environments/").
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
	content.Element(0).Object().ValueEqual("name", "prod")
	content.Element(1).Object().ValueEqual("name", "staging")
}
