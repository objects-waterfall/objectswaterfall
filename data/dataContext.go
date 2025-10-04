package data

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DataContext struct {
	Db     *sql.DB
	Driver string
}

var DbContext DataContext

func InitDbConnection() error {
	var err error
	driver := os.Getenv("DB_DRIVER")
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return err
	}

	DbContext = DataContext{
		Db:     db,
		Driver: driver,
	}

	DbContext.Db.SetMaxOpenConns(100)
	DbContext.Db.SetConnMaxIdleTime(10)
	return nil
}
