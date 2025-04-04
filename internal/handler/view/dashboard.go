package view

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

func DashboardPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/dashboard.html", echo.Map{})
	}
}
