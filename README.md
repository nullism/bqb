# bqb
Basic Query Builder

This project aims to provide a very lightweight and easy to use Query Builder
that provides an unescaped-first paradigm.

## Why?

* `bqb` does not require you to learn special syntax for operators.
* `bqb` makes `and`/`or` grouping function universally, and can be added to any clause.
* `bqb` is very small and quite fast.
* `bqb` is order-independent. Query components can be added in any order, which prevents side effects with complex query building logic.


## Examples

```golang
// Examples assume bqb has been imported
import "github.com/nullism/bqb"
```

Running examples:
```
$ go run examples/query/main.go
...
$ go run examples/insert/main.go
...
```

### Basic Select

```golang
q := bqb.Select("id, name, email").
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
q := bqb.Select("*").
    From("users").
    Where(
        bqb.And(
            bqb.V("email = ?", email),
            bqb.V("password = ?", password),
        ),
    ).Postgres()
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
bqb.Select("uuidv3_generate() as uuid", "u.id", "UPPER(u.name) as screamname", "u.age", "e.email").
    From("users u").
    Join("emails e ON e.user_id = u.id").
    JoinType(
        "LEFT OUTER JOIN",
        "friends f ON f.user_id = u.id",
    ).
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
    Limit(10).
    Postgres()
```

Produces

```sql
SELECT uuidv3_generate() as uuid, u.id, UPPER(u.name) as screamname, u.age, e.email
FROM users u
    JOIN emails e ON e.user_id = u.id
    LEFT OUTER JOIN friends f ON f.user_id = u.id
WHERE (
    (u.id IN ($1, $2, $3) AND e.email LIKE $4)
    OR
    (u.id IN ($5, $6, $7) AND e.email LIKE $8)
    OR
    u.id IN ($9, $10, $11, $12, $13, $14)
)
ORDER BY u.age DESC LIMIT 10
```
```
PARAMS: [1 3 5 %@gmail.com 2 4 6 %@hotmail.com 7 8 9 10 11 12]
```

### And / Or

And and Or can be used in any clause and not just `Where`. Subselects, `Having`, `OrderBy`, etc are all valid.
Separate clauses are assumed to be joined with `, ` without an `Or` or `And` call.


For example:

```golang
Select("*").From("patrons").
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
    ).Postgres()
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
q := bqb.Insert("my_table").
    Cols("name", "age", "current_time").
    Vals(bqb.V("?, ?, ?", "someone", 42, "2021-01-01 01:01:01Z")).
    Postgres()
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
q := bqb.Insert("my_table").
    Cols("name", "age", "current_time").
    Select(
        bqb.Select("b_name", "b_age", "b_time").
            From("b_table").
            Where(bqb.V("my_age > ?", 20)).
            Limit(10),
    ).Postgres()
```

Produces
```sql
INSERT INTO my_table (name, age, current_time)
SELECT b_name, b_age, b_time FROM b_table WHERE my_age > $1 LIMIT 10
```
```
PARAMS: [20]
```

### Basic Update

```golang
bqb.UpdatePsql().
    Update("my_table").
    Set(
        bqb.V("name = ?", "McCallister"),
        "age = 20", "current_time = CURRENT_TIMESTAMP()",
    ).
    Where(
        bqb.V("name = ?", "Mcallister"),
    )
```

Produces
```sql
UPDATE my_table SET name = $1, age = 20, current_time = CURRENT_TIMESTAMP()
WHERE name = $2
```
```
PARAMS: [McCallister Mcallister]
```

### Update with Sub-Queries

As with all clauses, sub-queries can be written inline as part of the clause or assigned
to a variable and used that way.

Note that the `Enclose()` method simply wraps the query in parentheses.

The `Concat()` method makes the following expressions join without any separator.

```golang
timeQ := bqb.QueryPsql().
    Select("timestamp").
    From("time_data").
    Where("is_current = true").
    Limit(1)

nameQ := bqb.QueryPsql().
    Select("name").
    From("users").
    Where(bqb.V("name LIKE ?", "%allister"))

bqb.Update("my_table").
    Set(
        bqb.V("name = ?", "McCallister"),
        "age = 20",
        bqb.Concat(
            "current_timestamp = ",
            timeQ.Enclose(),
        ),
    ).
    Where(
        bqb.Concat(
            "name IN ",
            nameQ.Enclose(),
        ),
    ).Postgres()
```

Produces
```sql
UPDATE my_table SET name = $1, age = 20, current_timestamp = (
    SELECT timestamp FROM time_data WHERE is_current = true LIMIT 1
) WHERE name IN (
    SELECT name FROM users WHERE name LIKE $2
)
```
```
PARAMS: [McCallister %allister]
```

### Using Text Queries

Sometimes it's easier to read inline queries than it
is to add another `bqb.Select()` to an existing query.
In these instances there's nothing wrong with writing an inline query as follows:

```golang
bqb.Select(
        "age as my_age",
        "(SELECT id FROM id_list WHERE user='me') as my_id",
    ).
...
```

### Create Table / Indexes

```golang
bqb.CreateTable("my_table").Cols("a VARCHAR(50) NOT NULL", "b BOOLEAN DEFAULT false")
```

Produces
```sql
CREATE TABLE my_table (
    a VARCHAR(50) NOT NULL,
    b BOOLEAN DEFAULT false
)
```

Subqueries are also possible in `CreateTable`

```golang
bqb.CreateTable("new_table").
    Cols("a INT NOT NULL DEFAULT 1", "b VARCHAR(50) NOT NULL").
    Select(
        bqb.Select("a", "b").From("other_table").Where("a IS NOT NULL"),
    )
```

Produces
```sql
CREATE TABLE new_table AS
SELECT a, b FROM other_table WHERE a IS NOT NULL
```

### Escaping `?`

Just use `??` instead of `?` in the query, for example:

```golang
Select("data->>'id' ?? '1234'") ...
```