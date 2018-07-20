package db

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This is required by GORM to enable postgresql support
	"github.com/sirupsen/logrus"
)

const (
	maxOpenConns = 10
	maxIdleConns = 3
)

// NewDBHandle provides the database handle to its callers
func NewDBHandle(dc DatabaseCredentials, logMode bool, log logrus.FieldLogger) (*gorm.DB, error) {
	connStr, err := dc.GetRW()
	if err != nil {
		return nil, err
	}
	return NewFromGormString(connStr, logMode, log)
}

// NewFromGormString creates a gorm db handler from a connection string
func NewFromGormString(connStr string, logMode bool, log logrus.FieldLogger) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.LogMode(logMode)
	db.SetLogger(getLogger(log))
	return db, nil
}

// GetRW get a read/write database
func (dc *DatabaseCredentials) GetRW() (string, error) {
	return dc.getConnStr(dc.Writers)
}

// GetRO get a read only database
func (dc *DatabaseCredentials) GetRO() (string, error) {
	return dc.getConnStr(dc.Readers)
}

func (dc *DatabaseCredentials) getConnStr(instances []DatabaseInstance) (string, error) {
	dbType, err := dc.getType()
	if err != nil {
		return "", err
	}
	i, err := getRandom(instances)
	if err != nil {
		return "", err
	}
	return buildConnStr(dbType, dc, i), nil
}

func (dc *DatabaseCredentials) getType() (Type, error) {
	switch strings.ToLower(dc.Type) {
	case "postgresql":
		return PostgreSQL, nil
	}
	return "", fmt.Errorf("unsupported DB type '%s'", dc.Type)
}

func (dc *DatabaseCredentials) getSslDefaultMode(value string) (string, error) {
	if len(value) > 0 {
		return value, nil
	}
	switch strings.ToLower(dc.Type) {
	case "postgresql":
		return PostgreSQLDefaultSslMode, nil
	}
	return "", fmt.Errorf("unsupported DB type '%s'", dc.Type)
}

func getRandom(instances []DatabaseInstance) (*DatabaseInstance, error) {
	max := len(instances)
	if max == 0 {
		return nil, errors.New("no suitable db instance found")
	}
	rand.Seed(time.Now().Unix())
	return &instances[rand.Intn(max)], nil
}

func buildConnStr(fmtStr Type, dc *DatabaseCredentials, i *DatabaseInstance) string {
	// build sslmode with default value according do bdd type
	var sslmode, _ = dc.getSslDefaultMode(i.Ssl)
	return fmt.Sprintf(string(fmtStr), dc.User, dc.Password, i.Host, i.Port, dc.Database, sslmode)
}
