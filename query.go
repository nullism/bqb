package bqb

import (
	"errors"
	"fmt"
	"strings"
)

// QueryPart holds a section of a Query.
type QueryPart struct {
	Text   string
	Params []any
	Errs   []error
}

// Query contains all the QueryParts for the query and is the primary
// struct of the bqb package.
type Query struct {
	Parts          []QueryPart
	OptionalPrefix string
}

// New returns an instance of Query with a single QueryPart.
func New(text string, args ...any) *Query {
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

// And joins the current QueryPart to the previous QueryPart with ' AND '.
func (q *Query) And(text string, args ...any) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(" AND ", text, args...)
}

// Comma joins the current QueryPart to the previous QueryPart with a comma.
func (q *Query) Comma(text string, args ...any) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(",", text, args...)
}

// Concat concatenates the current QueryPart to the previous QueryPart with a
// zero space string.
func (q *Query) Concat(text string, args ...any) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join("", text, args...)
}

// Empty returns true if the Query is nil or has a length > 0.
func (q *Query) Empty() bool {
	if q == nil {
		return true
	}
	return q.Len() == 0
}

// Join joins the current QueryPart to the previous QueryPart with `sep`.
func (q *Query) Join(sep, text string, args ...any) *Query {
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

// Or joins the current QueryPart to the previous QueryPart with ' OR '.
func (q *Query) Or(text string, args ...any) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(" OR ", text, args...)
}

// Print outputs the sql, parameters, and errors of a Query.
func (q *Query) Print() {
	sql, params, err := q.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	fmt.Printf("PARAMS: %v\n", params)
	fmt.Printf("ERROR: %v\n", err)
}

// Space joins the current QueryPart to the previous QueryPart with a space.
func (q *Query) Space(text string, args ...any) *Query {
	if q == nil {
		return New(text, args...)
	}
	return q.Join(" ", text, args...)
}

// ToMysql returns the sql placeholders with SQL (?) format used by MySQL
func (q *Query) ToMysql() (string, []any, error) {
	sql, params, err := q.toSql()
	if err != nil {
		return "", nil, err
	}
	sql, err = dialectReplace(MYSQL, sql, params)
	return sql, params, err
}

// ToPgsql returns the sql placeholders with dollarsign format used by postgres.
func (q *Query) ToPgsql() (string, []any, error) {
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

// ToSql returns the placeholders with question (?) format used by most
// databases such as sqlite, mysql, and others.
func (q *Query) ToSql() (string, []any, error) {
	sql, params, err := q.toSql()
	if err != nil {
		return "", nil, err
	}
	sql, err = dialectReplace(SQL, sql, params)
	return sql, params, err
}

func (q *Query) toSql() (string, []any, error) {
	if q == nil {
		return "", nil, errors.New("cannot get sql on nil Query")
	}
	var sql string
	var params []any

	if q.OptionalPrefix != "" && len(q.Parts) > 0 {
		sql = q.OptionalPrefix + " "
	}

	for _, p := range q.Parts {
		sql += p.Text
		params = append(params, p.Params...)

		if len(p.Errs) != 0 {
			return "", nil, errors.Join(p.Errs...)
		}
	}

	return strings.TrimSpace(sql), params, nil
}
