package main

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/nullism/bqb"
)

func basic() {
	println("===[ Example: Basic ]===")
	q := bqb.New("SELECT * FROM places WHERE id = ?", 1234)
	q.Print()
}

func builder() {
	println("===[ Example: Builder ]===")

	// Changing these values changes the output of the query
	getId := true
	getName := true
	lim := 10
	filterName := true
	filterNameNotNull := true
	filterAge := true

	// Optional queries will return nothing unless they have at least one query part.
	sel := bqb.Optional("SELECT")
	if getId {
		sel.Comma("id")
	}
	if getName {
		sel.Comma("name")
	}

	if sel.Len() > 0 {
		sel.Space("FROM my_table")
	}

	where := bqb.Optional("WHERE")

	if filterName {
		nameQ := bqb.New("name = ?", "name")
		if filterNameNotNull {
			nameQ.And("name IS NOT NULL")
		}
		where.And("(?)", nameQ)
	}

	if filterAge {
		where.And("age > ?", 21)
	}

	limit := bqb.Optional("LIMIT")

	if lim > 0 {
		limit.Space("?", lim)
	}

	query := bqb.New("? ? ?", sel, where, limit)

	query.Print()
}

func customTypes() {
	println("===[ Custom Types ]===")
	type myStruct struct {
		val string
	}
	q := bqb.New(
		"DELETE FROM my_table WHERE a = ?",
		&myStruct{val: "hello"},
	)
	q.Print()
}

func json() {
	println("===[ JSON ]===")
	q := bqb.New(
		"INSERT INTO my_table (json_map, json_list) VALUES (?, ?)",
		&bqb.JsonMap{"a": 1, "b": []string{"b", "c", "d"}},
		&bqb.JsonList{"a", 1, true, nil, &bqb.JsonMap{"a": "b"}},
	)
	q.Print()
}

func nilQuery() {
	println("===[ NIL Query ]===")
	var q *bqb.Query
	q.Print()
}

type valuerExample struct {
	Name string
	Age  int
}

func (v valuerExample) Value() (driver.Value, error) {
	return fmt.Sprintf("%v is %v years old", v.Name, v.Age), nil
}

func valuer() {
	println("===[ driver.Valuer ]===")
	example := &valuerExample{Name: "Bobby Tables", Age: 12}
	q := bqb.New("info text = ?", example)
	q.Print()
}

type embedderExample []string

func (e embedderExample) RawValue() string {
	parts := []string{}
	for _, v := range e {
		part := fmt.Sprintf(`"%v"`, strings.ReplaceAll(v, `"`, `""`))
		parts = append(parts, part)
	}
	return strings.Join(parts, ".")
}

func embedder() {
	println("===[ Embedder ]===")
	emb := embedderExample{"one", "two", `"three"`}
	q := bqb.New("embedded: ?, unembedded: ?", emb, "bound")
	q.Print()
}

func main() {
	customTypes()
	basic()
	builder()
	json()
	nilQuery()
	valuer()
	embedder()
}
