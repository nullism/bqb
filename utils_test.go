package bqb

import (
	"testing"
)

func Test_dialectReplace_unknown_dialect(t *testing.T) {
	const (
		testSql = "test-sql"
	)
	params := []any{1, 2, "a", "c"}
	sql, err := dialectReplace(Dialect("unknown"), testSql, params)

	if sql != "test-sql" {
		t.Errorf("unexpected sql statement: want %s got %s", testSql, sql)
	}

	if err != nil {
		t.Error("unknown dialect should not return an error")
	}
}
