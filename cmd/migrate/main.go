package main

import (
	"bLink-app/config"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// ---------- FLAGS ----------
	var action string
	var version int

	flag.StringVar(&action, "action", "up", "Migration action: up, down, force, version")
	flag.IntVar(&version, "version", -1, "Migration version for 'force' action")

	flag.Parse()
	// -----------------------------

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create migration driver with config
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    cfg.Postgres.Database,
	})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	// Migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}

	// ---------- EXECUTE ACTION ----------
	switch action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migration up:", err)
		}
		log.Println("Migration up completed successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migration down:", err)
		}
		log.Println("Migration down completed successfully")

	case "force":
		if version == -1 {
			log.Fatal("Version is required for force action. Example: -version=5")
		}

		if err := m.Force(version); err != nil {
			log.Fatal("Failed to force migration:", err)
		}
		log.Printf("Migration forced to version %d\n", version)

	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get migration version:", err)
		}
		log.Printf("Current version: %d, Dirty: %v\n", v, dirty)

	default:
		log.Fatal("Invalid action. Use: up, down, force, or version")
	}
}
