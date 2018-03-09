package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// DatabaseInstance host and port
type DatabaseInstance struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// DatabaseCredentials credentials
type DatabaseCredentials struct {
	Readers  []DatabaseInstance `json:"readers"`
	Writers  []DatabaseInstance `json:"writers"`
	Database string             `json:"database"`
	Password string             `json:"password"`
	User     string             `json:"user"`
	Type     string             `json:"type"`
}

// FromJSON unmarshall creds
func FromJSON(creds string) (*DatabaseCredentials, error) {
	dc := &DatabaseCredentials{}
	err := json.Unmarshal([]byte(creds), &dc)
	if err != nil {
		return nil, ErrUnmarshal
	}
	return dc, nil
}

// Type db type
type Type string

var (
	// PostgreSQL connect string
	PostgreSQL Type = "user=%s password=%s host=%s port=%d DB.name=%s sslmode=require"

	// ErrUnmarshal error
	ErrUnmarshal = errors.New("Unmarshaling error")

	// ErrNoInstanceFound instance not found
	ErrNoInstanceFound = errors.New("No suitable db instance found")
)

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
	return "", fmt.Errorf("Unknown DB type '%s'", dc.Type)
}

func getRandom(instances []DatabaseInstance) (*DatabaseInstance, error) {
	if len(instances) == 0 {
		return nil, ErrNoInstanceFound
	}
	return &instances[0], nil // TODO rnd
}

func buildConnStr(fmtStr Type, dc *DatabaseCredentials, i *DatabaseInstance) string {
	return fmt.Sprintf(string(fmtStr), dc.User, dc.Password, i.Host, i.Port, dc.Database)
}
