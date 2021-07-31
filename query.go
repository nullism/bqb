package bqb

import "fmt"

type Query struct {
	SE []Expr
	FE []Expr
	JE []Expr
	W  [][]Expr
	OB []Expr
	L  int
	O  int
	GB []Expr
	H  [][]Expr
}

func (q *Query) Select(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.SE = append(q.SE, newExprs...)
	return q
}

func (q *Query) From(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.FE = append(q.FE, newExprs...)
	return q
}

func (q *Query) Join(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.JE = append(q.JE, newExprs...)
	return q
}

func (q *Query) Where(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.W = append(q.W, newExprs)
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.L = limit
	return q
}

func (q *Query) Offset(offset int) *Query {
	q.O = offset
	return q
}

func (q *Query) GroupBy(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.GB = append(q.GB, newExprs...)
	return q
}

func (q *Query) Having(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.H = append(q.H, newExprs)
	return q
}

func (q *Query) OrderBy(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.OB = append(q.OB, newExprs...)
	return q
}

func (q *Query) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql(MYSQL)
	return sql, params, err
}

func (q *Query) ToPsql() (string, []interface{}, error) {
	sql, params, err := q.toSql(PGSQL)
	return sql, params, err
}

func (q *Query) Print(dialect string) {
	var sql string
	var params []interface{}
	var err error

	if dialect == PGSQL {
		sql, params, err = q.ToPsql()
	} else {
		sql, params, err = q.ToSql()
	}
	fmt.Printf("SQL: %v\nPARAMS: %v\nERROR: %v\n", sql, params, err)
}
