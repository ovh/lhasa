package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/ovh/lhasa/api/config"
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/logger"
	"github.com/ovh/lhasa/api/routing"
)

var (
	// version number of app catalog
	version = "0.0.0"
)

const (
	cmdCodeVersion     = "version"
	cmdCodeMigrate     = "migrate"
	cmdCodeMigrateUp   = "up"
	cmdCodeMigrateDown = "down"
	cmdCodeStart       = "start"

	dbRetryDuration = 5
)

var (
	application        = kingpin.New("appcatalog", "Application catalog.")
	flagAutoMigrations = application.Flag("auto-migrate", "Enables auto migrations (not for production use).").Envar("APPCATALOG_AUTO_MIGRATE").Bool()
	flagVerbose        = application.Flag("verbose", "Enables verbose mode.").Short('v').Envar("APPCATALOG_VERBOSE_MODE").Bool()
	flagDebug          = application.Flag("debug", "Enables debug mode (routing and sql logging).").Envar("APPCATALOG_DEBUG_MODE").Bool()
	flagQuiet          = application.Flag("quiet", "Enables quiet mode.").Short('q').Envar("APPCATALOG_QUIET_MODE").Bool()
	flagJSONOutput     = application.Flag("json", "Enables JSON output.").Envar("APPCATALOG_JSON_OUTPUT").Bool()
	flagConfigFile     = application.Flag("config", "Json configuration file").Default(".config.json").Envar("APPCATALOG_CONFIG_FILE").File()
	flagDBAlias        = application.Flag("db-alias", "Set alias to use in json configuration").Default("appcatalog-db").Envar("APPCATALOG_DB_ALIAS").String()

	cmdVersion = application.Command(cmdCodeVersion, "Shows version number.")

	cmdMigrate     = application.Command(cmdCodeMigrate, "Only run migrations and return (not for production use).")
	cmdMigrateUp   = cmdMigrate.Command(cmdCodeMigrateUp, "Runs migrations upward (default).").Default()
	cmdMigrateDown = cmdMigrate.Command(cmdCodeMigrateDown, "Runs migrations downward.")

	cmdStart            = application.Command(cmdCodeStart, "Starts application.").Default()
	flagStartPort       = cmdStart.Flag("port", "Listening port for the application.").Short('p').Envar("APPCATALOG_PORT").Default("8081").Uint()
	flagHateoasBasePath = cmdStart.Flag("hateoas-base-path", "Base path to use for Hateoas links").Envar("APPCATALOG_HATEOAS_BASE_PATH").Default("/api").String()
	flagUIBasePath      = cmdStart.Flag("ui-base-path", "Base path to use for UI redirections").Envar("APPCATALOG_UI_BASE_PATH").Default("/ui").String()
)

func main() {
	log := logger.NewLogger(*flagVerbose, *flagDebug, *flagQuiet, *flagJSONOutput)

	command, err := application.Parse(os.Args[1:])
	if err != nil {
		log.WithError(err).Fatal("cannot start appcatalog")
	}

	switch command {
	case cmdCodeVersion:
		fmt.Println(version)
		return
	case cmdCodeMigrate:
	case fmt.Sprintf("%s %s", cmdCodeMigrate, cmdCodeMigrateUp):
		parseConf(log)
		db := waitForDB(log)
		runMigrationsUp(db.DB(), log)
		return
	case fmt.Sprintf("%s %s", cmdCodeMigrate, cmdCodeMigrateDown):
		parseConf(log)
		db := waitForDB(log)
		runMigrationsDown(db.DB(), log)
		return
	case cmdCodeStart:
		parseConf(log)
		tm := db.NewTransactionManager(waitForDB(log), log)
		if *flagAutoMigrations {
			runMigrationsUp(tm.DB().DB(), log)
		}
		router := routers.NewRouter(tm, version, *flagHateoasBasePath, *flagUIBasePath, *flagDebug, log)
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", *flagStartPort),
			Handler: router,
		}
		panic(srv.ListenAndServe())
	}
}

func parseConf(log *logrus.Logger) {
	if err := config.LoadFromFile(*flagConfigFile); err != nil {
		log.WithError(err).Fatalf("cannot read configuration file")
	}
}

func waitForDB(log *logrus.Logger) *gorm.DB {
	dbHandle, err := config.NewDBHandle(*flagConfigFile, *flagDBAlias, *flagDebug, log)
	for err != nil {
		log.WithError(err).Errorf("cannot get DB handle, retrying in %d seconds", dbRetryDuration)
		time.Sleep(dbRetryDuration * time.Second)
		dbHandle, err = config.NewDBHandle(*flagConfigFile, *flagDBAlias, *flagDebug, log)
	}
	return dbHandle
}

func runMigrationsUp(datasource *sql.DB, log *logrus.Logger) {
	if err := db.MigrateUp(datasource, log); err != nil {
		log.WithError(err).Fatalf("cannot run migrations")
	}
}

func runMigrationsDown(datasource *sql.DB, log *logrus.Logger) {
	if err := db.MigrateDown(datasource, log); err != nil {
		log.WithError(err).Fatalf("cannot run migrations")
	}
}
