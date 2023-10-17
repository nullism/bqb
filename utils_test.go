package bqb

import (
	"fmt"
	"testing"
)

func Test_scanReplace(t *testing.T) {
	const (
		testReplace = "{{TEST_TOKEN}}"
	)
	type args struct {
		stmt    string
		replace string
		fn      replaceFn
	}
	type want struct {
		str string
		err error
	}

	for _, tt := range []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty statement",
			args: args{
				stmt:    "",
				replace: testReplace,
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "",
				err: nil,
			},
		},
		{
			name: "empty token statement",
			args: args{
				stmt:    "this tests an empty pattern token",
				replace: "",
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "this tests an empty pattern token",
				err: nil,
			},
		},
		{
			name: "no tokens",
			args: args{
				stmt:    "this tests no tokens",
				replace: testReplace,
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "this tests no tokens",
				err: nil,
			},
		},
		{
			name: "one front token",
			args: args{
				stmt:    fmt.Sprintf("%s this tests one token", testReplace),
				replace: testReplace,
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "0 this tests one token",
				err: nil,
			},
		},
		{
			name: "one boundary token",
			args: args{
				stmt:    fmt.Sprintf("this tests one%s token", testReplace),
				replace: testReplace,
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "this tests one0 token",
				err: nil,
			},
		},
		{
			name: "one end token",
			args: args{
				stmt:    fmt.Sprintf("this tests one token%s", testReplace),
				replace: testReplace,
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "this tests one token0",
				err: nil,
			},
		},
		{
			name: "several tokens",
			args: args{
				stmt:    fmt.Sprintf("this tests %s the token %s", testReplace, testReplace),
				replace: testReplace,
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "this tests 0 the token 1",
				err: nil,
			},
		},
		{
			name: "several tokens",
			args: args{
				stmt:    fmt.Sprintf("%s%sthis tests %s%s the token %s%s", testReplace, testReplace, testReplace, testReplace, testReplace, testReplace),
				replace: testReplace,
				fn: func(i int) string {
					return fmt.Sprintf("%v", i)
				},
			},
			want: want{
				str: "01this tests 23 the token 45",
				err: nil,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			str, err := scanReplace(tt.args.stmt, testReplace, tt.args.fn)
			if tt.want.str != str {
				t.Errorf("unexpected str: want '%s' got '%s'", tt.want.str, str)
			}

			if tt.want.err != err {
				t.Errorf("unexpected err: want '%s' got '%s'", tt.want.err, err)
			}

		})
	}
}

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
