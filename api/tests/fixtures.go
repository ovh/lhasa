package tests

import (
	"net/http/httptest"

	mocket "github.com/Selvatico/go-mocket"

	"github.com/jinzhu/gorm"

	"github.com/ovh/lhasa/api/logger"
	"github.com/ovh/lhasa/api/routing"
)

// StartTestHTTPServer starts a fake http server for testing purposes
func StartTestHTTPServer() *httptest.Server {
	log := logger.NewLogger(true, true, false, true)
	mocket.Catcher.Register()
	mocket.Catcher.Reset()
	mocket.Catcher.Logging = true
	db, _ := gorm.Open(mocket.DRIVER_NAME, "any_string")

	router := routers.NewRouter(db, "1.0.0", "/api", true, log)
	server := httptest.NewServer(router)
	return server
}
