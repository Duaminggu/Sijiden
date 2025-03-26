package helper

import (
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
