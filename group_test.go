package bqb

import (
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	q := Group(
		Select("*").From("a"),
		"UNION ALL",
		Select("*").From("b"),
	)
	sql, _, _ := q.ToSql()
	want := "SELECT * FROM a UNION ALL SELECT * FROM b"
	if want != sql {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func ExampleGroup() {
	q := Group(
		"CREATE TABLE my_table (",
		"name VARCHAR(50) NOT NULL",
		"age INTEGER NOT NULL DEFAULT 18",
		"is_authorized BOOLEAN NOT NULL DEFAULT false",
		")",
	)
	sql, _, _ := q.ToSql()
	fmt.Println(sql)
	// Output: CREATE TABLE my_table ( name VARCHAR(50) NOT NULL age INTEGER NOT NULL DEFAULT 18 is_authorized BOOLEAN NOT NULL DEFAULT false )
}
