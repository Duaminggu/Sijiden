package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RoleListPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/role/list.html", echo.Map{})
	}
}

func RoleDetailPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/role/detail.html", echo.Map{})
	}
}

func RoleCreatePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/role/form.html", echo.Map{})
	}
}

func RoleUpdatePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/role/form.html", echo.Map{})
	}
}
