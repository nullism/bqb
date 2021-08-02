package bqb

import (
	"fmt"
	"strings"
)

type Query struct {
	dialect string
	selects []Expr
	from    []Expr
	join    []Expr
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

func QueryPsql() *Query {
	return &Query{dialect: PGSQL}
}

func QueryMysql() *Query {
	return &Query{dialect: MYSQL}
}

func QuerySql() *Query {
	return &Query{dialect: SQL}
}

func QueryRaw() *Query {
	return &Query{dialect: RAW}
}

func (q *Query) Group(exprs ...interface{}) *Query {
	q.groups = append(q.groups, getExprs(exprs)...)
	return q
}

func (q *Query) Select(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.selects = append(q.selects, newExprs...)
	return q
}

func (q *Query) From(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.from = append(q.from, newExprs...)
	return q
}

func (q *Query) Join(exprs ...interface{}) *Query {
	// newExprs := getExprs(exprs)
	// q.join = append(q.join, newExprs...)
	q.JoinType("JOIN", exprs...)
	return q
}

func (q *Query) JoinType(kind string, exprs ...interface{}) *Query {
	for _, expr := range exprs {
		q.joins = append(
			q.joins, join{kind: kind, expr: intfToExpr(expr)},
		)
	}
	return q
}

func (q *Query) Where(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.where = append(q.where, newExprs...)
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

func (q *Query) GroupBy(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.groupBy = append(q.groupBy, newExprs...)
	return q
}

func (q *Query) Having(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.having = append(q.having, newExprs...)
	return q
}

func (q *Query) OrderBy(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.order = append(q.order, newExprs...)
	return q
}

func (q *Query) Enclose() *Query {
	q.enclose = true
	return q
}

func (q *Query) As(name string) *Query {
	q.as = name
	q.enclose = true
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

func (q *Query) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql()
	return dialectReplace(q.dialect, sql, params), params, err
}

func (q *Query) toSql() (string, []interface{}, error) {
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
		sql = fmt.Sprintf("%v as %v", sql, q.as)
	}

	return sql, params, nil
}
