package config

import (
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/security"
)

// Lhasa is the main config format
type Lhasa struct {
	DB         db.DatabaseCredentials `json:"appcatalog-db"`
	Policy     security.Policy        `json:"security"`
	LogHeaders []string               `json:"log-headers"`
}
