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

}