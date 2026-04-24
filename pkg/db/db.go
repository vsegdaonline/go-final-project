package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(64) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);`

func GetDBFile() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	return dbFile
}

var db *sql.DB

func Init(dbFile string) error {
	var install bool

	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		install = true
	} else if err != nil {
		return err
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	if install {
		_, err = db.Exec(schema)
	}
	return err
}
