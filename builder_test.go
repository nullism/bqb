package bqb

import (
	"fmt"
	"strings"
	"testing"
)

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

func Test_dialectReplace(t *testing.T) {
	mp := map[string][]interface{}{
		fmt.Sprintf("func(%v, %v)", paramPh, paramPh): {1, 2},
		fmt.Sprintf("%v", paramPh):                    {2},
		fmt.Sprintf("field = %v", paramPh):            {"some string"},
		fmt.Sprintf("val IS NOT %v", paramPh):         nil,
	}
	for text, val := range mp {
		newText := ""
		for _, dialect := range []string{PGSQL, MYSQL, SQL, RAW} {
			newText = dialectReplace(dialect, text, val)
			want := text
			for i, param := range val {
				switch d := dialect; d {
				case PGSQL:
					want = strings.Replace(want, paramPh, fmt.Sprintf("$%d", i+1), 1)
				case RAW:
					want = strings.Replace(want, paramPh, paramToRaw(param), 1)
				case SQL, MYSQL:
					want = strings.Replace(want, paramPh, "?", 1)
				default:
					t.Errorf("unknown dialect %v", d)
				}
			}
			if newText != want {
				t.Errorf("SQL %v not working: want: %q, got: %q", dialect, want, newText)
			}
		}
	}
}

func Test_getExprs(t *testing.T) {
	exprs := getExprs([]interface{}{"a", "b"})
	if exprs[0].F != "a" || exprs[1].F != "b" {
		t.Errorf("want: ['a', 'b'], got: %s, %s", exprs[0].F, exprs[1].F)
	}

}

func Test_intfToExpr(t *testing.T) {
	var want string

	expr := intfToExpr(1)
	if expr.F != "?" && expr.V[0] != 1 {
		t.Errorf("want: '?', got %q", expr.F)
	}

	expr = intfToExpr("abc")
	if expr.F != "abc" {
		t.Errorf("want: 'abc', got %q", expr.F)
	}

	expr = intfToExpr(Select("*").Where(V("a = ?", 1)))
	want = fmt.Sprintf("SELECT * WHERE a = %v", paramPh)
	if expr.F != want {
		t.Errorf("want: %q, got: %q", want, expr.F)
	}

	// Bad Values

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the expressions did not panic")
		}
	}()

	expr = intfToExpr([]int{1, 2, 3})
	expr = intfToExpr([]string{"a", "b"})
}

func Test_paramToRaw(t *testing.T) {
	want := "'string'"
	got := paramToRaw("string")
	if want != got {
		t.Errorf("want: %q, got: %q", want, got)
	}
}
