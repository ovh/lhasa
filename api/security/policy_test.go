package security_test

import (
	"net/http"
	"testing"

	"github.com/ovh/lhasa/api/security"
)

func TestBuildRolePolicy(t *testing.T) {
	data := []struct {
		policy  security.RolePolicy
		roles   []security.Role
		noRoles []security.Role
	}{
		{
			policy: security.BuildRolePolicy(
				security.Compile(security.Policy{
					"ROLE_ADMIN": {
						"X-Remote-User": {"john.doe"},
					},
					"ROLE_USER": {
						"X-Remote-User": {"*"},
					},
				}),
				&http.Request{Header: map[string][]string{"X-Remote-User": {"john.doe"}}}),
			roles: []security.Role{security.RoleAdmin, security.RoleUser},
		},
		{
			policy: security.BuildRolePolicy(
				security.Compile(security.Policy{
					"ROLE_ADMIN": {
						"X-Remote-User": {"john.doe"},
					},
					"ROLE_USER": {
						"X-Remote-User": {"*"},
					},
				}),
				&http.Request{Header: map[string][]string{"X-Remote-User": {"foo.bar"}}}),
			roles:   []security.Role{security.RoleUser},
			noRoles: []security.Role{security.RoleAdmin},
		},
		{
			policy: security.BuildRolePolicy(
				security.Compile(security.Policy{
					"ROLE_ADMIN": {
						"X-Ovh-Gateway-Source": {"foobar"},
					},
					"ROLE_USER": {
						"X-Remote-User": {"*"},
					},
				}),
				&http.Request{Header: map[string][]string{"X-Remote-User": {"foo.bar"}}}),
			roles:   []security.Role{security.RoleUser},
			noRoles: []security.Role{security.RoleAdmin},
		},
		{
			policy: security.BuildRolePolicy(
				security.Compile(security.Policy{
					"ROLE_ADMIN": {
						"X-Ovh-Gateway-Source": {"foobar"},
					},
				}),
				&http.Request{Header: map[string][]string{"X-Ovh-Gateway-Source": {"foobar"}}}),
			roles:   []security.Role{security.RoleAdmin},
			noRoles: []security.Role{security.RoleUser},
		},
	}

	for i, run := range data {
		if !run.policy.HasAll(run.roles...) {
			t.Errorf("Test %d - Expected %v - Found %v", i+1, run.roles, run.policy.ToSlice())
		}
		if run.policy.HasOne(run.noRoles...) {
			t.Errorf("Test %d - Expected %v - Found %v", i+1, run.roles, run.policy.ToSlice())
		}
	}
}
