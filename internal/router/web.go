package routes

import (
	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/internal/handler/view"
	"github.com/duaminggu/sijiden/internal/middleware"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

func RegisterWebRoutes(e *echo.Echo, client *ent.Client, store *session.SessionStore) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Sijiden App is running 🚀")
	})

	sijidenGroup := e.Group("/sijiden")
	sijidenGroup.Use(middleware.RequireLoginView(store))
	sijidenGroup.Use(middleware.RequireRole("admin", store))

	sijidenGroup.GET("", view.DashboardPage(store))

	sijidenUserGroup := sijidenGroup.Group("/users")
	sijidenUserGroup.GET("", view.UserListPage(store))
	sijidenUserGroup.GET("/create", view.UserCreatePage(store))
	sijidenUserGroup.GET("/update", view.UserUpdatePage(store))

	sijidenRoleGroup := sijidenGroup.Group("/roles")
	sijidenRoleGroup.GET("", view.RoleListPage())
	sijidenRoleGroup.GET("/create", view.RoleListPage())
	sijidenRoleGroup.GET("/update", view.RoleListPage())
	sijidenRoleGroup.GET("/:id/detail", view.RoleListPage())

	e.GET("/auth", view.AuthPage(store))

}
