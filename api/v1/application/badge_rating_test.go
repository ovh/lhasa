package application_test

import (
	"net/http"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/gavv/httpexpect"
	"github.com/ovh/lhasa/api/tests"
)

func mockSQLQuery(query string, jsonOutput []map[string]interface{}) *mocket.FakeResponse {
	return mocket.Catcher.
		NewMock().
		WithQuery(query).
		WithReply(jsonOutput)
}

func mockSelectBadges(jsonOutput []map[string]interface{}) *mocket.FakeResponse {
	return mockSQLQuery(`SELECT * FROM "badges"  WHERE "badges"."deleted_at" IS NULL`, jsonOutput)
}

func mockSelectApplications(jsonOutput []map[string]interface{}) *mocket.FakeResponse {
	return mockSQLQuery(`SELECT * FROM "releases"  WHERE`, jsonOutput)
}

func TestBadgeRatingList(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	mockSelectBadges(
		[]map[string]interface{}{
			{
				"id":    1,
				"slug":  "mybadge1",
				"title": "My Badge 1",
				"levels": []byte(`[
				{
					"id": "unset",
					"label": "Unknown", 
					"color": "gray",
					"isdefault": true
				}]`),
			},
			{
				"id":    2,
				"slug":  "mybadge2",
				"title": "My Badge 2",
				"levels": []byte(`[
				{
					"id": "unset",
					"label": "Unknown", 
					"color": "gray",
					"isdefault": true
				},
				{
					"id": "good",
					"label": "★★★★★", 
					"color": "green"
				}]`),
			},
		},
	)

	mockSelectApplications(
		[]map[string]interface{}{
			{
				"name":    "myapp",
				"domain":  "mydomain",
				"version": "1.2.3",
				"badge_ratings": []byte(`{
					"mybadge2": {
						"value": "good",
						"comment": "tests passed"
					}
				}`),
			},
		},
	)

	e := httpexpect.New(t, server.URL)
	e.GET("/api/v1/applications/mydomain/myapp/versions/1.0.0/badges").
		Expect().
		Status(http.StatusOK).
		JSON().Array().Equal([]map[string]interface{}{
		{
			"badgetitle": "My Badge 1",
			"value":      "unset",
			"comment":    "",
			"level": map[string]interface{}{
				"id":          "unset",
				"label":       "Unknown",
				"color":       "gray",
				"description": "",
				"isdefault":   true,
			},
			"_links": []interface{}{
				map[string]string{
					"href": "/api/v1/badges/mybadge1",
					"rel":  "badge",
				},
				map[string]string{
					"href": "/api/v1/applications/mydomain/myapp/versions/1.2.3",
					"rel":  "release",
				},
			},
		},
		{
			"badgetitle": "My Badge 2",
			"value":      "good",
			"comment":    "tests passed",
			"level": map[string]interface{}{
				"id":          "good",
				"label":       "★★★★★",
				"color":       "green",
				"description": "",
				"isdefault":   false,
			},
			"_links": []interface{}{
				map[string]string{
					"href": "/api/v1/badges/mybadge2",
					"rel":  "badge",
				},
				map[string]string{
					"href": "/api/v1/applications/mydomain/myapp/versions/1.2.3",
					"rel":  "release",
				},
			},
		},
	})

}

func TestBadgeRatingSetSuccess(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	mockSelectBadges(
		[]map[string]interface{}{
			{
				"id":    1,
				"slug":  "mybadge1",
				"title": "My Badge 1",
				"levels": []byte(`[
				{
					"id": "unset",
					"label": "Unknown", 
					"color": "gray",
					"isdefault": true
				}]`),
			},
			{
				"id":    2,
				"slug":  "mybadge2",
				"title": "My Badge 2",
				"levels": []byte(`[
				{
					"id": "unset",
					"label": "Unknown", 
					"color": "gray",
					"isdefault": true
				},
				{
					"id": "good",
					"label": "★★★★★", 
					"color": "green"
				}]`),
			},
		},
	)

	mockSelectApplications(
		[]map[string]interface{}{
			{
				"id": 1,
				"badge_ratings": []byte(`{
					"mybadge2": {
						"value": "good",
						"comment": "tests passed"
					}
				}`),
			},
		},
	)

	m := mocket.Catcher.NewMock().WithQuery(`UPDATE "releases" SET "domain" = ?, "name" = ?, "version" = ?, "properties" = ?, "manifest" = ?, "tags" = ?, "created_at" = ?, "updated_at" = ?, "deleted_at" = ?, "badge_ratings" = ?  WHERE "releases"."id" = ?`)

	bdgRating := `
	{
		"level":   "good",
		"comment": "This is a useless comment"
	}`
	e := httpexpect.New(t, server.URL)
	e.PUT("/api/v1/applications/mydomain/myapp/versions/1.0.0/badgeratings/mybadge2").
		WithHeader("Content-Type", "application/json").
		WithText(bdgRating).
		Expect().
		Status(http.StatusCreated)

	if m.Triggered == false {
		t.Fatal("expected the database to trigger an `UPDATE releases` statement but it has not been triggered")
	}
}

