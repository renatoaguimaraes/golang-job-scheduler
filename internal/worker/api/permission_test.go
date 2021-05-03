package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPermissionExistsMethodAndRole(t *testing.T) {
	permitted := HasPermission("/WorkerService/Start", []string{"admin"})

	assert.True(t, permitted)
}

func TestHasPermissionNotExistsMethod(t *testing.T) {
	permitted := HasPermission("/WorkerService/NotExists", []string{"admin"})

	assert.False(t, permitted)
}

func TestHasPermissionNotExistsRole(t *testing.T) {
	permitted := HasPermission("/WorkerService/Start", []string{"notexist"})

	assert.False(t, permitted)
}

func TestHasPermissionInvalidRole(t *testing.T) {
	permitted := HasPermission("/WorkerService/Start", []string{"user"})

	assert.False(t, permitted)
}

func TestParseRoles(t *testing.T) {
	roles := ParseRoles("admin,user")

	assert.ElementsMatch(t, []string{"admin", "user"}, roles)
}

func TestParseRolesLineBread(t *testing.T) {
	roles := ParseRoles("\nadmin,user")

	assert.ElementsMatch(t, []string{"admin", "user"}, roles)
}

func TestParseRolesCarriageReturn(t *testing.T) {
	roles := ParseRoles("admin,user\r")

	assert.ElementsMatch(t, []string{"admin", "user"}, roles)
}

func TestOidToString(t *testing.T) {
	oid := OidToString([]int{1, 2, 840, 10070, 8, 1})

	assert.Equal(t, "1.2.840.10070.8.1", oid)
}

func TestIsOidRole(t *testing.T) {
	isrole := IsOidRole("1.2.840.10070.8.1")

	assert.True(t, isrole)
}

func TestIsOidRoleInvalid(t *testing.T) {
	isrole := IsOidRole("9.1.820.105070.8.1")

	assert.False(t, isrole)
}
