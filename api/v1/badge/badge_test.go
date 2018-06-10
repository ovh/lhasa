package badge_test

import (
	"net/http"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/gavv/httpexpect"
	"github.com/ovh/lhasa/api/tests"
)

func TestBadgeAdd(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	emptyListReply := []map[string]interface{}{}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "badges"  WHERE`).WithReply(emptyListReply)

	e := httpexpect.New(t, server.URL)

	bdg := `
		{
			"name":   "My Shiny Badge",
			"type":   "enum",
			"levels": [
				{"id": "unset", "label": "Unknown", "color": "lightgray", "isdefault": true},
				{"id": "error", "label": "%d errors detected", "color": "red"},
				{"id": "ok", "label": "clean", "color": "green"}
			]
		}
	`

	e.PUT("/api/v1/badges/myshinybadge").
		WithHeader("Content-Type", "application/json").
		WithText(bdg).
		Expect().
		Status(http.StatusCreated)
}

func TestBadgeUpdate(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	nonEmptyListReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "badges"  WHERE`).WithReply(nonEmptyListReply)

	e := httpexpect.New(t, server.URL)
	bdg := `
		{
			"title":   "My Shiny Badge",
			"type":   "enum",
			"levels": [
				{"id": "unset", "label": "Unknown", "color": "lightgray", "isdefault":true},
				{"id": "error", "label": "%d errors detected", "color": "red"},
				{"id": "ok", "label": "clean", "color": "green"}
			]
		}
	`
	e.PUT("/api/v1/badges/myshinybadge").
		WithHeader("Content-Type", "application/json").
		WithText(bdg).
		Expect().
		Status(http.StatusOK)
}

func TestBadgeDelete(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	listReply := []map[string]interface{}{{"id": 12}}
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "badges"  WHERE`).WithReply(listReply)

	e := httpexpect.New(t, server.URL)

	e.DELETE("/api/v1/badges/myshinybadge").
		Expect().
		Status(http.StatusNoContent)
}

func TestBadgeList(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	countReply := []map[string]interface{}{{"count(*)": 2}}
	listReply := []map[string]interface{}{{"title": "myshinybadge"}, {"title": "myshinybadge2"}}

	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "badges"  WHERE`).WithReply(listReply)
	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "badges"  WHERE "badges"."deleted_at" IS NULL`).WithReply(countReply)

	e := httpexpect.New(t, server.URL)
	jsonObj := e.GET("/api/v1/badges/").
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
	content.Element(0).Object().ValueEqual("title", "myshinybadge")
	content.Element(1).Object().ValueEqual("title", "myshinybadge2")
}
