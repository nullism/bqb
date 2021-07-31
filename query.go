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

func (q *Query) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql(MYSQL)
	return sql, params, err
}

func (q *Query) ToPsql() (string, []interface{}, error) {
	sql, params, err := q.toSql(PGSQL)
	return sql, params, err
}

func (q *Query) ToRaw() (string, error) {
	sql, _, err := q.toSql(RAW)
	return sql, err
}

func (q *Query) Print() {
	var sql string
	var params []interface{}
	var err error
	dialect := q.dialect
	if dialect == PGSQL {
		sql, params, err = q.ToPsql()
	} else if dialect == RAW {
		sql, err = q.ToRaw()
	} else {
		sql, params, err = q.ToSql()
	}
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
