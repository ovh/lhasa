package db

// DatabaseInstance host and port
type DatabaseInstance struct {
	Port int    `json:"port"`
	Host string `json:"host"`
	Ssl  string `json:"sslmode"`
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

// Type db type
type Type string

const (
	// PostgreSQL connect string
	PostgreSQL Type = "user=%s password=%s host=%s port=%d DB.name=%s sslmode=%s"

	// PostgreSQLDefaultSslMode default ssl connect string
	PostgreSQLDefaultSslMode = "require"
)
