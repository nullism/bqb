package bqb

import (
	"fmt"
	"strings"
)

func exprGroup(dialect string, exprs [][]Expr, params []interface{}) (string, []interface{}) {
	var sql string

	if len(exprs) > 0 {
		for i, group := range exprs {
			sql += "("
			for n, expr := range group {
				f := expr.F
				for _, v := range expr.V {
					params = append(params, v)
					if dialect == pgsql {
						f = strings.Replace(f, "?", fmt.Sprintf("$%d", len(params)), 1)
					}
				}
				sql += fmt.Sprintf("%v", f)
				println(fmt.Sprintf("%v %T", expr.V, expr.V))

				if n+1 < len(group) {
					sql += " AND "
				}
			}
			if i+1 == len(exprs) {
				sql += ") "
			} else {
				sql += ") OR "
			}
		}
	}
	return sql, params
}

func (q *Query) toSql(dialect string) (string, []interface{}, error) {

	sql := ""
	var params []interface{}

	sql += fmt.Sprintf("SELECT %v ", q.S)

	if q.F != "" {
		sql += fmt.Sprintf("FROM %v ", q.F)
	}

	if len(q.J) > 0 {
		for _, j := range q.J {
			sql += fmt.Sprintf("JOIN %v ", j)
		}
	}

	if len(q.W) > 0 {
		sql += "WHERE "
		gsql, gparams := exprGroup(dialect, q.W, params)
		sql += gsql
		params = append(params, gparams...)
	}

	if q.GB != "" {
		sql += fmt.Sprintf("GROUP BY %v ", q.GB)

	}

	if len(q.H) > 0 {
		sql += "HAVING "
		hsql, hparams := exprGroup(dialect, q.H, params)
		sql += hsql
		params = append(params, hparams...)
	}

	if q.OB != "" {
		sql += fmt.Sprintf("ORDER BY %v ", q.OB)
	}

	if q.O != 0 {
		sql += fmt.Sprintf("OFFSET %d ", q.O)
	}

	if q.L != 0 {
		sql += fmt.Sprintf("LIMIT %d ", q.L)
	}

	return sql, params, nil
}

func (q *Query) ToSql() (string, []interface{}, error) {
	sql, params, err := q.toSql(mysql)
	return sql, params, err
}

func (q *Query) ToPsql() (string, []interface{}, error) {
	sql, params, err := q.toSql(pgsql)
	return sql, params, err
}

func (q *Query) Select(s string) *Query {
	q.S = s
	return q
}

func (q *Query) From(f string) *Query {
	q.F = f
	return q
}

func (q *Query) Join(j string) *Query {
	q.J = append(q.J, j)
	return q
}

func (q *Query) Where(exprs ...Expr) *Query {
	q.W = append(q.W, exprs)
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

func (q *Query) GroupBy(gb string) *Query {
	q.GB = gb
	return q
}

func (q *Query) Having(exprs ...Expr) *Query {
	q.H = append(q.H, exprs)
	return q
}

func (q *Query) OrderBy(expr string) *Query {
	q.OB = expr
	return q
}

func Valf(expr string, vals ...interface{}) Expr {
	e := Expr{
		F: expr,
		V: vals,
	}
	return e
}
