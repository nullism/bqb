package bqb

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestA(t *testing.T) {
}

func TestArrays(t *testing.T) {
	q := New("(?) (?) (?) (?) (?)", []string{"a", "b"}, []string{}, []*string{}, []int{1, 2}, []*int{})
	sql, params, _ := q.ToSql()

	if len(params) != 6 {
		t.Errorf("invalid params")
	}

	want := "(?,?) () (?) (?,?) (?)"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestJson(t *testing.T) {
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

	q = New("?", JsonMap{"a": func() {}})
	badsql, badargs, err := q.ToSql()
	if err == nil {
		t.Errorf("expected error with invalid JsonMap")
	}
	if !strings.Contains(err.Error(), "jsonify") {
		t.Errorf("invalid error from *JsonMap failure")
	}
	if badsql != "" {
		t.Errorf("expected emtpy SQL but got %v instead", badsql)
	}
	if len(badargs) > 0 {
		t.Errorf("expected no args but got %v instead", badargs)
	}
}

func TestJsonPointer(t *testing.T) {
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

	q = New("?", &JsonMap{"a": func() {}})
	_, _, err := q.ToSql()
	if err == nil {
		t.Errorf("expected error for invalid *JsonMap")
	}
	if !strings.Contains(err.Error(), "jsonify") {
		t.Errorf("invalid error from *JsonMap failure")
	}

}

func TestOptional(t *testing.T) {
	sel := Optional("you should not see this")

	if !sel.Empty() || sel.Len() != 0 {
		t.Errorf("Optional is not empty")
	}

	sql, _ := sel.ToRaw()

	if sql != "" {
		t.Errorf("Empty() not returning empty string")
	}

	sel.Space("but now you can")

	if sel.Empty() || sel.Len() == 0 {
		t.Errorf("Optional should not be empty when extended")
	}

	sql, _ = sel.ToRaw()
	want := "you should not see this but now you can"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestNils(t *testing.T) {
	var q *Query
	_, _, err := q.ToSql()
	if err == nil {
		t.Errorf("expected error for ToSql")
	}

	_, err = q.ToRaw()
	if err == nil {
		t.Errorf("expected error for ToRaw")
	}

	_, _, err = q.ToMysql()
	if err == nil {
		t.Errorf("expected error for ToMysql")
	}

	_, _, err = q.ToPgsql()
	if err == nil {
		t.Errorf("expected error for ToPgsql")
	}

	var qNil *Query
	qNil.And("test")
	_, _, err = qNil.ToSql()
	if err == nil {
		t.Errorf("expected error for qNil")
	}

	qNil.Comma("test")
	_, _, err = qNil.ToMysql()
	if err == nil {
		t.Errorf("expected error for qNil")
	}

	qNil.Concat("test")
	_, _, err = qNil.ToPgsql()
	if err == nil {
		t.Errorf("expected error for qNil")
	}

	if !qNil.Empty() {
		t.Errorf("expected Empty == true for qNil")
	}

	qNil.Join("", "test")
	_, _, err = qNil.ToPgsql()
	if err == nil {
		t.Errorf("expected error for qNil")
	}

	if qNil.Len() != 0 {
		t.Errorf("expected zero length of qNil")
	}

	qNil.Or("test")
	_, err = qNil.ToRaw()
	if err == nil {
		t.Errorf("expected error for qNil")
	}

	qNil.Space("test")
	_, err = qNil.ToRaw()
	if err == nil {
		t.Errorf("expected error for qNil")
	}

}

func TestParamsExtra(t *testing.T) {
	q := New("params ? ?", 1)
	_, _, err := q.ToSql()
	if err == nil {
		t.Errorf("no error for extra ?")
	}
	if !strings.Contains(err.Error(), "extra") {
		t.Errorf("got wrong error for extra ?")
	}
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
	sql, err := New("?", []any{"a", 1, true}).ToRaw()
	if err != nil {
		t.Errorf("got error %v", err)
	}

	want := "'a',1,true"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestParamsMissing(t *testing.T) {
	q := New("params ?", 1, 2)
	_, _, err := q.ToSql()
	if err == nil {
		t.Errorf("missing ? considered valid")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Errorf("got wrong error for missing ?")
	}
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

func TestQuery_And(t *testing.T) {
	q := New("a")
	q.And("b")
	q.And("c")

	sql, _, _ := q.ToSql()
	want := "a AND b AND c"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestQuery_Comma(t *testing.T) {
	q := New("a")
	q.Comma("b")

	sql, _, _ := q.ToSql()
	want := "a,b"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestQuery_Concat(t *testing.T) {
	q := New("a")
	q.Concat("b")

	if q.Len() != 2 {
		t.Errorf("invalid length for query: %v", q.Len())
	}

	sql, _, _ := q.ToSql()
	want := "ab"
	if sql != want {
		t.Errorf("want: %q, got: %q", want, sql)
	}
}

func TestQuery_Empty(t *testing.T) {
	child := Optional("EMPTY")
	parent := New("?", child)
	sql, _, _ := parent.ToSql()

	if sql != "" {
		t.Errorf("expected empty string in empty query: %q", sql)
	}
}

func TestQuery_Len(t *testing.T) {
	q := Optional("a")
	if q.Len() != 0 {
		t.Errorf("expected 0 length")
	}

	q.Comma("1")
	if q.Len() != 1 {
		t.Errorf("expected length of 1")
	}

	q.Comma("2")
	q.Comma("3")
	q.Comma("4")
	q.Comma("5")
	if q.Len() != 5 {
		t.Errorf("expected length of 5")
	}
}

func TestQuery_Or(t *testing.T) {
	q := New("a")
	q.Or("b")
	q.Or("c")

	sql, _ := q.ToRaw()
	want := "a OR b OR c"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQuery_Space(t *testing.T) {
	q := New("a")
	q.Space("b")

	sql, _, _ := q.ToSql()
	want := "a b"
	if sql != want {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQuery_ToMysql(t *testing.T) {
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

func TestQuery_ToMysqlTime(t *testing.T) {
	var names []string
	for i := 0; i < 10000; i++ {
		names = append(names, fmt.Sprintf("n%d", i))
	}

	q := New("TEST IN (?)", names)
	q2 := New("Secondary ?", q)
	q3 := New("Ternary ?", q2)

	start := time.Now()
	_, params, err := q3.ToMysql()
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	delta := time.Since(start)
	if delta.Seconds() > 1 {
		t.Errorf("ToMysql took too long to return: %v", delta.Seconds())
	}

	for i := 0; i < len(names); i++ {
		if params[i] != names[i] {
			t.Errorf("got: %v, want: %v", params[i], names[i])
		}
	}

}

func TestQuery_ToPgsql(t *testing.T) {
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

func TestQuery_ToRaw(t *testing.T) {
	q := New(
		"bool:? float:? int:? string:? []int:? []string:? Query:? Json:? nil:?",
		true, 1.5, 1, "2", []int{3, 3}, []string{"4", "4"}, New("5"), JsonMap{"6": 6}, nil,
	)
	sql, err := q.ToRaw()
	if err != nil {
		t.Errorf("got error from ToRaw(): %v", err)
	}

	want := "bool:true float:1.5 int:1 string:'2' []int:3,3 []string:'4','4' Query:5 Json:'{\"6\":6}' nil:NULL"
	if want != sql {
		t.Errorf("got: %q, want: %q", sql, want)
	}
}

func TestQuery_ToSql(t *testing.T) {
	text := "bool:? float:? int:? string:? []int:? []string:? Query:? Json:? nil:?"
	args := []any{
		true, 1.5, 1, "test", []int{1, 2}, []string{"a", "b", "c"}, New("(Q ?)", "q"), &JsonMap{"a": 1}, nil,
	}

	q := New(text, args...)
	sql, params, err := q.ToSql()
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	wantP := []any{
		true, 1.5, 1, "test", 1, 2, "a", "b", "c", "q", `{"a":1}`, nil,
	}
	want := "bool:? float:? int:? string:? []int:?,? []string:?,?,? Query:(Q ?) Json:? nil:?"
	if want != sql {
		t.Errorf("\n got: %q\nwant: %q", sql, want)
	}

	for i := range params {
		if params[i] != wantP[i] {
			t.Errorf("got: %v %T, want: %v %T", params[i], params[i], wantP[i], wantP[i])
		}
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

func TestQueryNil(t *testing.T) {
	var q *Query
	q2 := New("test ?", q)

	sql, params, _ := q2.ToSql()
	want := "test ?"
	if want != sql {
		t.Errorf("got: %q, want:%q", sql, want)
	}

	if params[0] != nil {
		t.Errorf("invalid param: %v", params[0])
	}
}

func TestQueryPrint(t *testing.T) {
	q := New("SELECT * FROM table WHERE name = ?", "name")
	q.Print()
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

	boolX := true

	intX := 1
	intsX := []int{2, 2}

	stringX := "s"
	stringsX := []string{"s1", "s2"}

	intp := &intX
	var intpn *int
	intsp := []*int{&intX, &intX}
	var intspn []*int

	stringp := &stringX
	var stringpn *string
	stringsp := []*string{&stringX, &stringX}
	var stringspn []*string

	jsonX := JsonMap{"a": 1}
	var jsonp *JsonMap

	text := "b:? - i:? ? - s:? ? - ip:? ? ? ? - sp:? ? ? ? - j:? ?"
	q := New(text,
		boolX,
		intX, intsX,
		stringX, stringsX,
		intp, intpn, intsp, intspn,
		stringp, stringpn, stringsp, stringspn,
		jsonX, jsonp)
	sql, _ := q.ToRaw()
	want := `b:true - i:1 2,2 - s:'s' 's1','s2' - ip:1 NULL 1,1 NULL - sp:'s' NULL 's','s' NULL - j:'{"a":1}' 'null'`
	if want != sql {
		t.Errorf("\ngot : %q\nwant: %q", sql, want)
	}
}

type valuer []string

func (v valuer) Value() (driver.Value, error) {
	// Force an error with nil for test purposes
	if v == nil {
		return nil, errors.New("error creating value")
	}
	return strings.Join(v, "/"), nil
}

func TestValuer(t *testing.T) {
	q := New("(?)", valuer{"a", "b", "c"})
	sql, params, _ := q.ToSql()

	if len(params) != 1 {
		t.Errorf("invalid params")
	}

	wantSql := "(?)"
	if sql != wantSql {
		t.Errorf("got: %q, want: %q", sql, wantSql)
	}

	wantParam := "a/b/c"
	if params[0].(string) != wantParam {
		t.Errorf("got: %q, want: %q", sql, wantParam)
	}

	sql, err := q.ToRaw()
	if err != nil {
		t.Errorf("got error from valuer ToRaw(): %v", err)
	}
	if sql != "('a/b/c')" {
		t.Errorf("got unexpected value from ToRaw(): %v", sql)
	}
}

func TestValuerError(t *testing.T) {
	var v valuer
	q := New("?", v)
	_, _, err := q.ToSql()

	wantError := "error creating value"
	if err.Error() != wantError {
		t.Errorf("got: %q, want: %q", err, wantError)
	}
}
