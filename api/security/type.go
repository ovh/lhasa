package security

// Role defines an applicative Role
type Role string

// Policy defines a match table between roles, http headers and allowed values
type Policy map[Role]map[string][]string

// RolePolicy defines an indexed list of Role
type RolePolicy map[Role]bool

// Existing roles
const (
	RoleAdmin Role = "ROLE_ADMIN"
	RoleUser       = "ROLE_USER"
)
