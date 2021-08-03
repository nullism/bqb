package bqb

import (
	"fmt"
	"strings"
)

type selectQ struct {
	dialect string
	selects []Expr
	from    []Expr
	joins   []join
	where   []Expr
	order   []Expr
	limit   int
	offset  int
	groupBy []Expr
	having  []Expr
	groups  []Expr
	as      string
	enclose bool
}

type join struct {
	kind string
	expr Expr
}

func Select(exprs ...interface{}) *selectQ {
	q := &selectQ{
		dialect: SQL,
		selects: getExprs(exprs),
	}
	return q
}

func (s *selectQ) Postgres() *selectQ {
	s.dialect = PGSQL
	return s
}

func (s *selectQ) Mysql() *selectQ {
	s.dialect = MYSQL
	return s
}

func (s *selectQ) Raw() *selectQ {
	s.dialect = RAW
	return s
}

func QueryPsql() *selectQ {
	return &selectQ{dialect: PGSQL}
}

func (q *selectQ) Group(exprs ...interface{}) *selectQ {
	q.groups = append(q.groups, getExprs(exprs)...)
	return q
}

func (q *selectQ) Select(exprs ...interface{}) *selectQ {
	newExprs := getExprs(exprs)
	q.selects = append(q.selects, newExprs...)
	return q
}

func (q *selectQ) From(exprs ...interface{}) *selectQ {
	newExprs := getExprs(exprs)
	q.from = append(q.from, newExprs...)
	return q
}

func (q *selectQ) Join(exprs ...interface{}) *selectQ {
	// newExprs := getExprs(exprs)
	// q.join = append(q.join, newExprs...)
	q.JoinType("JOIN", exprs...)
	return q
}

func (q *selectQ) JoinType(kind string, exprs ...interface{}) *selectQ {
	for _, expr := range exprs {
		q.joins = append(
			q.joins, join{kind: kind, expr: intfToExpr(expr)},
		)
	}
	return q
}

func (q *selectQ) Where(exprs ...interface{}) *selectQ {
	newExprs := getExprs(exprs)
	q.where = append(q.where, newExprs...)
	return q
}

func (q *selectQ) Limit(limit int) *selectQ {
	q.limit = limit
	return q
}

func (q *selectQ) Offset(offset int) *selectQ {
	q.offset = offset
	return q
}

func (q *selectQ) GroupBy(exprs ...interface{}) *selectQ {
	newExprs := getExprs(exprs)
	q.groupBy = append(q.groupBy, newExprs...)
	return q
}

func (q *selectQ) Having(exprs ...interface{}) *selectQ {
	newExprs := getExprs(exprs)
	q.having = append(q.having, newExprs...)
	return q
}

func (q *selectQ) OrderBy(exprs ...interface{}) *selectQ {
	newExprs := getExprs(exprs)
	q.order = append(q.order, newExprs...)
	return q
}

func (q *selectQ) Enclose() *selectQ {
	q.enclose = true
	return q
}

func (q *selectQ) As(name string) *selectQ {
	q.as = name
	q.enclose = true
	return q
}

func (q *selectQ) Print() {
	sql, params, err := q.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func (q *selectQ) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql()
	return dialectReplace(q.dialect, sql, params), params, err
}

func (q *selectQ) toSql() (string, []interface{}, error) {
	var params []interface{}
	sql := ""

	if len(q.selects) > 0 {
		sql += "SELECT "
		nsql, nparams := exprsToSql(q.selects)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(q.from) > 0 {
		sql += "FROM "
		nsql, nparams := exprsToSql(q.from)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(q.joins) > 0 {
		for _, join := range q.joins {
			sql += join.kind + " "
			sql += join.expr.F
			params = append(params, join.expr.V...)
			sql += " "
		}
	}

	if len(q.where) > 0 {
		sql += "WHERE "
		nsql, nparams := exprsToSql(q.where)
		sql += strings.Join(nsql, " OR ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(q.groupBy) > 0 {
		sql += "GROUP BY "
		nsql, nparams := exprsToSql(q.groupBy)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(q.having) > 0 {
		sql += "HAVING "
		nsql, nparams := exprsToSql(q.having)
		sql += strings.Join(nsql, " AND ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(q.order) > 0 {
		sql += "ORDER BY "
		nsql, nparams := exprsToSql(q.order)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if q.offset != 0 {
		sql += fmt.Sprintf("OFFSET %d ", q.offset)
	}

	if q.limit != 0 {
		sql += fmt.Sprintf("LIMIT %d ", q.limit)
	}

	if len(q.groups) > 0 {
		sql += "("
		nsql, nparams := exprsToSql(q.groups)
		sql += strings.Join(nsql, " ")
		params = append(params, nparams...)
		sql += ") "
	}

	sql = strings.TrimSpace(sql)

	if q.enclose {
		sql = fmt.Sprintf("(%v)", sql)
	}

	if q.as != "" {
		sql = fmt.Sprintf("%v AS %v", sql, q.as)
	}

	return sql, params, nil
}
