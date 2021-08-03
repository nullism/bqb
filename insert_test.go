package bqb

import "testing"

func TestInsert(t *testing.T) {
	q := Insert("table").Cols("a", "b").Vals("1", "2")
	sql, _, _ := q.ToSql()
	want := "INSERT INTO table (a, b) VALUES (1, 2)"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestInsertSelect(t *testing.T) {
	q := Insert("table").Cols("a", "b").Select(
		Select("a, b").From("table_b").Where(
			And("a IS NOT NULL", "b IS NOT NULL"),
		),
	)
	sql, _, _ := q.ToSql()
	want := "INSERT INTO table (a, b) SELECT a, b FROM table_b WHERE (a IS NOT NULL AND b IS NOT NULL)"
	if want != sql {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestInsertOnDuplicateKey(t *testing.T) {
	q := Insert("table").Cols("a", "b").Vals("1", "2").OnDuplicateKey("a = a + 1")
	sql, _, _ := q.ToSql()
	want := "INSERT INTO table (a, b) VALUES (1, 2) ON DUPLICATE KEY UPDATE a = a + 1"
	if want != sql {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}
