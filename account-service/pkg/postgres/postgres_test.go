package postgres

import (
	testPostgres "2margin.vn/account-service/test/containers/postgres"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"testing"
)

var container *testPostgres.Container

func TestMain(m *testing.M) {

	container = testPostgres.NewTestContainer(m)
	err := container.Run()
	if err != nil {
		log.Fatalf("TestMain faild %v", err)
	}

}

func TestNewSqlxDB(t *testing.T) {

	if container.Cfg == nil {
		t.Fatalf("NewSqlxDB Cfg is missing setup")
	}
	_, err := NewSqlxDB(container.Cfg)
	if err != nil {
		t.Fatalf("NewSqlxDB Error %v", err)
	}

}
