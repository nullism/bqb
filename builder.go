package bqb

import (
	"fmt"
	"strings"
)

const (
	PGSQL = "postgres"
	MYSQL = "mysql"
	RAW   = "raw"
	SQL   = "sql"

	paramPh = "{{xX_PARAM_Xx}}"
)

type Expr struct {
	F string
	V []interface{}
}

func GroupSep(sep string, enclose bool, exprs ...interface{}) Expr {
	var newFs []string
	var newV []interface{}
	for _, e := range exprs {
		expr := intfToExpr(e)
		newFs = append(newFs, expr.F)
		newV = append(newV, expr.V...)
	}
	pre := ""
	post := ""
	if enclose {
		pre = "("
		post = ")"
	}
	newF := pre + strings.Join(newFs, sep) + post
	return Expr{
		F: newF,
		V: newV,
	}
}

func Concat(exprs ...interface{}) Expr {
	return GroupSep("", false, exprs...)
}

func And(exprs ...interface{}) Expr {
	enclose := len(exprs) > 1
	return GroupSep(" AND ", enclose, exprs...)
}

func Or(exprs ...interface{}) Expr {
	enclose := len(exprs) > 1
	return GroupSep(" OR ", enclose, exprs...)
}

func V(expr string, vals ...interface{}) Expr {
	var params []interface{}
	tmpQ := "1xXX1_Y_2XXx2"
	newExpr := strings.ReplaceAll(expr, "??", tmpQ)

	for _, val := range vals {
		switch v := val.(type) {
		case []int:
			iparts := []string{}
			for _, intf := range v {
				iparts = append(iparts, paramPh)
				params = append(params, intf)
			}
			newPart := strings.Join(iparts, ", ")
			newExpr = strings.Replace(newExpr, "?", newPart, 1)
		case []string:
			iparts := []string{}
			for _, intf := range v {
				iparts = append(iparts, paramPh)
				params = append(params, intf)
			}
			newPart := strings.Join(iparts, ", ")
			newExpr = strings.Replace(newExpr, "?", newPart, 1)
		case []interface{}:
			iparts := []string{}
			for _, intf := range v {
				iparts = append(iparts, paramPh)
				params = append(params, intf)
			}
			newPart := strings.Join(iparts, ", ")
			newExpr = strings.Replace(newExpr, "?", newPart, 1)
		default:
			newExpr = strings.Replace(newExpr, "?", paramPh, 1)
			params = append(params, v)
		}
	}

	if strings.Contains(newExpr, "?") {
		panic(fmt.Sprintf("mismatched paramters for Valf: %v", expr))
	}

	return Expr{
		F: strings.ReplaceAll(newExpr, tmpQ, "?"),
		V: params,
	}
}

func exprsToSql(exprs []Expr) ([]string, []interface{}) {
	qs := []string{}
	var newP []interface{}

	for _, s := range exprs {
		expr := intfToExpr(s)
		if len(expr.V) > 0 {
			newP = append(newP, expr.V...)
		}
		qs = append(qs, expr.F)
	}
	return qs, newP
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
	case *Query:
		sql, params, err := v.toSql()
		if err != nil {
			panic("Error while parsing sub-query")
		}
		expr = Expr{F: sql, V: params}
	case string:
		v = strings.ReplaceAll(v, "??", "xXxXy__")
		if strings.Contains(v, "?") {
			panic(fmt.Sprintf("String value without parameters: %v", v))
		}
		v = strings.ReplaceAll(v, "xXxXy__", "?")
		expr = Expr{F: v}
	case int:
		expr = Expr{F: "?", V: []interface{}{v}}
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

func dialectReplace(dialect string, sql string, params []interface{}) string {
	for i, param := range params {
		if dialect == RAW {
			switch v := param.(type) {
			case nil:
				sql = strings.Replace(sql, paramPh, "NULL", 1)
			case int, bool:
				sql = strings.Replace(sql, paramPh, fmt.Sprintf("%v", v), 1)
			default:
				sql = strings.Replace(sql, paramPh, fmt.Sprintf("'%v'", v), 1)
			}
		} else if dialect == MYSQL || dialect == SQL {
			sql = strings.Replace(sql, paramPh, "?", 1)
		} else if dialect == PGSQL {
			sql = strings.Replace(sql, paramPh, fmt.Sprintf("$%d", i+1), 1)
		}
	}
	return sql
}
