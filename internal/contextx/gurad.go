package contextx

import (
	"github.com/labstack/echo/v4"
)

// HasRole checks whether the current user has the specified role
func HasRole(c echo.Context, target string) bool {
	roles := GetRoles(c)
	for _, role := range roles {
		if role == target {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has at least one of the specified roles
func HasAnyRole(c echo.Context, allowedRoles ...string) bool {
	userRoles := GetRoles(c)
	roleMap := make(map[string]struct{})
	for _, r := range userRoles {
		roleMap[r] = struct{}{}
	}
	for _, allowed := range allowedRoles {
		if _, ok := roleMap[allowed]; ok {
			return true
		}
	}
	return false
}

func IsLoggedIn(c echo.Context) bool {
	id := GetUserID(c)
	return id > 0
}
