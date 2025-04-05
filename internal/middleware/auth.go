package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/duaminggu/sijiden/internal/session"

	"github.com/labstack/echo/v4"
)

func RequireLoginView(store *session.SessionStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				return c.Redirect(302, "/sijiden/auth")
			}

			if userID, ok := store.Get(cookie.Value); ok && userID > 0 {
				c.Set("user_id", userID)
				return next(c)
			}

			return c.Redirect(302, "/sijiden/auth")
		}
	}
}

func RequireLoginAjax(store *session.SessionStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				log.Printf("Missing or empty session cookie")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}

			userID, ok := store.Get(cookie.Value)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
			c.Set("user_id", userID)

			// Inject roles
			rolesStr, ok := store.GetValue(cookie.Value, "roles")
			if ok && rolesStr != "" {
				roleList := strings.Split(rolesStr, ",")
				c.Set("roles", roleList)
			}

			return next(c)
		}
	}
}

func RequireCSRF(store *session.SessionStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "no session"})
			}

			sentToken := c.Request().Header.Get("X-CSRF-Token")
			expectedToken, ok := store.GetCSRF(cookie.Value)
			if !ok || sentToken != expectedToken {
				return c.JSON(http.StatusForbidden, echo.Map{"error": "invalid csrf token"})
			}

			return next(c)
		}
	}
}

func RequireRole(roleName string, store *session.SessionStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				return c.Redirect(302, "/auth")
			}

			rolesStr, ok := store.GetValue(cookie.Value, "roles")
			if !ok {
				return c.Redirect(302, "/auth")
			}

			roles := strings.Split(rolesStr, ",")
			for _, role := range roles {
				if role == roleName {
					return next(c)
				}
			}

			return c.String(http.StatusForbidden, "Forbidden: You don't have the required role.")
		}
	}
}
