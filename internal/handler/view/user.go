package view

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

func UserListPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/user/list.html", echo.Map{})
	}
}

func UserCreatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/user/form.html", echo.Map{})
	}
}

func UserUpdatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/user/form.html", echo.Map{})
	}
}
