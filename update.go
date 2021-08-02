package bqb

import (
	"fmt"
	"strings"
)

type Update struct {
	dialect string
	update  []Expr
	set     []Expr
	where   []Expr
}

func UpdatePsql() *Update {
	return &Update{dialect: PGSQL}
}

func UpdateSql() *Update {
	return &Update{dialect: SQL}
}

func UpdateMysql() *Update {
	return &Update{dialect: MYSQL}
}

func UpdateRaw() *Update {
	return &Update{dialect: PGSQL}
}

func (u *Update) Update(exprs ...interface{}) *Update {
	newExprs := getExprs(exprs)
	u.update = append(u.update, newExprs...)
	return u
}

func (u *Update) Set(exprs ...interface{}) *Update {
	newExprs := getExprs(exprs)
	u.set = append(u.set, newExprs...)
	return u
}

func (u *Update) Where(exprs ...interface{}) *Update {
	newExprs := getExprs(exprs)
	u.where = append(u.where, newExprs...)
	return u
}

func (u *Update) Print() {
	sql, params, err := u.ToSql()
	fmt.Printf("SQL: %v\n", sql)
	if len(params) > 0 {
		fmt.Printf("PARAMS: %v\n", params)
	}
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func (u *Update) ToSql() (string, []interface{}, error) {
	var sql string
	var params []interface{}

	if len(u.update) > 0 {
		sql += "UPDATE "
		nsql, nparams := exprsToSql(u.update)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(u.set) > 0 {
		sql += "SET "
		nsql, nparams := exprsToSql(u.set)
		sql += strings.Join(nsql, ", ")
		params = append(params, nparams...)
		sql += " "
	}

	if len(u.where) > 0 {
		sql += "WHERE "
		nsql, nparams := exprsToSql(u.where)
		sql += strings.Join(nsql, " ")
		params = append(params, nparams...)
		sql += " "
	}

	sql = dialectReplace(u.dialect, sql, params)
	return strings.TrimSpace(sql), params, nil
}
