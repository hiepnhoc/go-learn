package repository

import (
	"acbs.com.vn/account-service/db"
	"acbs.com.vn/account-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var (
	migrator *db.Migrate
)

const (
	DBDriver = "pgx"
)

func TestMain(m *testing.M) {
	//cfg, err := config.InitConfig()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//appLogger := logger.NewAppLogger(cfg.Logger)
	//appLogger.InitLogger()
	//appLogger.WithName("Account-Service")

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	_ = resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	migrator = db.Migrator(db.WithDSN(databaseUrl))

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		migrator.DB, err = sql.Open(DBDriver, databaseUrl)
		if err != nil {
			return err
		}
		return migrator.DB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Performing a migration this way means all tests in this package shares
	// the same db schema across all unit test.
	// If isolation is needed, then do away with using `testing.M`. Do a
	// migration for each test handler instead.
	migrator.Up()

	// We can access database with m.hostAndPort or m.databaseUrl
	// port changes everytime a new docker instance is run
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestRepository_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *models.Account
	}
	type want struct {
		lastInsertID int
		nameInsert   string
		err          error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				ctx: context.Background(),
				req: &models.Account{
					Name: "name",
				},
			},
			want: want{
				lastInsertID: 1,
				nameInsert:   "name",
				err:          nil,
			},
		}, {
			name: "adding a second book should return ID=2",
			args: args{
				ctx: context.Background(),
				req: &models.Account{
					Name: "name2",
				},
			},
			want: want{
				lastInsertID: 2,
				nameInsert:   "name2",
				err:          nil,
			},
		},
		{
			name: "empty strings",
			args: args{
				ctx: context.Background(),
				req: &models.Account{
					Name: "",
				},
			},
			want: want{
				lastInsertID: 0,
				err:          errors.New("repository.Book.Create"),
			},
		},
	}

	client := sqlxDBClient(migrator.DB)
	repo := NewAccountRepository(nil, nil, client)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.Create(test.args.ctx, test.args.req)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.nameInsert, got.Name)
		})
	}
}

func sqlxDBClient(db *sql.DB) *sqlx.DB {
	return sqlx.NewDb(db, DBDriver)
}
