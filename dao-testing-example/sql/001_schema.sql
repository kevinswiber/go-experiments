create table people (
  id integer primary key,
  name text,
  email text,
  photo text
);
create unique index people_email_idx on people (email);
