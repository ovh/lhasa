package db

// The db package's main purpose is to trigger a database connection
// and offer it as a global variable to all the package consumers.
// Note: This is not the best way to share this variable but

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This is required by GORM to enable postgresql support
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	dbconfig "github.com/ovh/lhasa/api/dbconfig"
)

const (
	maxOpenConns = 10
	maxIdleConns = 3
)

// NewFromVault provides the database handle to its callers
func NewFromVault(vaultAlias string, logMode bool, log *logrus.Logger) (*gorm.DB, error) {
	// Init vault
	connConfigStr, err := autovault.Secrets().Alias(vaultAlias)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read alias %s from vault", vaultAlias)
	}
	connConfig, err := fromJSON(connConfigStr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read JSON DBConfig from vault secret %s", vaultAlias)
	}
	connStr, err := connConfig.GetRW()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get a RW database from vault secret %s", vaultAlias)
	}
	return NewFromGormString(connStr, logMode, log)
}

// NewFromGormString creates a gorm db handler from a connection string
func NewFromGormString(connStr string, logMode bool, log *logrus.Logger) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.LogMode(logMode)
	db.SetLogger(gorm.Logger{LogWriter: log})
	return db, nil
}
