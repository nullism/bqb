package bqb

import (
	"fmt"
	"strings"
)

type Insert struct {
	dialect string
	into    []Expr
	sel     []Expr
	from    []Expr
	cols    []Expr
	vals    []Expr
	where   []Expr
}

func InsertPsql() *Insert {
	return &Insert{dialect: PGSQL}
}

func InsertSql() *Insert {
	return &Insert{dialect: SQL}
}

func InsertMysql() *Insert {
	return &Insert{dialect: MYSQL}
}

func InsertRaw() *Insert {
	return &Insert{dialect: PGSQL}
}

func (i *Insert) Into(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.into = append(i.into, newExprs...)
	return i
}

func (i *Insert) Select(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.sel = append(i.sel, newExprs...)
	return i
}

func (i *Insert) From(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.from = append(i.from, newExprs...)
	return i
}

func (i *Insert) Cols(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.cols = append(i.cols, newExprs...)
	return i
}

func (i *Insert) Vals(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.vals = append(i.vals, newExprs...)
	return i
}

func (i *Insert) Where(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.where = append(i.where, newExprs...)
	return i
}

func (i *Insert) Print() {
	sql, params, err := i.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func (i *Insert) ToSql() (string, []interface{}, error) {
	sql := "INSERT "
	var params []interface{}

	if len(i.into) > 0 {
		sql += "INTO "
		nsql, nparams := exprsToSql(i.into)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(i.cols) > 0 {
		sql += "("
		nsql, nparams := exprsToSql(i.cols)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += ") "
	}

	if len(i.vals) > 0 {
		sql += "("
		nsql, nparams := exprsToSql(i.vals)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += ") "
	}

	if len(i.sel) > 0 {
		sql += "SELECT "
		nsql, nparams := exprsToSql(i.sel)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(i.from) > 0 {
		sql += "FROM "
		nsql, nparams := exprsToSql(i.from)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(i.where) > 0 {
		sql += "WHERE "
		nsql, nparams := exprsToSql(i.where)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	sql = dialectReplace(i.dialect, sql, params)
	return sql, params, nil
}
