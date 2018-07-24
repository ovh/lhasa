package tests

import (
	"net/http/httptest"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/gobwas/glob"
	"github.com/ovh/lhasa/api/config"
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/security"

	"github.com/jinzhu/gorm"

	"github.com/ovh/lhasa/api/logger"
	"github.com/ovh/lhasa/api/routing"
)

// StartTestHTTPServer starts a fake http server for testing purposes
func StartTestHTTPServer() *httptest.Server {
	log := logger.NewLogger(true, true, false, false)
	mocket.Catcher.Register()
	mocket.Catcher.Reset()
	mocket.Catcher.Logging = true
	dbHandle, _ := gorm.Open(mocket.DRIVER_NAME, "any_string")
	tm := db.NewTransactionManager(dbHandle)
	c := config.Lhasa{Policy: security.Policy{"ROLE_ADMIN": {"X-Remote-User": {glob.MustCompile("*")}}}}
	router := routers.NewRouter(tm, c, "1.0.0", "/api", "/ui", "/", "./", true, log)
	server := httptest.NewServer(router)
	return server
}
