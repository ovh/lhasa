package db

// The db package's main purpose is to trigger a database connection
// and offer it as a global variable to all the package consumers.
// Note: This is not the best way to share this variable but

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This is required by GORM to enable postgresql support
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // This is required by GORM to enable sqlite support

	dbconfig "github.com/ovh/lhasa/api/dbconfig"
)

const (
	defaultTimeout = 5
	maxOpenConns   = 10
	maxIdleConns   = 3
	// DBSecretAlias db alias
	DBSecretAlias = "appcatalog-db"
)

var db = dbconnect()

func dbconnect() *gorm.DB {
	// Init vault
	connConfigStr, err := autovault.Secrets().Alias(DBSecretAlias)
	if err != nil {
		panic(err)
	}
	connConfig, err := dbconfig.FromJSON(connConfigStr)
	if err != nil {
		panic(err)
	}
	connStr, err := connConfig.GetRW()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	// Close the database after usage

	return db
}

// DB provides the database handle to its callers
func DB() *gorm.DB {
	return db
	// FIXME: Depending on the database keepalive config, we may
	// need Implement reconnection in case of connection loss
}
