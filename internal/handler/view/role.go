package view

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

func RoleListPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/role/list.html", echo.Map{})
	}
}

func RoleCreatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/role/form.html", echo.Map{})
	}
}

func RoleUpdatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/role/form.html", echo.Map{})
	}
}
