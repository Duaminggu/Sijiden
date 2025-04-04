package view

import (
	"net/http"

	"github.com/duaminggu/sijiden/internal/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func AuthPage(store *session.SessionStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID := uuid.NewString()
		csrfToken := uuid.NewString()

		store.SetCSRF(sessionID, csrfToken)

		c.SetCookie(&http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // Ganti ke true kalau pakai HTTPS
		})

		return c.Render(http.StatusOK, "sijiden/auth.html", echo.Map{
			"csrf_token": csrfToken,
		})
	}
}
