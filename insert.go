package bqb

import (
	"fmt"
	"strings"
)

type insert struct {
	dialect   string
	into      []Expr
	union     []Expr
	sel_query *selectQ
	cols      []Expr
	vals      []Expr
}

func Insert(exprs ...interface{}) *insert {
	return &insert{dialect: SQL, into: getExprs(exprs)}
}

func (i *insert) Postgres() *insert {
	i.dialect = PGSQL
	return i
}

func (i *insert) Mysql() *insert {
	i.dialect = MYSQL
	return i
}

func (i *insert) Raw() *insert {
	i.dialect = RAW
	return i
}

func (i *insert) Cols(exprs ...interface{}) *insert {
	newExprs := getExprs(exprs)
	i.cols = append(i.cols, newExprs...)
	return i
}

func (i *insert) Union(exprs ...interface{}) *insert {
	newExprs := getExprs(exprs)
	i.union = append(i.union, newExprs...)
	return i
}

func (i *insert) Select(q *selectQ) *insert {
	i.sel_query = q
	return i
}

func (i *insert) Vals(exprs ...interface{}) *insert {
	newExprs := getExprs(exprs)
	i.vals = append(i.vals, newExprs...)
	return i
}

func (i *insert) Print() {
	sql, params, err := i.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func (i *insert) ToSql() (string, []interface{}, error) {
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
		sql += "VALUES ("
		nsql, nparams := exprsToSql(i.vals)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += ") "
	}

	if i.sel_query != nil {
		qs, qp, err := i.sel_query.toSql()
		if err != nil {
			return "", nil, err
		}
		sql += qs
		params = append(params, qp...)
	}

	sql = dialectReplace(i.dialect, sql, params)
	return strings.TrimSpace(sql), params, nil
}
