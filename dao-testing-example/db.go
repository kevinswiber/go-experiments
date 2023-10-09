package dte

import (
	"database/sql"
	"embed"

	_ "modernc.org/sqlite"
)

//go:embed sql/*.sql
var ddlFiles embed.FS

type PeopleDb struct {
	conn *sql.DB
}

func NewDB(dbPath string) (*PeopleDb, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	schema, err := ddlFiles.ReadFile("sql/001_schema.sql")
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		db.Close()
		return nil, err
	}

	return &PeopleDb{conn: db}, nil
}

func (db *PeopleDb) Close() {
	db.conn.Close()
}
