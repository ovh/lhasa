package db

// The db package's main purpose is to trigger a database connection
// and offer it as a global variable to all the package consumers.
// Note: This is not the best way to share this variable but

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This is required by GORM to enable postgresql support
	"github.com/sirupsen/logrus"

	dbconfig "github.com/ovh/lhasa/api/dbconfig"
)

const (
	defaultTimeout = 5
	maxOpenConns   = 10
	maxIdleConns   = 3
)

var db *gorm.DB

func connectFromVault(log *logrus.Logger, vaultAlias string) *gorm.DB {
	// Init vault
	connConfigStr, err := autovault.Secrets().Alias(vaultAlias)
	if err != nil {
		log.WithError(err).Fatalf("cannot read alias %s from vault", vaultAlias)
	}
	connConfig, err := fromJSON(connConfigStr)
	if err != nil {
		log.WithError(err).Fatalf("cannot read JSON DBConfig from vault secret %s", vaultAlias)
	}
	connStr, err := connConfig.GetRW()
	if err != nil {
		log.WithError(err).Fatalf("cannot get a RW database from vault secret %s", vaultAlias)
	}

	return NewFromGormString(connStr, log, vaultAlias)
}

// NewFromGormString creates a gorm db handler from a connection string
func NewFromGormString(connStr string, log *logrus.Logger, vaultAlias string) *gorm.DB {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.WithError(err).Fatalf("cannot open database from vault secret %s with gorm", vaultAlias)
	}
	return db
}

// NewFromVault provides the database handle to its callers
func NewFromVault(log *logrus.Logger, logMode bool, vaultAlias string) *gorm.DB {
	db = connectFromVault(log, vaultAlias)
	db.LogMode(logMode)
	db.SetLogger(gorm.Logger{LogWriter: log})
	return db
}
