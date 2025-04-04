package main

import (
	"os"
	"strconv"

	"github.com/duaminggu/sijiden/internal/db"
	routes "github.com/duaminggu/sijiden/internal/router"
	sessionstore "github.com/duaminggu/sijiden/internal/session"
	"github.com/duaminggu/sijiden/internal/template"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo/v4"
)

func main() {
	client := db.NewClient()

	store := sessionstore.NewStore()

	e := echo.New()
	// Register routes
	routes.RegisterWebRoutes(e, client, store)
	routes.RegisterAjaxRoutes(e, client, store)
	routes.RegisterComponentRoutes(e)

	isDev := os.Getenv("IS_DEV")

	e.Renderer = template.SijidenRenderer(stringToBoolWithDefault(isDev, false))
	e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":1234"))
}

func stringToBoolWithDefault(s string, defaultValue bool) bool {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return value
}
