package security

// Role defines an applicative Role
type Role string

// Policy defines a match table between roles, http headers and allowed values
type Policy map[Role]map[string][]string

// CompiledPolicy is a policy where values have been compiled as globs
type CompiledPolicy map[Role]map[string][]interface{}

// RolePolicy defines an indexed list of Role
type RolePolicy map[Role]bool

// Existing roles
const (
	RoleAdmin Role = "ROLE_ADMIN"
	RoleUser       = "ROLE_USER"
)

// User is a human user
type User struct {
	name string
}
