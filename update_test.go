package bqb

import "testing"

func TestUpdate(t *testing.T) {
	q := Update("hello")
	sql, _, _ := q.ToSql()
	good := "UPDATE hello"
	if sql != good {
		t.Errorf("Update = %v, want %v", sql, good)
	}
}

func TestUpdateAdvanced(t *testing.T) {
	q := Update("table").Set("a = 1", "b = 2").Where("c = 3")
	sql, _, _ := q.ToSql()
	want := "UPDATE table SET a = 1, b = 2 WHERE c = 3"
	if sql != want {
		t.Errorf("Update = %q, want %q", sql, want)
	}
}
