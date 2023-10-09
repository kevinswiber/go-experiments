package dte

import (
	"database/sql"
)

const insertPersonSql = `
insert into people (id, name, email, photo) values (
	@id,
	@name,
	@email,
	@photo
)
on conflict do update set
	name = @name,
	email = @email,
	photo = @photo
`

const getPersonSql = `
select
  id,
  name,
  email,
  photo
from people
where id = @id
`

type Person struct {
	Id    int
	Name  string
	Email string
	Photo string
}

func (db *PeopleDb) AddPerson(person *Person) error {
	_, err := db.conn.Exec(
		insertPersonSql,
		sql.Named("id", person.Id),
		sql.Named("name", person.Name),
		sql.Named("email", person.Email),
		sql.Named("photo", person.Photo),
	)

	return err
}

func (db *PeopleDb) GetPerson(id int) (*Person, error) {
	rows, err := db.conn.Query(getPersonSql, sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() == false {
		return nil, sql.ErrNoRows
	}

	var person Person
	err = rows.Scan(
		&person.Id,
		&person.Name,
		&person.Email,
		&person.Photo,
	)
	if err != nil {
		return nil, err
	}

	return &person, nil
}
