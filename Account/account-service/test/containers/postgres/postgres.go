package postgres

import (
	"2margin.vn/account-service/config"
	"2margin.vn/account-service/migrate"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"time"

	"log"
	"os"
	"strings"
	"testing"
)

var (
	PostgresqlUser     = "postgres"
	PostgresqlPassword = "postgres"
	PostgresqlDbname   = "postgres"
)

type Container struct {
	Cfg *config.Config
	m   *testing.M
}

func NewTestContainer(m *testing.M) *Container {
	return &Container{m: m}
}

func (c *Container) Run() error {

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	pool.MaxWait = 120 * time.Second

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a Container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_PASSWORD=" + PostgresqlPassword,
			"POSTGRES_USER=" + PostgresqlUser,
			"POSTGRES_DB=" + PostgresqlDbname,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped Container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	hostAndPort := resource.GetHostPort("5432/tcp")

	var items = strings.Split(hostAndPort, ":")

	_ = resource.Expire(120) // Tell docker to hard kill the Container in 120 seconds

	cfgPostgres := config.Postgres{PostgresqlHost: items[0], PostgresqlPort: items[1], PostgresqlUser: PostgresqlUser, PostgresqlPassword: PostgresqlPassword, PostgresqlDbname: PostgresqlDbname, PostgresqlSSLMode: false, PgDriver: "pgx"}

	cfg := &config.Config{Postgres: cfgPostgres}
	c.Cfg = cfg
	migrator, err := migrate.Migrator(cfg)
	if err != nil {
		return err
	}

	// exponential backoff-retry, because the application in the Container might not be ready to accept connections yet

	// Performing a migration this way means all tests in this package shares
	// the same db schema across all unit test.
	// If isolation is needed, then do away with using `testing.M`. Do a
	// migration for each test handler instead.
	migrator.Up()

	// We can access database with m.hostAndPort or m.databaseUrl
	// port changes everytime a new docker instance is run
	code := c.m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)

	return nil

}
