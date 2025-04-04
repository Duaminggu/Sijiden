package contextx

import (
	"github.com/duaminggu/sijiden/ent"
	"github.com/labstack/echo/v4"
)

// GetUserID safely extracts user_id from echo.Context
func GetUserID(c echo.Context) int {
	userID := c.Get("user_id")
	if id, ok := userID.(int); ok && id > 0 {
		return id
	}
	return 0
}

func GetUserIDWithStatus(c echo.Context) (int, bool) {
	userID := c.Get("user_id")
	id, ok := userID.(int)
	if !ok || id <= 0 {
		return 0, false
	}
	return id, true
}

func GetRoles(c echo.Context) []string {
	roles, ok := c.Get("roles").([]string)
	if !ok {
		return []string{}
	}
	return roles
}

func GetUser(c echo.Context) (*ent.User, bool) {
	user, ok := c.Get("user").(*ent.User)
	return user, ok
}
