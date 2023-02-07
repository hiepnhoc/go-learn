package postgres

import (
	"acbs.com.vn/account-service/utility/database"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func NewSqlx(cfg *Config) (*sqlx.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	database.Alive(db.DB)

	return db, nil
}
