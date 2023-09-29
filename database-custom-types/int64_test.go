package dbtypes

import (
	"database/sql"
	"testing"
)

func Test_int64Value(t *testing.T) {
	testDb := InitDB(t)
	if testDb == nil {
		return
	}
	db := testDb.db

	t.Cleanup(func() {
		testDb.Close()
	})

	t.Run("writes and reads a root value to the database", func(t *testing.T) {
		var val RootLong = 42

		db.Exec("insert into foo (id, data) values (1, @data)", sql.Named("data", val))

		rows, err := db.Query("select data from foo where id = 1")
		if err != nil {
			t.Error(err)
			return
		}

		for rows.Next() {
			var readVal RootLong
			err := rows.Scan(&readVal)
			if err != nil {
				t.Error(err)
				break
			}

			if readVal != 42 {
				t.Errorf("expected 42, to %q", readVal)
				break
			}
		}
	})

	t.Run("writes and reads a child value to the database", func(t *testing.T) {
		var val ChildLong = 42

		db.Exec("insert into foo (id, data) values (2, @data)", sql.Named("data", val))

		rows, err := db.Query("select data from foo where id = 2")
		if err != nil {
			t.Error(err)
			return
		}

		for rows.Next() {
			var readVal ChildLong
			err := rows.Scan(&readVal)
			if err != nil {
				t.Error(err)
				break
			}

			if readVal != 42 {
				t.Errorf("expected 42, to %q", readVal)
				break
			}
		}
	})
}
