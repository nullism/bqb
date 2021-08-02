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

func TestQueryIn(t *testing.T) {
	q := Select("a, b").From("table").Where(V("a IN (?)", []int{1, 2, 3}))
	sql, params, _ := q.ToSql()
	want := "SELECT a, b FROM table WHERE a IN (?, ?, ?)"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
	if len(params) != 3 {
		t.Errorf("expected 3 params, got: %v", len(params))
	}
	if params[0] != 1 || params[1] != 2 || params[2] != 3 {
		t.Errorf("got invalid params: %v", params)
	}
}

func TestQueryPrint(t *testing.T) {
	q := Select("*").From("table").Where(V("field = ?", 1))
	q.Print() // @TODO: capture standard out?
}

func TestQueryValues(t *testing.T) {
	q := Select("a.a", "b.a").From("a_table a").Join("b_table b ON b.a = a.a").
		Where(
			And(
				V("b.a IS NOT ?", nil),
				V("a.a > ?", 10),
				V("a.name LIKE ?", "test%"),
			),
		).Limit(10).Postgres()
	sql, params, _ := q.ToSql()
	if len(params) != 3 {
		t.Errorf("expected 3 params, got: %v", len(params))
	}

	if params[0] != nil || params[1] != 10 || params[2] != "test%" {
		t.Errorf("invalid params: %v", params)
	}

	want := "SELECT a.a, b.a FROM a_table a JOIN b_table b ON b.a = a.a WHERE (b.a IS NOT $1 AND a.a > $2 AND a.name LIKE $3) LIMIT 10"
	if want != sql {
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
