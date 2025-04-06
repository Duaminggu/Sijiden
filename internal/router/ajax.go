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
	roleHandler := &ajax.RoleHandler{Client: client, Store: store}
	// Route login harus di luar middleware RequireLoginAPI
	sijidenGroup := e.Group("/ajax/sijiden")

	sijidenGroup.POST("/login", ajax.Login(client, store), middleware.RequireCSRF(store))

	sijidenUserGroup := sijidenGroup.Group("/users")
	sijidenUserGroup.GET("/count", userHandler.CountUsers, middleware.RequireLoginAjax(store))

	sijidenRoleGroup := sijidenGroup.Group("/roles")
	sijidenRoleGroup.Use(middleware.RequireLoginAjax(store))
	sijidenRoleGroup.Use(middleware.RequireRole("admin", store))
	sijidenRoleGroup.GET("", roleHandler.List)
	sijidenRoleGroup.POST("/create", roleHandler.Create)
	sijidenRoleGroup.GET("/:id", roleHandler.Detail)
	sijidenRoleGroup.PUT("/:id/update", roleHandler.Update)
	sijidenRoleGroup.DELETE("/:id", roleHandler.Delete)

	// Hanya route yang sudah butuh login dikelompokkan ke sini
	ajaxGroup := e.Group("/ajax", middleware.RequireLoginAjax(store))
	ajaxGroup.GET("/user", ajax.GetUsers(client))

	// sijiden := api.Group("/sijiden")
	// tambahkan route lain di sini yang butuh login
}
