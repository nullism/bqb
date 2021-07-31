package bqb

import (
	"fmt"
	"strings"
)

func getExprs(exprs []interface{}) []Expr {
	newExprs := []Expr{}
	for _, intf := range exprs {
		newExprs = append(newExprs, intfToExpr(intf))
	}
	return newExprs
}

func intfToExpr(intf interface{}) Expr {
	var expr Expr
	switch v := intf.(type) {
	case string:
		expr = Expr{F: v}
	case Expr:
		expr = v
	case []string:
		expr = Expr{F: strings.Join(v, " ")}
	case int:
		expr = Expr{F: "?", V: []interface{}{v}}
	default:
		expr = Expr{F: "UNKNOWN"}
	}
	return expr
}

func dialectFormat(dialect string, stmt string, params []interface{}) string {
	if dialect == pgsql {
		return strings.Replace(stmt, "?", fmt.Sprintf("$%d", len(params)), 1)
	}
	return stmt
}

func exprGroup(dialect string, exprs [][]Expr, params []interface{}) (string, []interface{}) {
	var sql string
	var newP []interface{}
	if len(exprs) > 0 {
		for i, group := range exprs {
			sql += "("
			for n, expr := range group {
				f := expr.F
				for _, v := range expr.V {
					newP = append(newP, v)
					// if dialect == pgsql {
					// 	f = strings.Replace(f, "?", fmt.Sprintf("$%d", len(params)+len(newP)), 1)
					// }
					f = dialectFormat(dialect, f, append(newP, params...))
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
	return sql, newP
}

func (q *Query) toSql(dialect string) (string, []interface{}, error) {

	sql := ""
	var params []interface{}

	if len(q.SE) > 0 {
		sql += "SELECT "
		sels := []string{}
		for _, s := range q.SE {
			expr := intfToExpr(s)
			if len(expr.V) > 0 {
				params = append(params, expr.V...)
			}
			sels = append(sels, dialectFormat(dialect, expr.F, params))
		}
		sql += strings.Join(sels, ", ")
		sql += " "
	}

	if len(q.FE) > 0 {
		sql += "FROM "
		tables := []string{}
		for _, s := range q.FE {
			expr := intfToExpr(s)
			if len(expr.V) > 0 {
				params = append(params, expr.V...)
			}
			tables = append(tables, dialectFormat(dialect, expr.F, params))
		}
		sql += strings.Join(tables, ", ")
		sql += " "
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

func (q *Query) Join(j string) *Query {
	q.J = append(q.J, j)
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

func (q *Query) GroupBy(gb string) *Query {
	q.GB = gb
	return q
}

func (q *Query) Having(exprs ...interface{}) *Query {
	newExprs := getExprs(exprs)
	q.H = append(q.H, newExprs)
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
