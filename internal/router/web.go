package routes

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/middleware"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterWebRoutes(e *echo.Echo, store *session.SessionStore) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Sijiden App is running 🚀")
	})

	admin := e.Group("/sijiden")

	admin.GET("/login", func(c echo.Context) error {
		sessionID := uuid.NewString()
		csrfToken := uuid.NewString()
		store.SetCSRF(sessionID, csrfToken)

		c.SetCookie(&http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		})

		return c.Render(200, "login.html", echo.Map{
			"csrf_token": csrfToken,
		})
	})

	admin.GET("", func(c echo.Context) error {
		userID := c.Get("user_id")
		if userID == nil {
			return c.Redirect(302, "/login")
		}

		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/login")
		}

		csrfToken := uuid.NewString()
		store.SetCSRF(cookie.Value, csrfToken)

		return c.Render(200, "todos.html", echo.Map{
			"username":   "User",
			"csrf_token": csrfToken,
		})
	}, middleware.RequireLoginView(store))

	// Custom route start here

	e.GET("/login", func(c echo.Context) error {
		sessionID := uuid.NewString()
		csrfToken := uuid.NewString()
		store.SetCSRF(sessionID, csrfToken)

		c.SetCookie(&http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		})

		return c.Render(200, "login.html", echo.Map{
			"csrf_token": csrfToken,
		})
	})

	e.GET("/login", func(c echo.Context) error {
		sessionID := uuid.NewString()
		csrfToken := uuid.NewString()
		store.SetCSRF(sessionID, csrfToken)

		c.SetCookie(&http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		})

		return c.Render(200, "login.html", echo.Map{
			"csrf_token": csrfToken,
		})
	})

	e.GET("/register", func(c echo.Context) error {
		sessionID := uuid.NewString()
		csrfToken := uuid.NewString()
		store.SetCSRF(sessionID, csrfToken)

		c.SetCookie(&http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		})

		return c.Render(200, "register.html", echo.Map{
			"csrf_token": csrfToken,
		})
	})

	e.GET("/logout", func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err == nil {
			store.Delete(cookie.Value)
			c.SetCookie(&http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
		}
		return c.Redirect(302, "/login")
	})

}
