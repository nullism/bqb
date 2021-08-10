package bqb

import (
	"fmt"
	"strings"
	"testing"
)

func TestA(t *testing.T) {

}

func TestArrays(t *testing.T) {
	q := New("(?) (?) (?) (?)", []string{"a", "b"}, []*string{}, []int{1, 2}, []*int{})
	sql, params, _ := q.ToSql()

	if len(params) != 6 {
		t.Errorf("invalid params")
	}

	want := "(?,?) (?) (?,?) (?)"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestSpace(t *testing.T) {
	q := New("a")
	q.Space("b")

	sql, _, _ := q.ToSql()
	want := "a b"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestAnd(t *testing.T) {
	q := New("a")
	q.And("b")
	q.And("c")

	sql, _, _ := q.ToSql()
	want := "a AND b AND c"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestComma(t *testing.T) {
	q := New("a")
	q.Comma("b")

	sql, _, _ := q.ToSql()
	want := "a,b"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestConcat(t *testing.T) {
	q := New("a")
	q.Concat("b")

	sql, _, _ := q.ToSql()
	want := "ab"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestEmpty(t *testing.T) {
	sel := Optional("you should not see this")
	sql, _ := sel.ToRaw()

	if sql != "" {
		t.Errorf("Empty() not returning empty string")
	}

	sel.Space("but now you can")

	sql, _ = sel.ToRaw()
	want := "you should not see this but now you can"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestParamsExtra(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if !strings.Contains(r.(string), "extra") {
				t.Errorf("invalid panic for missing params: %v", r)
			}
		}
	}()

	New("params ? ?", 1)
	t.Errorf("extra ? considered valid")
}

func TestParamsFunc(t *testing.T) {
	q := New("?", func(x int) int { return x })
	sql, err := q.ToRaw()
	if err == nil {
		t.Errorf("got nil error for invalid raw parameter")
	}

	if !strings.Contains(fmt.Sprint(err), "func(int) int") {
		t.Errorf("got incorrect error %v", err)
	}

	if sql != "" {
		t.Errorf("got non-empty value from ToRaw(): %q", sql)
	}

}

func TestParamsInterfaceList(t *testing.T) {
	sql, err := New("?", []interface{}{"a", 1, true}).ToRaw()
	if err != nil {
		t.Errorf("got error %v", err)
	}

	want := "'a',1,true"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestParamsJson(t *testing.T) {

	sql, _ := New(
		"INSERT INTO my_table (json_map,json_list) VALUES (?,?)",
		JsonMap{"a": 1, "b": []string{"a", "b"}},
		JsonList{"string", 1, true},
	).ToRaw()

	want := `INSERT INTO my_table (json_map,json_list) ` +
		`VALUES ('{"a":1,"b":["a","b"]}','["string",1,true]')`
	if sql != want {
		t.Errorf("\n got: %q\nwant: %q", sql, want)
	}

	q := New("INSERT INTO foo (json) VALUES (?)", JsonMap{"a": "test", "b": []int{1, 2}})

	sql, params, _ := q.ToSql()
	want = "INSERT INTO foo (json) VALUES (?)"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}

	pwant := `{"a":"test","b":[1,2]}`
	if params[0] != pwant {
		t.Errorf("want: %q, got: %q", pwant, params[0])
	}

	q = New("a = ?", JsonList{"a", 1, true})
	sql, params, _ = q.ToSql()
	jlpwant := `["a",1,true]`

	if params[0] != jlpwant {
		t.Errorf("got: %q, want: %q", params[0], jlpwant)
	}

	jlwant := "a = ?"
	if sql != jlwant {
		t.Errorf("got: %q, want: %q", sql, jlwant)
	}

	defer func() {
		if r := recover(); r != nil {
			if !strings.Contains(r.(string), "jsonify") {
				t.Errorf("invalid panic for missing params: %v", r)
			}
		}
	}()
	New("?", JsonMap{"a": func() {}})
	t.Errorf("didn't panic")
}

func TestParamsJsonPointer(t *testing.T) {
	q := New("INSERT INTO foo (json) VALUES (?)", &JsonMap{"a": "test", "b": []int{1, 2}})

	sql, params, _ := q.ToSql()
	want := "INSERT INTO foo (json) VALUES (?)"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}

	pwant := `{"a":"test","b":[1,2]}`
	if params[0] != pwant {
		t.Errorf("want: %q, got: %q", pwant, params[0])
	}

	q = New("a = ?", &JsonList{"a", 1, 2})
	sql, params, _ = q.ToSql()
	jlpwant := `["a",1,2]`

	if params[0] != jlpwant {
		t.Errorf("got: %q, want: %q", params[0], jlpwant)
	}

	jlwant := "a = ?"
	if sql != jlwant {
		t.Errorf("got: %q, want: %q", sql, jlwant)
	}

	defer func() {
		if r := recover(); r != nil {
			if !strings.Contains(r.(string), "jsonify") {
				t.Errorf("invalid panic for missing params: %v", r)
			}
		}
	}()
	New("?", &JsonMap{"a": func() {}})
	t.Errorf("didn't panic")
}

func TestParamsMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if !strings.Contains(r.(string), "missing") {
				t.Errorf("invalid panic for missing params: %v", r)
			}
		}
	}()

	New("params ?", 1, 2)
	t.Errorf("missing ? considered valid")
}

