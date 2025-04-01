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

	admin.GET("/auth", func(c echo.Context) error {
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

		return c.Render(200, "sijiden/auth.html", echo.Map{
			"csrf_token": csrfToken,
		})
	})

	admin.GET("", func(c echo.Context) error {
		userID := c.Get("user_id")
		if userID == nil {
			return c.Redirect(302, "/sijiden/auth")
		}

		_, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(302, "/sijiden/auth")
		}

		return c.Render(200, "sijiden/dashboard.html", echo.Map{
			"username": "User",
		})
	}, middleware.RequireLoginView(store))

}
