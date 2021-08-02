package bqb

import "testing"

func TestCreateTable(t *testing.T) {
	q := CreateTable("new_table").Cols("a INT NOT NULL DEFAULT 1", "b VARCHAR(50) NOT NULL")
	sql, _, _ := q.ToSql()
	want := "CREATE TABLE new_table (a INT NOT NULL DEFAULT 1, b VARCHAR(50) NOT NULL)"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestCreateTableFrom(t *testing.T) {
	q := CreateTable("new_table").
		Cols("a INT NOT NULL DEFAULT 1", "b VARCHAR(50) NOT NULL").
		Select(
			Select("a", "b").From("other_table").Where("a IS NOT NULL"),
		)
	sql, _, _ := q.ToSql()
	want := "CREATE TABLE new_table (a INT NOT NULL DEFAULT 1, b VARCHAR(50) NOT NULL) " +
		"AS SELECT a, b FROM other_table WHERE a IS NOT NULL"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}
