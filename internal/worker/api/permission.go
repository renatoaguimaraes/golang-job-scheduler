package api

import (
	"strconv"
	"strings"
)

// oidRole oid identifier used to store user roles
const oidRole string = "1.2.840.10070.8.1"

// permissions map initialization
var permissions = map[string][]string{
	"/WorkerService/Start":  {"admin"},
	"/WorkerService/Stop":   {"admin"},
	"/WorkerService/Query":  {"admin", "user"},
	"/WorkerService/Stream": {"admin", "user"},
}

// HasPermission verifies the permission given a method and user roles
func HasPermission(method string, roles []string) bool {
	permission, ok := permissions[method]
	if !ok {
		return false
	}
	for _, role := range roles {
		for _, value := range permission {
			if role == value {
				return true
			}
		}
	}
	return false
}

// IsOidRole validates the role oid
func IsOidRole(oid string) bool {
	return oidRole == oid
}

// ParseRoles split the roles string by comma
func ParseRoles(roles string) []string {
	return strings.Split(strings.TrimSpace(roles), ",")
}

// OidToString convert the int[] to string with
// point separator between the values
func OidToString(oid []int) string {
	var strs []string
	for _, value := range oid {
		strs = append(strs, strconv.Itoa(value))
	}
	return strings.Join(strs, ".")
}
