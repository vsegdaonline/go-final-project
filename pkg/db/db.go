package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(64) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);`

func Init(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	defer func() {
		err = db.Close()
	}()
	if err != nil {
		return err
	}
	err = db.Ping()
	return err
}
