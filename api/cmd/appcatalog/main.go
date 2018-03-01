package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/routers"
	"github.com/ovh/lhasa/api/v0/models"
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
	flagDebug          = application.Flag("debug", "Enables debug mode.").Envar("APPCATALOG_DEBUG_MODE").Bool()
	flagPort           = application.Flag("port", "Listening port for the application.").Short('p').Envar("APPCATALOG_PORT").Default("8081").Uint()
	cmdVersion         = application.Command(cmdCodeVersion, "Shows version number.")
	cmdMigrate         = application.Command(cmdCodeMigrate, "Only run migrations and return (not for production use).")
	cmdStart           = application.Command(cmdCodeStart, "Starts application.").Hidden().Default()
	ginMode            = gin.ReleaseMode
)

func main() {
	command := kingpin.MustParse(application.Parse(os.Args[1:]))

	db := db.DB()
	db.LogMode(*flagDebug)
	if *flagDebug {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	applicationRepository := repositories.NewApplicationRepository(db)
	applicationVersionRepository := repositories.NewApplicationVersionAwareRepository(db)

	if *flagAutoMigrations || command == cmdCodeMigrate {
		if err := models.MigrateApplications(); err != nil {
			panic(err)
		}
		if err := applicationRepository.Migrate(); err != nil {
			panic(err)
		}
	}

	switch command {
	case cmdCodeVersion:
		fmt.Println(version)
		return
	case cmdCodeMigrate:
		return
	}

	router := routers.NewRouter(applicationRepository, applicationVersionRepository, version)
	panic(router.Run(fmt.Sprintf(":%d", *flagPort)))
}
