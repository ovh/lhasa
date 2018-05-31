package handlers_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/ovh/lhasa/api/tests"
)

func TestAPIRoot(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	jsonObj := e.GET("/api").
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	jsonObj.Keys().ContainsOnly("_links")
	linksArray := jsonObj.Value("_links").Array()

	version1 := linksArray.Element(0).Object()
	version1.Value("href").String().Equal("/api/v1")
	version1.Value("rel").String().Equal("v1")

	unsecured := linksArray.Element(1).Object()
	unsecured.Value("href").String().Equal("/api/unsecured")
	unsecured.Value("rel").String().Equal("unsecured")
}

func TestPing(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/api/unsecured/mon").
		Expect().
		Status(http.StatusOK).
		Body().Equal("\"OK\"")
}

func TestVersion(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/api/unsecured/version").
		Expect().
		Status(http.StatusOK).
		Body().Equal("\"1.0.0\"")
}

func TestOpenAPI(t *testing.T) {
	server := tests.StartTestHTTPServer()
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/api/unsecured/openapi.json").
		Expect().
		Status(http.StatusOK).JSON().Object()
	//TODO check some stuff in the openapi json

	e.GET("/api/unsecured/openapi.yaml").
		Expect().
		Status(http.StatusOK)
}