func TestQuery(t *testing.T) {
	q := New("SELECT * FROM table").
		Space("WHERE a = ? AND b = ?", 1, 2).
		Space("OR (b = 2 AND c = ?)", "hellos")

	sql, params, err := q.ToSql()

	if err != nil {
		t.Errorf("got error %v", err)
	}

	want := "SELECT * FROM table WHERE a = ? AND b = ? OR (b = 2 AND c = ?)"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}

	if len(params) != 3 {
		t.Errorf("got incorrect param count: %v", len(params))
	}
}

func TestQueryBuilding(t *testing.T) {
	sel := Optional("SELECT")

	sel.Comma("name")
	sel.Comma("id")

	from := Optional("FROM")
	from.Space("my_table")

	where := Optional("WHERE")

	adultCond := Q()
	adultCond.And("name = ?", "adult")
	adultCond.And("age > ?", 20)

	if len(adultCond.Parts) > 0 {
		where.And("(?)", adultCond)
	}

	where.Or("(name = ? AND age < ?)", "youth", 21)

	q := New("? ? ? LIMIT ?", sel, from, where, 10)

	sql, _ := q.ToRaw()
	want := "SELECT name,id FROM my_table WHERE (name = 'adult' AND age > 20) OR (name = 'youth' AND age < 21) LIMIT 10"

	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQueryLiteralQ(t *testing.T) {
	q := New("json->>field ?? val AND val = ?", "asdf")
	sql, _, _ := q.ToPgsql()
	want := "json->>field ? val AND val = $1"
	if want != sql {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQueryMysql(t *testing.T) {
	q := New("SELECT * FROM table WHERE a = ? AND b = ?", 1, "b")
	sql, params, _ := q.ToMysql()
	if len(params) != 2 {
		t.Errorf("expected two parameters, got %v", len(params))
	}

	want := "SELECT * FROM table WHERE a = ? AND b = ?"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQueryPgsql(t *testing.T) {
	q := New("SELECT name,").
		Space("(SELECT * FROM other_table WHERE name = ?) as other_name", "test").
		Space("FROM table LIMIT ?", 1)

	sql, params, _ := q.ToPgsql()
	if len(params) != 2 {
		t.Errorf("got incorrect param count: %v", len(params))
	}

	want := "SELECT name, (SELECT * FROM other_table WHERE name = $1) as other_name FROM table LIMIT $2"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQueryPrint(t *testing.T) {
	q := New("SELECT * FROM table WHERE name = ?", "name")
	q.Print()
}

func TestQueryRaw(t *testing.T) {

	q := New(
		"bool:? float:? int:? string:? []int:? []string:? Query:? Json:? nil:?",
		true, 1.5, 1, "2", []int{3, 3}, []string{"4", "4"}, New("5"), JsonMap{"6": 6}, nil,
	)
	sql, _ := q.ToRaw()

	want := "bool:true float:1.5 int:1 string:'2' []int:3,3 []string:'4','4' Query:5 Json:'{\"6\":6}' nil:NULL"
	if want != sql {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQuerySubquery(t *testing.T) {

	q := New(
		"SELECT name, (?) as time, age",
		New("SELECT time FROM logins LIMIT 1"),
	)
	q.Space("FROM users").
		Space("WHERE id NOT IN (?)", []string{"a", "b", "c", "d"}).
		Space("AND name NOT LIKE ?", "admin%").
		Space("LIMIT 1")

	sql, params, err := q.ToPgsql()

	if err != nil {
		t.Errorf("got error: %v", err)
	}

	if len(params) != 5 {
		t.Errorf("want 5 params, got %v", len(params))
	}

	want := "SELECT name, (SELECT time FROM logins LIMIT 1) as time, age FROM users WHERE id NOT IN ($1,$2,$3,$4) AND name NOT LIKE $5 LIMIT 1"
	if want != sql {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQueryTypes(t *testing.T) {

	bool_ := true

	int_ := 1
	ints_ := []int{2, 2}

	string_ := "s"
	strings_ := []string{"s1", "s2"}

	intp := &int_
	var intpn *int
	intsp := []*int{&int_, &int_}
	var intspn []*int

	stringp := &string_
	var stringpn *string
	stringsp := []*string{&string_, &string_}
	var stringspn []*string

	json_ := JsonMap{"a": 1}
	var jsonp *JsonMap

	text := "b:? - i:? ? - s:? ? - ip:? ? ? ? - sp:? ? ? ? - j:? ?"
	q := New(text,
		bool_,
		int_, ints_,
		string_, strings_,
		intp, intpn, intsp, intspn,
		stringp, stringpn, stringsp, stringspn,
		json_, jsonp)
	sql, _ := q.ToRaw()
	want := `b:true - i:1 2,2 - s:'s' 's1','s2' - ip:1 NULL 1,1 NULL - sp:'s' NULL 's','s' NULL - j:'{"a":1}' 'null'`
	if want != sql {
		t.Errorf("\ngot : %q\nwant: %q", sql, want)
	}
}
