package migrate

import (
	"2margin.vn/account-service/config"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"math/rand"
	"time"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type migrate struct {
	sqlDB *sql.DB
}

func Migrator(cfg *config.Config) (*migrate, error) {

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlDbname,
		cfg.Postgres.PostgresqlPassword,
	)
	db, err := sqlx.Open(cfg.Postgres.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}
	m := &migrate{sqlDB: db.DB}

	goose.SetBaseFS(embedMigrations)

	return m, nil
}

func (m *migrate) Up() {
	alive(m.sqlDB)
	if err := goose.Up(m.sqlDB, "migrations"); err != nil {
		panic(err)
	}
}

func (m *migrate) Down() {
	if err := goose.Down(m.sqlDB, "migrations"); err != nil {
		panic(err)
	}
}

func alive(db *sql.DB) {
	log.Println("connecting to database... ")
	for {
		// Ping by itself is un-reliable, the connections are cached. This
		// ensures that the database is still running by executing a harmless
		// dummy query against it.
		//
		// Also, we perform an exponential backoff when querying the database
		// to spread our retries.
		_, err := db.Exec("SELECT true")
		if err == nil {
			log.Println("database connected")
			return
		}

		base, capacity := time.Second, time.Minute
		for backoff := base; err != nil; backoff <<= 1 {
			if backoff > capacity {
				backoff = capacity
			}

			// A pseudo-random number generator here is fine. No need to be
			// cryptographically secure. Ignore with the following comment:
			/* #nosec */
			jitter := rand.Int63n(int64(backoff * 3))
			sleep := base + time.Duration(jitter)
			time.Sleep(sleep)
			_, err := db.Exec("SELECT true")
			if err == nil {
				log.Println("database connected")
				return
			}
		}
	}
}
