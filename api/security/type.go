package security

import (
	"encoding/json"

	"github.com/gobwas/glob"
)

// Role defines an applicative Role
type Role string

// Policy defines a match table between roles, http headers and allowed values
type Policy map[Role]map[string][]interface{}

// RolePolicy defines an indexed list of Role
type RolePolicy map[Role]bool

// Existing roles
const (
	RoleAdmin        Role = "ROLE_ADMIN"
	RoleUser              = "ROLE_USER"
	RoleBadgeCreator      = "ROLE_BADGE_CREATOR"
)

// UnmarshalJSON implements json.Unmarshaler interface
func (p Policy) UnmarshalJSON(raw []byte) error {
	var jsonPolicy map[Role]map[string][]string
	if err := json.Unmarshal(raw, &jsonPolicy); err != nil {
		return err
	}
	for role, headers := range jsonPolicy {
		p[role] = map[string][]interface{}{}
		for header, patterns := range headers {
			p[role][header] = []interface{}{}
			for _, pattern := range patterns {
				var v interface{} = pattern
				g, err := glob.Compile(pattern)
				if err == nil {
					v = g
				}
				p[role][header] = append(p[role][header], v)
			}
		}
	}
	return nil
}
