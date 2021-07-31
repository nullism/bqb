# bqb
Basic Query Builder

This project aims to provide a very lightweight and easy to use Query Builder
that provides an unescaped-first paradigm.

## Why?

* `bqb` does not require you to learn special syntax for operators
* `bqb` makes `and`/`or` grouping simple to understand
* `bqb` is very small, and quite fast


## Examples

```golang
// Examples assume bqb has been imported
import "github.com/nullism/bqb"
```

### Basic Select

```golang
q := bqb.New(bqb.PGSQL).
    Select("id, name, email").
    From("users").
    Where("email LIKE '%@yahoo.com'")
sql, params, err := q.ToSql()
```

Produces

```sql
SELECT id, name, email FROM users WHERE (email LIKE '%@yahoo.com')
```

### Bind Variables

Often times user-provided information is used to generate queries.
In this case, you should _always_ wrap those values in the `Valf()` function.

```golang
email := "foo@bar.com"
password := "p4ssw0rd"
q := bqb.New(bqb.PGSQL).
    Select("*").
    From("users").
    Where(
        bqb.Valf("email = ?", email),
        bqb.Valf("password = ?", password),
    )
```

Produces
```sql
SELECT * FROM users WHERE (email = $1 AND password = $2)
```
```
PARAMS: [foo@bar.com p4ssw0rd]
```

### Select With Join

```golang
q := bqb.New(bqb.PGSQL).
    Select("uuidv3_generate() as uuid", "u.id", "UPPER(u.name) as screamname", "u.age", "e.email").
    From("users u").
    Join("emails e ON e.user_id = u.id").
    Where(
        bqb.Valf("u.id IN (?, ?, ?) AND e.email LIKE ?", 1, 3, 5, "%@gmail.com"),
        bqb.Valf("u.id IN (?, ?, ?) AND e.email LIKE ?", 2, 4, 6, "%@hotmail.com"),
        bqb.Valf("u.id IN (?)", []int{7, 8, 9, 10, 11, 12}),
    ).
    OrderBy("u.age DESC").
    Limit(10)

```

Produces

```sql
SELECT uuidv3_generate() as uuid, u.id, UPPER(u.name) as screamname, u.age, e.email
FROM users u JOIN emails e ON e.user_id = u.id
WHERE (u.id IN ($1, $2, $3) AND e.email LIKE $4)
OR (u.id IN ($5, $6, $7) AND e.email LIKE $8)
OR (u.id IN ($9, $10, $11, $12, $13, $14))
ORDER BY u.age DESC LIMIT 10
```
```
PARAMS: [1 3 5 %@gmail.com 2 4 6 %@hotmail.com 7 8 9 10 11 12]
```

### And / Or

And are or have been simplified to a great deal in `bqb`.
Separate clauses are assumed to be joined with `OR`.


For example:

```golang
bqb.New(bqb.PGSQL).Select("*").From("patrons").
    Where(
        "drivers_license IS NOT NULL AND (age > 20 AND age < 60)",
        "drivers_license IS NULL AND age >= 60",
        "is_known = true",
    )
```

Produces

```sql
SELECT * FROM patrons WHERE
    (drivers_license IS NOT NULL AND age > 20 AND age < 60)
    OR
    (drivers_license IS NULL AND age >= 60)
    OR
    (is_known = true)
```

### Escaping `?`

Just use `??` instead of `?` in the query, for example:

```golang
Select("data->>'id' ?? '1234'") ...
```