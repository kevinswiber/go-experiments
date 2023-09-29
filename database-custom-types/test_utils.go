package dbtypes

import (
	"database/sql"
	"io"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

type TestDb struct {
	db   *sql.DB
	file string
}

func (td *TestDb) Close() {
	td.db.Close()
	os.Remove(td.file)
}

func InitDB(t *testing.T) *TestDb {
	dest, err := os.CreateTemp("", "test_db.*.sqlite")
	if err != nil {
		t.Error(err)
		return nil
	}

	result := &TestDb{
		file: dest.Name(),
	}

	src, err := os.Open("db.sqlite.fixture")
	if err != nil {
		t.Error(err)
		return nil
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		t.Error(err)
		dest.Close()
		return nil
	}
	dest.Close()

	db, err := sql.Open("sqlite", result.file)
	if err != nil {
		t.Error(err)
		return nil
	}

	result.db = db
	return result
}
