package routes

import (
	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/internal/handler"
	"github.com/duaminggu/sijiden/internal/middleware"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

func RegisterAjaxRoutes(e *echo.Echo, client *ent.Client, store *session.SessionStore) {
	// Route login harus di luar middleware RequireLoginAPI
	e.POST("/ajax/sijiden/login", handler.Login(client, store), middleware.RequireCSRF(store))

	// Hanya route yang sudah butuh login dikelompokkan ke sini
	api := e.Group("/ajax", middleware.RequireLoginAPI(store))
	api.GET("/user", handler.GetUsers(client))

	// sijiden := api.Group("/sijiden")
	// tambahkan route lain di sini yang butuh login
}
