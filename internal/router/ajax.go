package routes

import (
	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/internal/handler"
	"github.com/duaminggu/sijiden/internal/middleware"
	"github.com/duaminggu/sijiden/internal/session"
	"github.com/labstack/echo/v4"
)

func RegisterAjaxRoutes(e *echo.Echo, client *ent.Client, store *session.SessionStore) {
	api := e.Group("/ajax", middleware.RequireLoginAPI(store))
	api.GET("/user", handler.GetUsers(client))
}
