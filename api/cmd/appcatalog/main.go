package main

import (
	"fmt"
	"os"

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
	cmdCodeVersion = "version"
	cmdCodeMigrate = "migrate"
	cmdCodeStart   = "start"
)

var (
	application        = kingpin.New("appcatalog", "Application catalog.")
	flagAutoMigrations = application.Flag("auto-migrate", "Enables auto migrations (not for production use).").Envar("APPCATALOG_AUTO_MIGRATE").Bool()
	flagVerbose        = application.Flag("verbose", "Enables verbose mode.").Short('v').Envar("APPCATALOG_VERBOSE_MODE").Bool()
	flagDebug          = application.Flag("debug", "Enables debug mode (routing and sql logging).").Envar("APPCATALOG_DEBUG_MODE").Bool()
	flagQuiet          = application.Flag("quiet", "Enables quiet mode.").Short('q').Envar("APPCATALOG_QUIET_MODE").Bool()
	flagJSONOutput     = application.Flag("json", "Enables JSON output.").Envar("APPCATALOG_JSON_OUTPUT").Bool()
	flagDBVaultAlias   = cmdStart.Flag("db-vault-alias", "Set vault alias to use").Default("appcatalog-db").Envar("APPCATALOG_DB_VAULT_ALIAS").String()

	cmdVersion = application.Command(cmdCodeVersion, "Shows version number.")

	cmdMigrate = application.Command(cmdCodeMigrate, "Only run migrations and return (not for production use).")

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
		db := db.NewFromVault(log, *flagDebug, *flagDBVaultAlias)
		applicationRepository := repositories.NewApplicationRepository(db)
		runMigrations(applicationRepository, log)
		return
	case cmdCodeStart:
		db := db.NewFromVault(log, *flagDebug, *flagDBVaultAlias)
		applicationRepository := repositories.NewApplicationRepository(db)
		if *flagAutoMigrations {
			runMigrations(applicationRepository, log)
		}
		router := routers.NewRouter(applicationRepository, version, *flagDebug, log)
		panic(router.Run(fmt.Sprintf(":%d", *flagStartPort)))
	}
}

func runMigrations(applicationRepository *repositories.ApplicationRepository, log *logrus.Logger) {
	if err := applicationRepository.Migrate(); err != nil {
		log.WithError(err).Fatalf("cannot run migrations")
	}
}

func configureLogger(isVerbose, isDebug, isQuiet, isJSON bool) *logrus.Logger {
	log := logrus.New()
	log.Level = logrus.WarnLevel
	log.Out = os.Stderr
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
