package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"embed"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed resources
var dbScript embed.FS

type LatLng []float64

func (ll LatLng) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("(")
	buf.WriteString(strconv.FormatFloat(ll[0], 'f', -1, 64))
	buf.WriteString(",")
	buf.WriteString(strconv.FormatFloat(ll[1], 'f', -1, 64))
	buf.WriteString(")")
	return buf.String()
}

func (ll LatLng) Value() (driver.Value, error) {
	return ll.String(), nil
}

func (ll *LatLng) Scan(value any) error {
	if value == nil {
		return nil
	}

	rv := reflect.TypeOf(value)
	if rv.Name() != "string" {
		return fmt.Errorf("value must be a string, got: %s", rv.Name())
	}

	str, err := driver.String.ConvertValue(value)
	if err != nil {
		return err
	}
	if str == "" {
		return nil
	}

	parts := strings.Split(strings.Trim(str.(string), "()"), ",")
	if len(parts) != 2 {
		return fmt.Errorf("expected (float64,float64) but received: %s", str)
	}

	val1, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}

	val2, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}

	*ll = LatLng{val1, val2}
	return nil
}

const insertSql = `
insert into foo (id, data)
values (@id, @data)
on conflict do update 
set data = @data
`

const selectSql = `
select data from foo where id = 1
`

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := initDb()
	if err != nil {
		return err
	}
	defer db.Close()

	latLng := LatLng{42.12345, -42.54321}

	_, err = db.Exec(insertSql, sql.Named("id", 1), sql.Named("data", latLng))
	if err != nil {
		return err
	}

	rows, err := db.Query(selectSql)
	if err != nil {
		return err
	}

	for rows.Next() {
		var dbLatLng LatLng
		err = rows.Scan(&dbLatLng)
		fmt.Printf("dbLatLng = %+v\n", dbLatLng)
	}

	return nil
}

func initDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./db.sqlite")
	if err != nil {
		return nil, err
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		db.Close()
		return nil, err
	}

	fsDriver, err := iofs.New(dbScript, "resources")
	if err != nil {
		db.Close()
		return nil, err
	}

	migrator, err := migrate.NewWithInstance("iofs", fsDriver, "sqlite", driver)
	if err != nil {
		db.Close()
		return nil, err
	}
	err = migrator.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return db, nil
		}
		db.Close()
		return nil, err
	}

	return db, nil
}
