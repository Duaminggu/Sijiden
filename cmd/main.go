package main

import (
	"github.com/duaminggu/sijiden/internal/db"
	routes "github.com/duaminggu/sijiden/internal/router"
	sessionstore "github.com/duaminggu/sijiden/internal/session"
	"github.com/duaminggu/sijiden/internal/template"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func main() {
	client := db.NewClient()

	store := sessionstore.NewStore()

	e := echo.New()
	// Register routes
	routes.RegisterWebRoutes(e, store)
	routes.RegisterAjaxRoutes(e, client, store)
	routes.RegisterComponentRoutes(e)

	e.Renderer = template.NewRenderer()
	e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":1234"))
}
