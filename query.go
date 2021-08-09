package bqb

import (
	"fmt"
	"strings"
)

type part struct {
	Text   string
	Params []interface{}
}

type Query struct {
	Parts   []part
	Prepend string
}

func New(text string, args ...interface{}) *Query {
	q := &Query{}
	q.Parts = append(q.Parts, makePart(text, args...))
	return q
}

func Empty(prep ...string) *Query {
	return &Query{
		Prepend: strings.Join(prep, " "),
	}
}

func (q *Query) And(text string, args ...interface{}) *Query {
	return q.Join(" AND ", text, args...)
}

func (q *Query) Comma(text string, args ...interface{}) *Query {
	return q.Join(",", text, args...)
}

func (q *Query) Concat(text string, args ...interface{}) *Query {
	return q.Join("", text, args...)
}

func (q *Query) Join(sep, text string, args ...interface{}) *Query {
	if len(q.Parts) > 0 {
		q.Parts = append(q.Parts, makePart(sep+text, args...))
	} else {
		q.Parts = append(q.Parts, makePart(text, args...))
	}

	return q
}

func (q *Query) Or(text string, args ...interface{}) *Query {
	return q.Join(" OR ", text, args...)
}

func (q *Query) Space(text string, args ...interface{}) *Query {
	return q.Join(" ", text, args...)
}

func (q *Query) Print() {
	sql, params, err := q.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	fmt.Printf("PARAMS: %v\n", params)
	fmt.Printf("ERROR: %v\n", err)
}

func (q *Query) ToMysql() (string, []interface{}, error) {
	sql, params, _ := q.toSql()
	sql, err := dialectReplace(MYSQL, sql, params)
	return sql, params, err
}

func (q *Query) ToPgsql() (string, []interface{}, error) {
	sql, params, _ := q.toSql()
	sql, err := dialectReplace(PGSQL, sql, params)
	return sql, params, err
}

func (q *Query) ToRaw() (string, error) {
	sql, params, _ := q.toSql()
	sql, err := dialectReplace(RAW, sql, params)
	return sql, err
}

func (q *Query) ToSql() (string, []interface{}, error) {
	sql, params, _ := q.toSql()
	sql, err := dialectReplace(SQL, sql, params)
	return sql, params, err
}

func (q *Query) toSql() (string, []interface{}, error) {
	var sql string
	var params []interface{}

	if q.Prepend != "" && len(q.Parts) > 0 {
		sql = q.Prepend + " "
	}

	for _, p := range q.Parts {
		sql += p.Text
		params = append(params, p.Params...)
	}

	return strings.TrimSpace(sql), params, nil
}
