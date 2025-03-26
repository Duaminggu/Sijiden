package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func RegisterComponentRoutes(e *echo.Echo) {
	e.GET("/components/:name", func(c echo.Context) error {
		name := c.Param("name")
		return c.Render(200, fmt.Sprintf("component_%s.html", name), echo.Map{
			"username": "User",
		})
	})
}
