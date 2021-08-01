package bqb

import (
	"fmt"
	"strings"
)

type Insert struct {
	dialect   string
	into      []Expr
	union     []Expr
	sel       []Expr
	sel_query *Query
	from      []Expr
	cols      []Expr
	vals      []Expr
	where     []Expr
	limit     int
	offset    int
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

func (i *Insert) Cols(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.cols = append(i.cols, newExprs...)
	return i
}

func (i *Insert) Union(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.union = append(i.union, newExprs...)
	return i
}

func (i *Insert) Select(q *Query) *Insert {
	i.sel_query = q
	return i
}

func (i *Insert) Vals(exprs ...interface{}) *Insert {
	newExprs := getExprs(exprs)
	i.vals = append(i.vals, newExprs...)
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

	if len(i.union) > 0 {
		sql += "UNION "
		nsql, nparams := exprsToSql(i.union)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(i.vals) > 0 {
		sql += "("
		nsql, nparams := exprsToSql(i.vals)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += ") "
	}

	if i.sel_query != nil {
		qs, qp, err := i.sel_query.ToSql()
		if err != nil {
			return "", nil, err
		}
		sql += qs
		params = append(params, qp...)
	}

	sql = dialectReplace(i.dialect, sql, params)
	return sql, params, nil
}
