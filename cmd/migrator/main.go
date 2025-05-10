package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dbURL, migrationPath, migrationTable string

	flag.StringVar(&dbURL, "db-url", "", "PostgreSQL connection URL (postgres://user:pass@host:port/dbname)")
	flag.StringVar(&migrationPath, "migrations-path", "./migrations", "path to migration files")
	flag.StringVar(&migrationTable, "migrations-table", "schema_migrations", "name for migrations table")
	flag.Parse()

	if dbURL == "" {
		dbURL = os.Getenv("DB_URL")
		if dbURL == "" {
			log.Fatal("Database URL is required. Use -db-url flag or DB_URL environment variable")
		}
	}

	if migrationPath == "" {
		log.Fatal("Migrations path is required")
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		fmt.Sprintf("%s?x-migrations-table=%s&sslmode=disable", dbURL, migrationTable))
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}
	defer m.Close()

	log.Println("Applying migrations...")
	err = m.Up()
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		log.Println("No new migrations to apply")
	case err != nil:
		log.Fatalf("Migration failed: %v", err)
	default:
		log.Println("Migrations applied successfully")
	}
}
