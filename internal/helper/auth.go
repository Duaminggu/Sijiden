package helper

import (
	"github.com/labstack/echo/v4"
)

func GetCurrentUserAndRoles(c echo.Context) (int, []string) {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return 0, nil
	}

	roles, ok := c.Get("roles").([]string)
	if !ok {
		return userID, []string{}
	}

	return userID, roles
}
