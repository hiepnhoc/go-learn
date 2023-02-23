package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type PostgresConnector struct {
}

func (p *PostgresConnector) GetConnection() (db *gorm.DB, err error) {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort, err := strconv.ParseInt(os.Getenv("db_port"), 10, 0)

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable port=%d password=%s", dbHost, username, dbName, dbPort, password)

	fmt.Println(dbURI)
	return gorm.Open("postgres", dbURI)
}
