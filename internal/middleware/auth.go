package middleware

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/session"

	"github.com/labstack/echo/v4"
)

func RequireLoginView(store *session.SessionStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				return c.Redirect(302, "/login")
			}

			if userID, ok := store.Get(cookie.Value); ok && userID > 0 {
				c.Set("user_id", userID)
				return next(c)
			}

			return c.Redirect(302, "/login")
		}
	}
}

func RequireLoginAPI(store *session.SessionStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}

			if userID, ok := store.Get(cookie.Value); ok && userID > 0 {
				c.Set("user_id", userID)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}
	}
}
