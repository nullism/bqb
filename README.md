# Basic Query Builder

[![Tests Status](https://github.com/nullism/bqb/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/nullism/bqb/actions/workflows/tests.yml) [![GoDoc](https://godoc.org/github.com/nullism/bqb?status.svg)](https://godoc.org/github.com/nullism/bqb) [![code coverage](coverage.svg)](https://github.com/nullism/bqb/actions/workflows/tests.yml) [![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/nullism/bqb) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

# Compatibility

This has been tested with sqlite, PostGres, and MySQL, using `database/sql`, `pq`, `pgx`, and `sqlx`.
By the nature of how it works it should be fully compatible with any DB interface and database that uses `?` or `$` parameter syntax.

# Why

1. Simple, lightweight, and fast
2. Supports any and all syntax by the nature of how it works
3. Doesn't require learning special syntax or operators
4. 100% test coverage

# Examples

## Basic

```golang
q := bqb.New("SELECT * FROM places WHERE id = ?", 1234)
sql, params, err := q.ToSql()
```

Produces

```sql
SELECT * FROM places WHERE id = ?
```

```
PARAMS: [1234]
```

## Postgres - ToPgsql()

Just call the `ToPgsql()` method instead of `ToSql()` to convert the query to Postgres syntax

```golang
q := bqb.New("DELETE FROM users").
    Space("WHERE id = ? OR name IN (?)", 7, []string{"delete", "remove"}).
    Space("LIMIT ?", 5)
sql, params, err := q.ToPgsql()
```

Produces

```sql
DELETE FROM users WHERE id = $1 OR name IN ($2, $3) LIMIT $4
```

```
PARAMS: [7, "delete", "remove", 5]
```

## Raw - ToRaw()

_Obvious warning: You should not use this for user input_

The `ToRaw()` call returns a string with the values filled in rather than parameterized

```golang
q := New("a = ?, b = ?, c = ?", "my a", 1234, nil)
sql, err := q.ToRaw()
```

Produces

```
a = 'my a', b = 1234, c = NULL
```

## Types

```golang
q := bqb.New(
    "int:? string:? []int:? []string:? Query:? JsonMap:? nil:? []intf:?",
    1, "2", []int{3, 3}, []string{"4", "4"}, bqb.New("5"), bqb.JsonMap{"6": 6}, nil, []interface{}{"a",1,true},
)
sql, _ := q.ToRaw()
```

Produces

```
int:1 string:'2' []int:3,3 []string:'4','4' Query:5 JsonMap:'{"6":6}' nil:NULL []intf:'a',1,true
```

### driver.Valuer

The [driver.Valuer](https://pkg.go.dev/database/sql/driver#Valuer) interface is supported for types that are able to convert
themselves to a sql driver value. See [examples/main.go:valuer](./examples/main.go#L102).

```
q := bqb.New("?", valuer)
```

### Embedder

BQB provides an Embedder interface for directly replacing `?` with a string returned by the `RawValue` method on the Embedder implementation.

This can be useful for changing sort direction or embedding table and column names. See [examples/main.go:embedder](./examples/main.go#L122) for an example.

_Note: Since this is a raw value, special attention should be paid to ensure user-input is checked and sanitized._

## Query IN

Arguments of type `[]string`,`[]*string`, `[]int`,`[]*int`, or `[]interface{}` are automatically expanded.

```golang
    q := bqb.New(
        "strs:(?) *strs:(?) ints:(?) *ints:(?) intfs:(?)",
        []string{"a", "b"}, []*string{}, []int{1, 2}, []*int{}, []interface{}{3, true},
    )
    sql, params, _ := q.ToSql()
```

Produces

```
SQL: strs:(?,?) *strs:(?) ints:(?,?) *ints:(?) intfs:(?,?)
PARAMS: [a b <nil> 1 2 <nil> 3 true]
```

## Json Arguments

There are two helper structs, `JsonMap` and `JsonList` to make JSON conversion a little simpler.

```golang

sql, err := bqb.New(
    "INSERT INTO my_table (json_map, json_list) VALUES (?, ?)",
    bqb.JsonMap{"a": 1, "b": []string{"a","b","c"}},
    bqb.JsonList{"string",1,true,nil},
).ToRaw()
```

Produces

```sql
INSERT INTO my_table (json_map, json_list)
VALUES ('{"a": 1, "b": ["a","b","c"]}', '["string",1,true,null]')
```

## Query Building

Since queries are built in an additive way by reference rather than value, it's easy to mutate a query without
having to reassign the result.

### Basic Example

```golang
sel := bqb.New("SELECT")

...

// later
sel.Space("id")

...

// even later
sel.Comma("age").Comma("email")
```

Produces

```sql
SELECT id,age,email
```

### Advanced Example

The `Optional(string)` function returns a query that resolves to an empty string if no query parts have
been added via methods on the query instance. For example `q := Optional("SELECT")` will resolve to
an empty string unless parts have been added by one of the methods,
e.g `q.Space("* FROM my_table")` would make `q.ToSql()` resolve to `SELECT * FROM my_table`.

```golang

sel := bqb.Optional("SELECT")

if getName {
    sel.Comma("name")
}

if getId {
    sel.Comma("id")
}

if !getName && !getId {
    sel.Comma("*")
}

from := bqb.New("FROM my_table")

where := bqb.Optional("WHERE")

if filterAdult {
    adultCond := bqb.New("name = ?", "adult")
    if ageCheck {
        adultCond.And("age > ?", 20)
    }
    where.And("(?)", adultCond)
}

if filterChild {
    where.Or("(name = ? AND age < ?)", "youth", 21)
}

q := bqb.New("? ? ?", sel, from, where).Space("LIMIT ?", 10)
```

Assuming all values are true, the query would look like:

```sql
SELECT name,id FROM my_table WHERE (name = 'adult' AND age > 20) OR (name = 'youth' AND age < 21) LIMIT 10
```

If `getName` and `getId` are false, the query would be

```sql
SELECT * FROM my_table WHERE (name = 'adult' AND age > 20) OR (name = 'youth' AND age < 21) LIMIT 10
```

If `filterAdult` is `false`, the query would be:

```sql
SELECT name,id FROM my_table WHERE (name = 'youth' AND age < 21) LIMIT 10
```

If all values are `false`, the query would be:

```sql
SELECT * FROM my_table LIMIT 10
```

## Methods

Methods on the bqb `Query` struct follow the same pattern.

All query-modifying methods take a string (the query text) and variable length interface (the query args).

For example `q.And("abc")` will add `AND abc` to the query.

Take the following

```golang
q := bqb.Optional("WHERE")
q.Empty() // returns true
q.Len() // returns 0
q.Space("1 = 2") // query is now WHERE 1 = 2
q.Empty() // returns false
q.Len() // returns 1
q.And("b") // query is now WHERE 1 = 2 AND b
q.Or("c") // query is now WHERE 1 = 2 AND b OR c
q.Concat("d") // query is now WHERE 1 = 2 AND b OR cd
q.Comma("e") // query is now WHERE 1 = 2 AND b OR cd,e
q.Join("+", "f") // query is now WHERE 1 = 2 AND b OR cd,e+f
```

Valid `args` include `string`, `int`, `floatN`, `*Query`, `[]int`, `Embedder`, `Embedded`, `driver.Valuer` or `[]string`.

# Frequently Asked Questions

## Is there more documentation?

It's not really necessary because the API is so tiny and public methods are documented in code,
see [godoc](https://godoc.org/github.com/nullism/bqb).
However, you can check out the [tests](./query_test.go) and [examples](./examples/main.go) to see the variety of usages.

## Why not just use a string builder?

Bqb provides several benefits over a string builder:

For example let's say we use the string builder way to build the following:

```golang
var params []interface{}
var whereParts []string
q := "SELECT * FROM my_table "
if filterAge {
    params = append(params, 21)
    whereParts = append(whereParts, fmt.Sprintf("age > $%d ", len(params)))
}

if filterBobs {
    params = append(params, "Bob%")
    whereParts = append(whereParts, fmt.Sprintf("name LIKE $%d ", len(params)))
}

if len(whereParts) > 0 {
    q += "WHERE " + strings.Join(whereParts, " AND ") + " "
}

if limit != nil {
    params = append(params, limit)
    q += fmt.Sprintf("LIMIT $%d", len(params))
}

// SELECT * FROM my_table WHERE age > $1 AND name LIKE $2 LIMIT $3
```

Some problems with that approach

1. You must perform a string join for the various parts of the where clause
2. You must remember to include a trailing or leading space for each clause
3. You have to keep track of parameter count (for Postgres anyway)
4. It's kind of ugly

The same logic can be achieved with `bqb` a bit more cleanly

```golang
q := bqb.New("SELECT * FROM my_table")
where := bqb.Optional("WHERE")
if filterAge {
    where.And("age > ?", 21)
}

if filterBobs {
    where.And("name LIKE ?", "Bob%")
}

q.Space("?", where)

if limit != nil {
    q.Space("LIMIT ?", limit)
}

// SELECT * FROM my_table WHERE age > $1 AND name LIKE $2 LIMIT $3
```

Both methods will allow you to remain close to the SQL, however the `bqb` approach will

1. Easily adapt to MySQL or Postgres without changing parameters
2. Hide the "WHERE" clause if both `filterBobs` and `filterAge` are false

## Why not use a full query builder?

Take the following _typical_ query example:

```golang
q := qb.Select("*").From("users").Where(qb.And{qb.Eq{"name": "Ed"}, qb.Gt{"age": 21}})
```

Vs the bqb way:

```golang
q := bqb.New("SELECT * FROM users WHERE name = ? AND age > ?", "ed", 21)
```

## Okay, so a simple query it might make sense to use something like `bqb`, but what about grouped queries?

A query builder can handle this in multiple ways, a fairly common pattern might be:

```golang
q := qb.Select("name").From("users")

and := qb.And{}

if checkAge {
    and = append(and, qb.Gt{"age": 21})
}

if checkName {
    or := qb.Or{qb.Eq{"name":"trusted"}}
    if nullNameOkay {
        or = append(or, qb.Is{"name": nil})
    }
    and = append(and, or)
}

q = q.Where(and)

// SELECT name FROM users WHERE age > 21 AND (name = 'trusted' OR name IS NULL)
```

Contrast that with the `bqb` approach:

```golang

q := bqb.New("SELECT name FROM users")

where := bqb.Optional("WHERE")

if checkAge {
    where.And("age > ?", 21)
}

if checkName {
    or := bqb.New("name = ?", "trusted")
    if nullNameOkay {
        or.Or("name IS ?", nil)
    }
    where.And("(?)", or)
}

q.Space("?", where)

// SELECT name FROM users WHERE age > 21 AND (name = 'trusted' OR name IS NULL)
```

It seems to be a matter of taste as to which method appears cleaner.
