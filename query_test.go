package bqb

import "testing"

func TestQuery(t *testing.T) {
	q := Select("*").From("table").Where("field IS NULL")
	sql, _, _ := q.ToSql()
	want := "SELECT * FROM table WHERE field IS NULL"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestSubQuery(t *testing.T) {
	q := Select(Select("1").Enclose()).From("table")
	sql, _, _ := q.ToSql()
	want := "SELECT (SELECT 1) FROM table"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}
