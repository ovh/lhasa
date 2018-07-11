package security

import (
	"net/http"

	"github.com/gobwas/glob"
)

// HasOne returns a gin handler that checks the request against the given role
func (policy RolePolicy) HasOne(roles ...Role) bool {
	for _, role := range roles {
		if policy[role] {
			return true
		}
	}
	return false
}

// HasAll returns a gin handler that checks the request against the given role
func (policy RolePolicy) HasAll(roles ...Role) bool {
	for _, role := range roles {
		if !policy[role] {
			return false
		}
	}
	return true
}

// ToSlice returns a slice of roles
func (policy RolePolicy) ToSlice() (roles []Role) {
	for role, ok := range policy {
		if ok {
			roles = append(roles, role)
		}
	}
	return
}

// Compile a policy while compiling values as glob
func Compile(policy Policy) CompiledPolicy {
	p := CompiledPolicy{}
	for role, headers := range policy {
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
	return p
}

// BuildRolePolicy returns a map of Role matching the given http request
func BuildRolePolicy(policy CompiledPolicy, r *http.Request) RolePolicy {
	if r == nil {
		return nil
	}
	roles := map[Role]bool{}
	for role, headers := range policy {
		for header, patterns := range headers {
			for _, pattern := range patterns {
				if checkPattern(pattern, r.Header.Get(header)) {
					roles[role] = true
				}
			}
		}
	}
	return roles
}

func checkPattern(pattern interface{}, value string) bool {
	switch p := pattern.(type) {
	case glob.Glob:
		return p.Match(value)
	case string:
		return value == p
	}
	return false
}
