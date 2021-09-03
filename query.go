package bqb

import (
	"errors"
	"fmt"
	"strings"
)

type QueryPart struct {
	Text   string
	Params []interface{}
}

type Query struct {
	Parts          []QueryPart
	OptionalPrefix string
}

// New returns an instance of Query with a single QueryPart
func New(text string, args ...interface{}) *Query {
	q := Q()
	q.Parts = append(q.Parts, makePart(text, args...))
	return q
}

// Q returns a new empty Query
func Q() *Query {
	return &Query{}
}

// Optional returns a query object that has a conditional prefix which only
// resolves when at least one QueryPart has been added.
func Optional(prefix string) *Query {
	return &Query{
		OptionalPrefix: prefix,
	}
}

// And joins the current QueryPart to the previous QueryPart with ' AND '
func (q *Query) And(text string, args ...interface{}) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(" AND ", text, args...)
}

// Comma joins the current QueryPart to the previous QueryPart with a comma
func (q *Query) Comma(text string, args ...interface{}) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(",", text, args...)
}

// Concat concatenates the current QueryPart to the previous QueryPart with a
// zero space string
func (q *Query) Concat(text string, args ...interface{}) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join("", text, args...)
}

// Join joins the current QueryPart to the previous QueryPart with `sep`
func (q *Query) Join(sep, text string, args ...interface{}) *Query {
	if q == nil {
		return New(text, args...)
	}
	if len(q.Parts) > 0 {
		q.Parts = append(q.Parts, makePart(sep+text, args...))
	} else {
		q.Parts = append(q.Parts, makePart(text, args...))
	}

	return q
}

// Len returns the length of Query.Parts
func (q *Query) Len() int {
	if q == nil {
		return 0
	}
	return len(q.Parts)
}

// Or joins the current QueryPart to the previous QueryPart with ' OR '
func (q *Query) Or(text string, args ...interface{}) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(" OR ", text, args...)
}

// Print outputs the sql, parameters, and errors of a Query
func (q *Query) Print() {
	sql, params, err := q.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	fmt.Printf("PARAMS: %v\n", params)
	fmt.Printf("ERROR: %v\n", err)
}

// Space joins the current QueryPart to the previous QueryPart with a space
func (q *Query) Space(text string, args ...interface{}) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(" ", text, args...)
}

func (q *Query) ToMysql() (string, []interface{}, error) {
	sql, params, err := q.toSql()
	if err != nil {
		return "", nil, err
	}
	sql, err = dialectReplace(MYSQL, sql, params)
	return sql, params, err
}

// ToPgsql returns the sql placeholders with dollarsign format used by postgres
func (q *Query) ToPgsql() (string, []interface{}, error) {
	sql, params, err := q.toSql()
	if err != nil {
		return "", nil, err
	}
	sql, err = dialectReplace(PGSQL, sql, params)
	return sql, params, err
}

// ToRaw returns a string which the parameters have been resolved added
// as correctly as possible.
func (q *Query) ToRaw() (string, error) {
	sql, params, err := q.toSql()
	if err != nil {
		return "", err
	}

	sql, err = dialectReplace(RAW, sql, params)
	return sql, err
}

func (q *Query) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql()
	if err != nil {
		return "", nil, err
	}
	sql, err = dialectReplace(SQL, sql, params)
	return sql, params, err
}

func (q *Query) toSql() (string, []interface{}, error) {
	if q == nil {
		return "", nil, errors.New("cannot get sql on nil Query")
	}
	var sql string
	var params []interface{}

	if q.OptionalPrefix != "" && len(q.Parts) > 0 {
		sql = q.OptionalPrefix + " "
	}

	for _, p := range q.Parts {
		sql += p.Text
		params = append(params, p.Params...)
	}

	return strings.TrimSpace(sql), params, nil
}
