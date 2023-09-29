This experiment is used to test serializing and deserializing types to
databases. There is an sqlite fixture database included that contains a
single table:

```sql
create table foo (
    id integer primary key,
    data blob
)
```
