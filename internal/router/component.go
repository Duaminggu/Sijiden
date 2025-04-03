package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func RegisterComponentRoutes(e *echo.Echo) {
	sijiden := e.Group("/sijiden/components")
	sijiden.GET("/navbar", func(c echo.Context) error {
		return c.Render(200, "sijiden/components/navbar.html", echo.Map{
			"username": "User",
		})
	})

	e.GET("/components/:name", func(c echo.Context) error {
		name := c.Param("name")
		return c.Render(200, fmt.Sprintf("components/%s.html", name), echo.Map{
			"username": "User",
		})
	})

}
