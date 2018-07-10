package security

import (
	"net/http"
)

// Various security headers
const (
	HeaderGatewaySource = "X-Ovh-Gateway-Source"
	HeaderRequestID     = "X-Request-Id"
	HeaderRemoteUser    = "X-Remote-User"
)

// HasOneRoleOf returns a gin handler that checks the request against the given role
func (policy RolePolicy) HasOneRoleOf(roles ...Role) bool {
	for _, role := range roles {
		if policy[role] {
			return true
		}
	}
	return false
}

// HasAllRoles returns a gin handler that checks the request against the given role
func (policy RolePolicy) HasAllRoles(roles ...Role) bool {
	for _, role := range roles {
		if !policy[role] {
			return false
		}
	}
	return true
}

func (policy RolePolicy) ToSlice() (roles []Role) {
	for role, ok := range policy {
		if ok {
			roles = append(roles, role)
		}
	}
	return
}

// BuildRolePolicy returns a map of Role matching the given http request
func BuildRolePolicy(policy Policy, r *http.Request) RolePolicy {
	if r == nil {
		return nil
	}
	roles := map[Role]bool{}
	for role, headers := range policy {
		for header, values := range headers {
			for _, value := range values {
				if r.Header.Get(header) == value {
					roles[role] = true
				}
			}
		}
	}
	if r.Header.Get(HeaderRemoteUser) != "" {
		roles[RoleUser] = true
	}
	return roles
}
