package config

// The db package's main purpose is to trigger a database connection
// and offer it as a global variable to all the package consumers.
// Note: This is not the best way to share this variable but

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This is required by GORM to enable postgresql support
	"github.com/ovh/lhasa/api/db"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	maxOpenConns = 10
	maxIdleConns = 3
)

var (
	// Config store all configuration items
	config = make(map[string]interface{})
)

// read a file
func readFile(configFile string) (string, error) {
	b, err := ioutil.ReadFile(configFile) // just pass the file name
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// LoadFromFile extract configuration file
func LoadFromFile(configFile *os.File) error {
	// Init config file
	connConfigStr, err := readFile(configFile.Name())
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(connConfigStr), &config)
}

// ExtractValue return value from key
func ExtractValue(key string) interface{} {
	// Transform it to string
	return config[key]
}

// ExtractKey return o json formated value (in pretty format)
func extractKey(key string) (string, error) {
	// Transform it to string
	payload, err := json.MarshalIndent(config[key], "", "\t")
	if err != nil {
		return "{}", errors.Wrapf(err, "Unable to marshal:", err)
	}
	return string(payload), nil
}

// NewDBHandle provides the database handle to its callers
func NewDBHandle(configFile *os.File, dbAlias string, logMode bool, log *logrus.Logger) (*gorm.DB, error) {
	connConfigStr, err := extractKey(dbAlias)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read alias %s from configuration file", configFile)
	}
	connConfig, err := db.FromJSON(connConfigStr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read JSON DBConfig from db alias %s", dbAlias)
	}
	connStr, err := connConfig.GetRW()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get a RW database from db alias %s", dbAlias)
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
	db.SetLogger(gorm.Logger{LogWriter: dbLogWriter{log}})
	return db, nil
}

type dbLogWriter struct {
	log *logrus.Logger
}

// Println implements gorm's LogWriter interface
func (l dbLogWriter) Println(v ...interface{}) {
	if l.log == nil {
		return
	}
	if len(v) != 5 {
		return
	}

	var message interface{} = v
	file := strings.Split(strings.TrimSuffix(strings.TrimPrefix(v[0].(string), "\u001b[35m("), ")\u001b[0m"), ":")
	durationStr := strings.TrimSuffix(strings.TrimPrefix(v[2].(string), " \u001b[36;1m["), "]\u001b[0m ")
	rows, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(v[4].(string), " \n\u001b[36;31m["), " rows affected or returned ]\u001b[0m "))
	duration, _ := time.ParseDuration(durationStr)
	line, _ := strconv.Atoi(file[1])
	fields := logrus.Fields{
		"file":         file[0],
		"line":         line,
		"duration":     duration.Seconds(),
		"full_message": v[3],
		"rows":         rows,
	}
	switch v[3].(type) {
	case string:
		message = fmt.Sprintf("sql query: %s...", v[3].(string)[:int(math.Min(50, float64(len(v[3].(string)))))])
		l.log.WithFields(fields).Debug(message)
	case error:
		l.log.WithFields(fields).WithError(v[3].(error)).Error(message)
	default:
		l.log.WithFields(fields).Debug(message)
	}
}
