package db

import (
	"context"
	"database/sql"
	"log"
	"os"

	"entgo.io/ent/dialect"
	entSql "entgo.io/ent/dialect/sql"
	"github.com/duaminggu/sijiden/ent"
	"github.com/duaminggu/sijiden/ent/migrate"
	_ "github.com/go-sql-driver/mysql"
)

func NewClient() *ent.Client {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root@tcp(127.0.0.1:3306)/todogo?parseTime=True"
	}

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed opening mysql connection: %v", err)
	}

	drv := entSql.OpenDB(dialect.MySQL, sqlDB)
	client := ent.NewClient(ent.Driver(drv))

	// Jalankan migrasi
	ctx := context.Background()
	if err := client.Schema.Create(ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed creating schema: %v", err)
	}

	return client
}