func TestBadgeRatingSetLevelNotFound(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	mockSelectBadges(
		[]map[string]interface{}{
			{
				"id":    1,
				"slug":  "mybadge1",
				"title": "My Badge 1",
				"levels": []byte(`[
				{
					"id": "unset",
					"label": "Unknown", 
					"color": "gray",
					"isdefault": true
				}]`),
			},
		},
	)

	mockSelectApplications(
		[]map[string]interface{}{
			{
				"badge_ratings": []byte(`{
					"mybadge2": {
						"value": "good",
						"comment": "tests passed"
					}
				}`),
			},
		},
	)

	bdgRating := `
	{
		"level":   "leveldoesnotexist",
		"comment": "This is a useless comment"
	}`

	e := httpexpect.New(t, server.URL)
	e.PUT("/api/v1/applications/mydomain/myapp/versions/1.0.0/badgeratings/mybadge2").
		WithHeader("Content-Type", "application/json").
		WithText(bdgRating).
		Expect().
		Status(http.StatusNotFound).
		JSON().
		Object().
		ValueEqual("error", "there is no level id `leveldoesnotexist` not found")
}

func TestBadgeRatingSetBadgeNotFound(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	mockSelectBadges([]map[string]interface{}{})

	mockSelectApplications(
		[]map[string]interface{}{
			{
				"badge_ratings": []byte(`{
					"mybadge2": {
						"value": "good",
						"comment": "tests passed"
					}
				}`),
			},
		},
	)

	bdgRating := `
	{
		"level":   "compliant",
		"comment": "This is a useless comment"
	}`

	e := httpexpect.New(t, server.URL)
	e.PUT("/api/v1/applications/mydomain/myapp/versions/1.0.0/badgeratings/mybadge2").
		WithHeader("Content-Type", "application/json").
		WithText(bdgRating).
		Expect().
		Status(http.StatusNotFound).
		JSON().
		Object().
		ValueEqual("error", "entity v1.Badge[slug=mybadge2] does not exist: entity v1.Badge[slug=mybadge2] does not exist")
}

func TestBadgeRatingSetAppNotFound(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	mocket.Catcher.
		NewMock().
		WithQuery(`SELECT * FROM "applications"  WHERE "applications"."deleted_at" IS NULL`).
		WithArgs(int64(27)).
		WithQueryException()

	bdgRating := `
	{
		"level":   "compliant",
		"comment": "This is a useless comment"
	}`

	e := httpexpect.New(t, server.URL)
	e.PUT("/api/v1/applications/mydomain/myapp/versions/1.0.0/badgeratings/mybadge2").
		WithHeader("Content-Type", "application/json").
		WithText(bdgRating).
		Expect().
		Status(http.StatusNotFound).
		JSON().
		Object().
		ValueEqual("error", "entity v1.Release[domain=mydomain name=myapp version=1.0.0] does not exist: entity v1.Release[domain=mydomain name=myapp version=1.0.0] does not exist")
}

func TestBadgeRatingRequestMalformedJSON(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.PUT("/api/v1/applications/mydomain/myapp/versions/1.0.0/badgeratings/mybadge2").
		WithHeader("Content-Type", "application/json").
		WithText("{THIS IS A MALFORMED JSON}").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().
		ValueEqual("error", "binding error: error parsing request body")
}

func TestBadgeRatingDeleteSuccess(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	mockSelectBadges(
		[]map[string]interface{}{
			{
				"id":    1,
				"slug":  "mybadge2",
				"title": "My Badge 2",
				"levels": []byte(`[
				{
					"id", "unset",
					"label": "Unknown", 
					"color": "gray",
					"isdefault": true
				}]`),
			},
		},
	)

	mockSelectApplications(
		[]map[string]interface{}{
			{
				"badge_ratings": []byte(`{
					"mybadge2": {
						"value": "thisisadeprecatedvalue",
						"comment": ""
					}
				}`),
			},
		},
	)

	m := mocket.Catcher.NewMock().WithQuery(`INSERT INTO "releases"`)

	e := httpexpect.New(t, server.URL)
	e.DELETE("/api/v1/applications/mydomain/myapp/versions/1.0.0/badgeratings/mybadge2").
		Expect().
		Status(http.StatusOK)
	if m.Triggered == false {
		t.Fatal("expected the database to trigger an `INSERT INTO releases` statement but it has not been triggered")
	}
}
