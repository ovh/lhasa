package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/routers"
	"github.com/ovh/lhasa/api/v1/repositories"
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
	flagDBVaultAlias   = application.Flag("db-vault-alias", "Set vault alias to use").Default("appcatalog-db").Envar("APPCATALOG_DB_VAULT_ALIAS").String()

	cmdVersion = application.Command(cmdCodeVersion, "Shows version number.")

	cmdMigrate     = application.Command(cmdCodeMigrate, "Only run migrations and return (not for production use).")
	cmdMigrateUp   = cmdMigrate.Command(cmdCodeMigrateUp, "Runs migrations upward (default).").Default()
	cmdMigrateDown = cmdMigrate.Command(cmdCodeMigrateDown, "Runs migrations downward.")

	cmdStart      = application.Command(cmdCodeStart, "Starts application.").Default()
	flagStartPort = cmdStart.Flag("port", "Listening port for the application.").Short('p').Envar("APPCATALOG_PORT").Default("8081").Uint()
)

func main() {
	command := kingpin.MustParse(application.Parse(os.Args[1:]))
	log := configureLogger(*flagVerbose, *flagDebug, *flagQuiet, *flagJSONOutput)

	switch command {
	case cmdCodeVersion:
		fmt.Println(version)
		return
	case cmdCodeMigrate:
	case fmt.Sprintf("%s %s", cmdCodeMigrate, cmdCodeMigrateUp):
		db := waitForDB(log)
		runMigrationsUp(db.DB(), log)
		return
	case fmt.Sprintf("%s %s", cmdCodeMigrate, cmdCodeMigrateDown):
		db := waitForDB(log)
		runMigrationsDown(db.DB(), log)
		return
	case cmdCodeStart:
		db := waitForDB(log)
		applicationRepository := repositories.NewApplicationRepository(db)
		if *flagAutoMigrations {
			runMigrationsUp(db.DB(), log)
		}
		router := routers.NewRouter(applicationRepository, db.DB(), version, *flagDebug, log)
		panic(router.Run(fmt.Sprintf(":%d", *flagStartPort)))
	}
}

func waitForDB(log *logrus.Logger) *gorm.DB {
	dbHandle, err := db.NewFromVault(*flagDBVaultAlias, *flagDebug, log)
	for err != nil {
		log.WithError(err).Errorf("cannot get DB handle, retrying in %d seconds", dbRetryDuration)
		time.Sleep(dbRetryDuration * time.Second)
		dbHandle, err = db.NewFromVault(*flagDBVaultAlias, *flagDebug, log)
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

func configureLogger(isVerbose, isDebug, isQuiet, isJSON bool) *logrus.Logger {
	log := logrus.New()
	log.Level = logrus.WarnLevel
	if isVerbose {
		log.Level = logrus.InfoLevel
	}
	if isQuiet {
		log.Level = logrus.FatalLevel
	}
	if isDebug {
		log.Level = logrus.DebugLevel
	}
	if isJSON {
		log.Formatter = &logrus.JSONFormatter{}
	}
	return log
}
