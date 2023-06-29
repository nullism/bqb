package bqb

import (
	"reflect"
	"strings"
	"testing"
)

type embedder []string

func (e embedder) RawValue() string {
	return strings.Join(e, ".")
}

type sortEmbedder string

const (
	down sortEmbedder = "down"
	up   sortEmbedder = "up"
)

func (s sortEmbedder) RawValue() string {
	if s == down {
		return "DESC"
	}
	if s == up {
		return "ASC"
	}
	panic("invalid sort direction: " + s)
}

func TestEmbedder(t *testing.T) {
	emb := embedder{"one", "two", "three"}
	want := "one.two.three"

	if emb.RawValue() != want {
		t.Errorf("Embedder error: want=%v got=%v", want, emb.RawValue())
	}

	q := New("SELECT ? FROM ? WHERE ?=?", embedder{"id"}, embedder{"schema", "table"}, embedder{"name"}, "bound")
	sql, args, err := q.ToSql()

	if err != nil {
		t.Errorf("got error: %v", err)
	}

	want = "SELECT id FROM schema.table WHERE name=?"
	if want != sql {
		t.Errorf("\n got:%v\nwant:%v", sql, want)
	}

	wantArgs := []any{"bound"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("\n got:%v\nwant:%v", args, wantArgs)
	}

	sortq := New("SELECT * FROM my_table ORDER BY name ?,?", down, up)
	want = "SELECT * FROM my_table ORDER BY name DESC,ASC"
	got, _ := sortq.ToRaw()
	if got != want {
		t.Errorf("\n got:%v\nwant:%v", got, want)
	}

	sql, err = q.ToRaw()
	if err != nil {
		t.Errorf("got error from ToRaw(): %v", err)
	}
	want = "SELECT id FROM schema.table WHERE name='bound'"
	if sql != want {
		t.Errorf("\n got:%v\nwant:%v", sql, want)
	}
}

func TestEmbedded(t *testing.T) {
	q := New("SELECT ? FROM ? WHERE ?=?", Embedded("col"), Embedded("table"), Embedded("id"), "123")
	sql, args, err := q.ToSql()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want := "SELECT col FROM table WHERE id=?"
	if want != sql {
		t.Errorf("\n got:%v\nwant:%v", sql, want)
	}

	if args[0] != "123" {
		t.Errorf("got unexpected args: %v", args)
	}

	sql, err = q.ToRaw()
	if err != nil {
		t.Errorf("got error from ToRaw(): %v", err)
	}
	want = "SELECT col FROM table WHERE id='123'"
	if sql != want {
		t.Errorf("\n got:%v\nwant:%v", sql, want)
	}

}
