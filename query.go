package bqb

import (
	"fmt"
	"strings"
)

type select_ struct {
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

func Select(exprs ...interface{}) *select_ {
	q := &select_{
		dialect: SQL,
		selects: getExprs(exprs),
	}
	return q
}

func QueryPsql() *select_ {
	return &select_{dialect: PGSQL}
}

func QueryMysql() *select_ {
	return &select_{dialect: MYSQL}
}

func QuerySql() *select_ {
	return &select_{dialect: SQL}
}

func QueryRaw() *select_ {
	return &select_{dialect: RAW}
}

func (q *select_) Group(exprs ...interface{}) *select_ {
	q.groups = append(q.groups, getExprs(exprs)...)
	return q
}

func (q *select_) Select(exprs ...interface{}) *select_ {
	newExprs := getExprs(exprs)
	q.selects = append(q.selects, newExprs...)
	return q
}

func (q *select_) From(exprs ...interface{}) *select_ {
	newExprs := getExprs(exprs)
	q.from = append(q.from, newExprs...)
	return q
}

func (q *select_) Join(exprs ...interface{}) *select_ {
	// newExprs := getExprs(exprs)
	// q.join = append(q.join, newExprs...)
	q.JoinType("JOIN", exprs...)
	return q
}

func (q *select_) JoinType(kind string, exprs ...interface{}) *select_ {
	for _, expr := range exprs {
		q.joins = append(
			q.joins, join{kind: kind, expr: intfToExpr(expr)},
		)
	}
	return q
}

func (q *select_) Where(exprs ...interface{}) *select_ {
	newExprs := getExprs(exprs)
	q.where = append(q.where, newExprs...)
	return q
}

func (q *select_) Limit(limit int) *select_ {
	q.limit = limit
	return q
}

func (q *select_) Offset(offset int) *select_ {
	q.offset = offset
	return q
}

func (q *select_) GroupBy(exprs ...interface{}) *select_ {
	newExprs := getExprs(exprs)
	q.groupBy = append(q.groupBy, newExprs...)
	return q
}

func (q *select_) Having(exprs ...interface{}) *select_ {
	newExprs := getExprs(exprs)
	q.having = append(q.having, newExprs...)
	return q
}

func (q *select_) OrderBy(exprs ...interface{}) *select_ {
	newExprs := getExprs(exprs)
	q.order = append(q.order, newExprs...)
	return q
}

func (q *select_) Enclose() *select_ {
	q.enclose = true
	return q
}

func (q *select_) As(name string) *select_ {
	q.as = name
	q.enclose = true
	return q
}

func (q *select_) Print() {
	sql, params, err := q.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func (q *select_) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql()
	return dialectReplace(q.dialect, sql, params), params, err
}

func (q *select_) toSql() (string, []interface{}, error) {
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
