package bqb

import "testing"

func TestAnd(t *testing.T) {
	expr := And("1", "2")
	want := "(1 AND 2)"
	if want != expr.F {
		t.Errorf("want: %q, got: %q", want, expr.F)
	}
}

func TestConcat(t *testing.T) {
	expr := Concat("1", "2", "3")
	want := "123"
	if want != expr.F {
		t.Errorf("want: %q, got: %q", want, expr.F)
	}
}
