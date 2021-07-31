package bqb

import (
	"fmt"
	"strings"
)

const (
	PGSQL   = "postgres"
	MYSQL   = "mysql"
	RAW     = "raw"
	paramPh = "xX_PARAM_Xx"
)

type Expr struct {
	F string
	V []interface{}
}

func Valf(expr string, vals ...interface{}) Expr {
	expr = strings.ReplaceAll(expr, "?", paramPh)
	e := Expr{
		F: expr,
		V: vals,
	}
	return e
}

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
	default:
		panic(fmt.Sprintf("Unsupported expression type: %T", v))
	}
	return expr
}

func exprGroup(exprs [][]Expr) (string, []interface{}) {
	var sql string
	var params []interface{}
	if len(exprs) > 0 {
		for i, group := range exprs {
			sql += "("
			for n, expr := range group {
				sql += fmt.Sprintf("%v", expr.F)
				println(fmt.Sprintf("%v %T", expr.V, expr.V))
				params = append(params, expr.V...)
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

	if len(q.SE) > 0 {
		sql += "SELECT "
		sels := []string{}
		for _, s := range q.SE {
			expr := intfToExpr(s)
			if len(expr.V) > 0 {
				params = append(params, expr.V...)

			}
			sels = append(sels, expr.F)
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
			tables = append(tables, expr.F)
		}
		sql += strings.Join(tables, ", ")
		sql += " "
	}

	if len(q.JE) > 0 {
		sql += "JOIN "
		tables := []string{}
		for _, s := range q.JE {
			expr := intfToExpr(s)
			if len(expr.V) > 0 {
				params = append(params, expr.V...)
			}
			tables = append(tables, expr.F)
		}
		sql += strings.Join(tables, ", ")
		sql += " "
	}
	fmt.Printf("param count is %d\n", len(params))

	if len(q.W) > 0 {
		sql += "WHERE "
		gsql, p := exprGroup(q.W)
		sql += gsql
		params = append(params, p...)
	}

	if len(q.GB) > 0 {
		sql += "GROUP BY "
		gbs := []string{}
		for _, s := range q.GB {
			expr := intfToExpr(s)
			if len(expr.V) > 0 {
				params = append(params, expr.V...)
			}
			gbs = append(gbs, expr.F)
		}
		sql += strings.Join(gbs, ", ")
		sql += " "
	}

	if len(q.H) > 0 {
		sql += "HAVING "
		hsql, hparams := exprGroup(q.H)
		sql += hsql
		params = append(params, hparams...)
	}

	if len(q.OB) > 0 {
		sql += "ORDER BY "
		obs := []string{}
		for _, s := range q.OB {
			expr := intfToExpr(s)
			if len(expr.V) > 0 {
				params = append(params, expr.V...)
			}
			obs = append(obs, expr.F)
		}
		sql += strings.Join(obs, ", ")
		sql += " "
	}

	if q.O != 0 {
		sql += fmt.Sprintf("OFFSET %d ", q.O)
	}

	if q.L != 0 {
		sql += fmt.Sprintf("LIMIT %d ", q.L)
	}

	fmt.Printf("param count %d\n", len(params))

	for i, p := range params {
		if dialect == RAW {
			switch v := p.(type) {
			case int:
				sql = strings.Replace(sql, paramPh, fmt.Sprintf("%v", v), 1)
			default:
				sql = strings.Replace(sql, paramPh, fmt.Sprintf("'%v'", v), 1)
			}
		} else if dialect == MYSQL {
			sql = strings.Replace(sql, paramPh, "?", 1)
		} else if dialect == PGSQL {
			sql = strings.Replace(sql, paramPh, fmt.Sprintf("$%d", i+1), 1)
		}

	}

	return sql, params, nil
}
