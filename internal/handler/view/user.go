package view

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserListPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "sijiden/user/list.html", echo.Map{})
	}
}

func UserCreatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		csrfToken := uuid.NewString()
		store.SetCSRF(cookie.Value, csrfToken)
		return c.Render(http.StatusOK, "sijiden/user/form.html", echo.Map{
			"csrf_token": csrfToken,
		})
	}
}

func UserUpdatePage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		csrfToken := uuid.NewString()
		store.SetCSRF(cookie.Value, csrfToken)

		id := c.Param("id")
		return c.Render(http.StatusOK, "sijiden/user/form.html", echo.Map{
			"csrf_token": csrfToken,
			"user_id":    id,
		})
	}
}

func UserDetailPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		return c.Render(http.StatusOK, "sijiden/user/detail.html", echo.Map{
			"user_id": id,
		})
	}
}
