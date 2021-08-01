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
	where   []Expr
	order   []Expr
	limit   int
	offset  int
	groupBy []Expr
	having  []Expr
}

func QueryPsql() *Query {
	return &Query{dialect: PGSQL}
}

func QueryMysql() *Query {
	return &Query{dialect: MYSQL}
}

func QueryRaw() *Query {
	return &Query{dialect: RAW}
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
	newExprs := getExprs(exprs)
	q.join = append(q.join, newExprs...)
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

	sql := ""
	var params []interface{}

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

	if len(q.join) > 0 {
		sql += "JOIN "
		nsql, nparams := exprsToSql(q.join)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
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

	for i, p := range params {
		if q.dialect == RAW {
			switch v := p.(type) {
			case nil:
				sql = strings.Replace(sql, paramPh, "NULL", 1)
			case int, bool:
				sql = strings.Replace(sql, paramPh, fmt.Sprintf("%v", v), 1)
			default:
				sql = strings.Replace(sql, paramPh, fmt.Sprintf("'%v'", v), 1)
			}
		} else if q.dialect == MYSQL || q.dialect == SQL {
			sql = strings.Replace(sql, paramPh, "?", 1)
		} else if q.dialect == PGSQL {
			sql = strings.Replace(sql, paramPh, fmt.Sprintf("$%d", i+1), 1)
		}

	}

	return sql, params, nil
}
