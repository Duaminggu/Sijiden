package routes

import (
	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/internal/handler/ajax"
	"github.com/duaminggu/sijiden/internal/middleware"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

func RegisterAjaxRoutes(e *echo.Echo, client *ent.Client, store *session.SessionStore) {
	userHandler := &ajax.UserHandler{Client: client}
	// Route login harus di luar middleware RequireLoginAPI
	sijidenGroup := e.Group("/ajax/sijiden")
	sijidenGroup.POST("/login", ajax.Login(client, store), middleware.RequireCSRF(store))

	sijidenGroup.GET("/users/count", userHandler.CountUsers)

	// Hanya route yang sudah butuh login dikelompokkan ke sini
	ajaxGroup := e.Group("/ajax", middleware.RequireLoginAPI(store))
	ajaxGroup.GET("/user", ajax.GetUsers(client))

	// sijiden := api.Group("/sijiden")
	// tambahkan route lain di sini yang butuh login
}
