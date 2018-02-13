package db

// The db package's main purpose is to trigger a database connection
// and offer it as a global variable to all the package consumers.
// Note: This is not the best way to share this variable but

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This is required by GORM to enable postgresql support
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // This is required by GORM to enable sqlite support
)

var db = dbconnect()

func dbconnect() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./test.db")
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
