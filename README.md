# bqb
Basic Query Builder

This project aims to provide a very lightweight and easy to use Query Builder
that provides an unescaped-first paradigm.

## Why?

* `bqb` does not require you to learn special syntax for operators.
* `bqb` makes `and`/`or` grouping function universally, and can be added to any clause.
* `bqb` is very small and quite fast.


## Examples

```golang
// Examples assume bqb has been imported
import "github.com/nullism/bqb"
```

### Basic Select

```golang
q := bqb.QueryPsql().
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
In this case, you should _always_ wrap those values in the `V()` function.

```golang
email := "foo@bar.com"
password := "p4ssw0rd"
q := bqb.QueryPsql().
    Select("*").
    From("users").
    Where(
        bqb.And(
            bqb.V("email = ?", email),
            bqb.V("password = ?", password),
        ),
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
q := bqb.QueryPsql().
    Select("uuidv3_generate() as uuid", "u.id", "UPPER(u.name) as screamname", "u.age", "e.email").
    From("users u").
    Join("emails e ON e.user_id = u.id").
    Where(
        bqb.Or(
            bqb.And(
                bqb.V("u.id IN (?, ?, ?)", 1, 3, 5),
                bqb.V("e.email LIKE ?", "%@gmail.com"),
            ),
            bqb.And(
                bqb.V("u.id IN (?, ?, ?)", 2, 4, 6),
                bqb.V("e.email LIKE ?", "%@yahoo.com"),
            ),
            bqb.V("u.id IN (?)", []int{7, 8, 9, 10, 11, 12}),
        ),
    ).
    OrderBy("u.age DESC").
    Limit(10)
```

Produces

```sql
SELECT uuidv3_generate() as uuid, u.id, UPPER(u.name) as screamname, u.age, e.email
FROM users u
JOIN emails e ON e.user_id = u.id
WHERE (
    (u.id IN ($1, $2, $3) AND e.email LIKE $4)
    OR
    (u.id IN ($5, $6, $7) AND e.email LIKE $8)
    OR
    u.id IN ($9, $10, $11, $12, $13, $14)
) ORDER BY u.age DESC LIMIT 10
```
```
PARAMS: [1 3 5 %@gmail.com 2 4 6 %@hotmail.com 7 8 9 10 11 12]
```

### And / Or

And and Or can be used in any clause and not just `Where`. Subselects, `Having`, `OrderBy`, etc are all valid.
Separate clauses are assumed to be joined with `, ` without an `Or` or `And` call.


For example:

```golang
bqb.QueryPsql().Select("*").From("patrons").
    Where(
        bqb.Or(
            bqb.And(
                "drivers_license IS NOT NULL",
                bqb.And("age > 20", "age < 60"),
            ),
            bqb.And(
                "drivers_license IS NULL",
                "age >= 60",
            ),
            "is_known = true",
        ),
    )
```

Produces

```sql
SELECT * FROM patrons WHERE (
    (
        drivers_license IS NOT NULL AND (
            age > 20 AND age < 60
        )
    )
    ) OR (
    drivers_license IS NULL AND age >= 60
    ) OR is_known = true
)
```

### Basic Insert

```golang
q := bqb.InsertPsql().
    Into("my_table").
    Cols("name", "age", "current_time").
    Vals(bqb.V("?, ?, ?", "someone", 42, "2021-01-01 01:01:01Z"))
```

Produces
```sql
INSERT INTO my_table (name, age, current_time) ($1, $2, $3)
```
```
PARAMS: [someone 42 2021-01-01 01:01:01Z]
```


### Insert .. Select

```golang
q := bqb.InsertPsql().
    Into("my_table").
    Cols("name", "age", "current_time").
    Select(
        bqb.QueryPsql().
            Select("b_name", "b_age", "b_time").
            From("b_table").
            Where(bqb.V("my_age > ?", 20)).
            Limit(10),
    )
```

Produces
```sql
INSERT INTO my_table (name, age, current_time) SELECT b_name, b_age, b_time FROM b_table WHERE my_age > $1 LIMIT 10
```
```
PARAMS: [20]
```

### Escaping `?`

Just use `??` instead of `?` in the query, for example:

```golang
Select("data->>'id' ?? '1234'") ...
```