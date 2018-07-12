package config

import (
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/security"
)

// Lhasa is the main config format
type Lhasa struct {
	DB         db.DatabaseCredentials  `json:"appcatalog-db"`
	Security   security.Policy         `json:"security"`
	Policy     security.CompiledPolicy `json:"-"`
	LogHeaders []string                `json:"log-headers"`
}
