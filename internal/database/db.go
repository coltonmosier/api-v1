package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/coltonmosier/api-v1/internal/sqlc"
	"github.com/go-sql-driver/mysql"
)

var (
	equeries *sqlc.Queries
	lqueries *sqlc.Queries
)

func InitEquipmentDatabase() (*sqlc.Queries, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_EQUIPMENT_DB"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error pinging database")
		return nil, err
	}

	equeries = sqlc.New(db)

	return equeries, nil
}

func InitLoggingDatabase() (*sqlc.Queries, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_LOG_DB"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error pinging database")
		return nil, err
	}

	lqueries = sqlc.New(db)

	return lqueries, nil
}
