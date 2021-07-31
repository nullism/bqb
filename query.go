package bqb

import "fmt"

type Query struct {
	dialect string
	SE      []Expr
	FE      []Expr
	JE      []Expr
	W       [][]Expr
	OB      []Expr
	L       int
	O       int
	GB      []Expr
	H       [][]Expr
}

func New(dialect string) *Query {
	return &Query{dialect: dialect}
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

func (q *Query) Print() {
	sql, params, err := q.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
