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
	"github.com/duaminggu/sijiden/internal/seeder"
	_ "github.com/go-sql-driver/mysql"
)

func NewClient() *ent.Client {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalf("failed to connect, because DB_DSN is empty")
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

	// Setelah migrasi
	if os.Getenv("SEED_ENABLED") == "true" {
		log.Println("📦 Seeding enabled. Running seed...")
		if err := seeder.SeedIfNeeded(client); err != nil {
			log.Fatalf("failed seeding data: %v", err)
		}
	}

	return client
}
