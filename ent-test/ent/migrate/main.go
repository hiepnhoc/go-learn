//go:build ignore

package main

import (
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/entc/integration/migrate/versioned/migrate"
	"log"
	"os"

	atlas "ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
)

func main() {
	ctx := context.Background()
	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := sqltool.NewFlywayDir("ent/migrate/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}
	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}
	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	err = migrate.NamedDiff(ctx, "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable", os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
