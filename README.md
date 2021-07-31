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
    Select("id", "name", "email").
    From("users").
    Where("email LIKE '%@yahoo.com'")
```

Produces

```sql
SELECT id, name, email FROM users WHERE (email LIKE $1)
```
```
PARAMS: [%@yahoo.com]
```

### Select With Join

```golang
q := bqb.New(bqb.PGSQL).
    Select("uuidv3_generate() as uuid", "u.id", "UPPER(u.name) as screamname", "u.age", "e.email").
    From("users u").
    Join("emails e ON e.user_id = u.id").
    Where(
        bqb.Valf("u.id IN (?, ?, ?)", 1, 3, 5),
        bqb.Valf("e.email LIKE ?", "%@gmail.com"),
    ).
    Where(
        bqb.Valf("u.id IN (?, ?, ?)", 2, 4, 6),
        bqb.Valf("e.email LIKE ?", "%@hotmail.com"),
    ).
    Where(
        // you can pass arrays for IN / ANY clauses
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

And queries are applied to everything within the same `Where`, and Or queries
are assumed for consequtive `Where` clauses.

For example:

```golang
bqb.New(bqb.PGSQL).Select("*").From("patrons").
Where(
    "drivers_license IS NOT NULL",
    "age > 20",
    "age < 60",
).
Where(
    "drivers_license IS NULL",
    "age > 60",
).
Where(
    "is_known = true",
)
```

Produces

```sql
SELECT * FROM patrons WHERE
    (drivers_license IS NOT NULL AND age > 20 AND age < 60)
    OR
    (drivers_license IS NULL AND age > 60)
    OR
    (is_known = true)
```

### Escaping `?`

Just use `??` instead of `?` in the query, for example:

```golang
Select("data->>'id' ?? '1234'") ...
```