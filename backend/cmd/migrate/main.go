package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

type MigrationConfig struct {
	DatabaseURL    string
	MigrationsPath string
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading: %v", err)
	}

	var config MigrationConfig
	
	// Default values
	config.DatabaseURL = os.Getenv("DATABASE_URL")
	config.MigrationsPath = "file://./migrations"
	
	// Override with flags if provided
	flag.StringVar(&config.DatabaseURL, "database-url", config.DatabaseURL, "Database connection URL")
	flag.StringVar(&config.MigrationsPath, "migrations-path", config.MigrationsPath, "Path to migration files")
	flag.Parse()
	
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	
	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	
	command := args[0]
	
	// Create migrate instance
	m, err := migrate.New(config.MigrationsPath, config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()
	
	switch command {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to run")
		} else {
			fmt.Println("Migrations applied successfully")
		}
		
	case "down":
		err = m.Steps(-1)
		if err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Last migration rolled back successfully")
		
	case "reset":
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration reset failed: %v", err)
		}
		fmt.Println("All migrations rolled back successfully")
		
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
		
	case "force":
		if len(args) < 2 {
			log.Fatal("Force command requires a version argument")
		}
		var version int
		_, err := fmt.Sscanf(args[1], "%d", &version)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		err = m.Force(version)
		if err != nil {
			log.Fatalf("Force migration failed: %v", err)
		}
		fmt.Printf("Forced migration to version %d\n", version)
		
	case "goto":
		if len(args) < 2 {
			log.Fatal("Goto command requires a version argument")
		}
		var targetVersion uint
		_, err := fmt.Sscanf(args[1], "%d", &targetVersion)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		err = m.Migrate(targetVersion)
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration to version %d failed: %v", targetVersion, err)
		}
		fmt.Printf("Migrated to version %d\n", targetVersion)
		
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database Migration CLI")
	fmt.Println()
	fmt.Println("Usage: go run cmd/migrate/main.go [flags] <command> [args]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -database-url string    Database connection URL (default: $DATABASE_URL)")
	fmt.Println("  -migrations-path string Path to migration files (default: file://./migrations)")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up          Run all pending migrations")
	fmt.Println("  down        Rollback last migration")
	fmt.Println("  reset       Rollback all migrations")
	fmt.Println("  version     Show current migration version")
	fmt.Println("  force N     Force migration to version N")
	fmt.Println("  goto N      Migrate to specific version N")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go up")
	fmt.Println("  go run cmd/migrate/main.go down")
	fmt.Println("  go run cmd/migrate/main.go version")
	fmt.Println("  go run cmd/migrate/main.go force 1")
	fmt.Println("  go run cmd/migrate/main.go goto 3")
}