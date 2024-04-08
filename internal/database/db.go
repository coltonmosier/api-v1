package database

import (
	"database/sql"
	"log"
	"os"

    "github.com/go-sql-driver/mysql"
)

func InitEquipmentDatabase() *sql.DB {
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_EQUIPTMENT_DB"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func InitLoggingDatabase() *sql.DB {
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_EQUIPTMENT_DB"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db
}
