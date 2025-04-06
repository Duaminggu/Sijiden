package view

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RoleListPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		csrfToken := uuid.NewString()
		store.SetCSRF(cookie.Value, csrfToken)
		return c.Render(http.StatusOK, "sijiden/role/list.html", echo.Map{
			"csrf_token": csrfToken,
		})
	}
}

func RoleDetailPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		csrfToken := uuid.NewString()
		store.SetCSRF(cookie.Value, csrfToken)

		id := c.Param("id")

		return c.Render(http.StatusOK, "sijiden/role/detail.html", echo.Map{
			"csrf_token": csrfToken,
			"role_id":    id,
		})
	}
}

func RoleCreatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		csrfToken := uuid.NewString()
		store.SetCSRF(cookie.Value, csrfToken)
		return c.Render(http.StatusOK, "sijiden/role/form.html", echo.Map{
			"csrf_token": csrfToken,
		})
	}
}

func RoleUpdatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		csrfToken := uuid.NewString()
		store.SetCSRF(cookie.Value, csrfToken)

		id := c.Param("id")

		return c.Render(http.StatusOK, "sijiden/role/form.html", echo.Map{
			"csrf_token": csrfToken,
			"role_id":    id, // dikirim ke frontend
		})
	}
}
